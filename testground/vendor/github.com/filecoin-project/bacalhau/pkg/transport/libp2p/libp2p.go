package libp2p

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"time"

	realsync "sync"

	sync "github.com/lukemarsden/golang-mutex-tracer"

	"github.com/filecoin-project/bacalhau/pkg/config"
	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/filecoin-project/bacalhau/pkg/system"
	"github.com/filecoin-project/bacalhau/pkg/transport"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const JobEventChannel = "bacalhau-job-event"

type LibP2PTransport struct {
	// Cleanup manager for resource teardown on exit:
	cm *system.CleanupManager

	subscribeFunctions   []transport.SubscribeFn
	mutex                sync.RWMutex
	host                 host.Host
	port                 int
	peers                []string
	pubSub               *pubsub.PubSub
	jobEventTopic        *pubsub.Topic
	jobEventSubscription *pubsub.Subscription
	privateKey           crypto.PrivKey
}

func NewTransport(ctx context.Context, cm *system.CleanupManager, port int, peers []string) (*LibP2PTransport, error) {
	ctx, span := system.GetTracer().Start(ctx, "pkg/transport/libp2p.NewTransport")
	defer span.End()

	ctx, cancel := context.WithCancel(ctx)
	cm.RegisterCallback(func() error {
		cancel()
		return nil
	})

	usePeers := []string{}

	for _, p := range peers {
		if p != "" {
			usePeers = append(usePeers, p)
		}
	}

	prvKey, err := config.GetPrivateKey(fmt.Sprintf("private_key.%d", port))
	if err != nil {
		return nil, err
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	if err != nil {
		return nil, err
	}

	h, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		return nil, err
	}

	pgParams := pubsub.NewPeerGaterParams(
		0.33, //nolint:gomnd
		pubsub.ScoreParameterDecay(2*time.Minute),  //nolint:gomnd
		pubsub.ScoreParameterDecay(10*time.Minute), //nolint:gomnd
	)
	ps, err := pubsub.NewGossipSub(ctx, h, pubsub.WithPeerExchange(true), pubsub.WithPeerGater(pgParams))
	if err != nil {
		return nil, err
	}

	jobEventTopic, err := ps.Join(JobEventChannel)
	if err != nil {
		return nil, err
	}

	jobEventSubscription, err := jobEventTopic.Subscribe()
	if err != nil {
		return nil, err
	}

	libp2pTransport := &LibP2PTransport{
		cm:                   cm,
		subscribeFunctions:   []transport.SubscribeFn{},
		host:                 h,
		port:                 port,
		peers:                usePeers,
		privateKey:           prvKey,
		pubSub:               ps,
		jobEventTopic:        jobEventTopic,
		jobEventSubscription: jobEventSubscription,
	}

	libp2pTransport.mutex.EnableTracerWithOpts(sync.Opts{
		Threshold: 10 * time.Millisecond,
		Id:        "LibP2PTransport.mutex",
	})
	return libp2pTransport, nil
}

/*

  public api

*/

func (t *LibP2PTransport) HostID(ctx context.Context) (string, error) {
	return t.host.ID().String(), nil
}

func (t *LibP2PTransport) GetPeers(ctx context.Context) (map[string][]peer.ID, error) {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/transport/libp2p.GetPeers")
	defer span.End()

	response := map[string][]peer.ID{}
	for _, topic := range t.pubSub.GetTopics() {
		peers := t.pubSub.ListPeers(topic)
		response[topic] = peers
	}
	return response, nil
}

func (t *LibP2PTransport) Start(ctx context.Context) error {
	ctx, span := system.GetTracer().Start(ctx, "pkg/transport/libp2p.Start")
	defer span.End()

	if len(t.subscribeFunctions) == 0 {
		panic("Programming error: no subscribe func, please call Subscribe immediately after constructing interface")
	}

	t.cm.RegisterCallback(func() error {
		return t.Shutdown(ctx)
	})

	err := t.connectToPeers(ctx)
	if err != nil {
		return err
	}

	go t.listenForEvents(ctx)

	log.Trace().Msg("Libp2p transport has started")

	return nil
}

func (t *LibP2PTransport) Shutdown(ctx context.Context) error {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/transport/libp2p.Shutdown")
	defer span.End()

	closeErr := t.host.Close()

	if closeErr != nil {
		log.Error().Msgf("Libp2p transport had error stopping: %s", closeErr.Error())
	} else {
		log.Debug().Msg("Libp2p transport has stopped")
	}

	return nil
}

func (t *LibP2PTransport) Publish(ctx context.Context, ev model.JobEvent) error {
	return t.writeJobEvent(ctx, ev)
}

func (t *LibP2PTransport) Subscribe(ctx context.Context, fn transport.SubscribeFn) {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/transport/libp2p.Subscribe")
	defer span.End()

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.subscribeFunctions = append(t.subscribeFunctions, fn)
}

func (t *LibP2PTransport) Encrypt(ctx context.Context, data, libp2pKeyBytes []byte) ([]byte, error) {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/transport/libp2p.Encrypt")
	defer span.End()

	unmarshalledPublicKey, err := crypto.UnmarshalPublicKey(libp2pKeyBytes)
	if err != nil {
		return nil, err
	}
	publicKeyBytes, err := unmarshalledPublicKey.Raw()
	if err != nil {
		return nil, err
	}
	genericPublicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, ok := genericPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("could not cast public key to RSA")
	}
	return rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		rsaPublicKey,
		data,
		nil,
	)
}

func (t *LibP2PTransport) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/transport/libp2p.Decrypt")
	defer span.End()

	privateKeyBytes, err := t.privateKey.Raw()
	if err != nil {
		return nil, err
	}
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptOAEP(
		sha512.New(),
		rand.Reader,
		rsaPrivateKey,
		data,
		nil,
	)
}

/*

  libp2p

*/

func (t *LibP2PTransport) connectToPeers(ctx context.Context) error {
	ctx, span := system.GetTracer().Start(ctx, "pkg/transport/libp2p.Subscribe")
	defer span.End()

	if len(t.peers) == 0 {
		return nil
	}

	for _, peerAddress := range t.peers {
		maddr, err := multiaddr.NewMultiaddr(peerAddress)
		if err != nil {
			return err
		}

		// Extract the peer ID from the multiaddr.
		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			return err
		}

		t.host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
		err = t.host.Connect(ctx, *info)
		if err != nil {
			return err
		}
		log.Trace().Msgf("Libp2p transport connected to: %s", peerAddress)
	}

	return nil
}

/*

  pub / sub

*/

// we wrap our events on the wire in this envelope so
// we can pass our tracing context to remote peers
type jobEventEnvelope struct {
	SentTime  time.Time              `json:"sent_time"`
	JobEvent  model.JobEvent         `json:"job_event"`
	TraceData propagation.MapCarrier `json:"trace_data"`
}

func (t *LibP2PTransport) writeJobEvent(ctx context.Context, event model.JobEvent) error {
	ctx, span := system.GetTracer().Start(ctx, "pkg/transport/libp2p.writeJobEvent")
	defer span.End()

	traceData := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, &traceData)

	bs, err := json.Marshal(jobEventEnvelope{
		JobEvent:  event,
		TraceData: traceData,
		SentTime:  time.Now(),
	})
	if err != nil {
		return err
	}

	log.Trace().Msgf("Sending event %s: %s", event.EventName.String(), string(bs))
	return t.jobEventTopic.Publish(ctx, bs)
}

func (t *LibP2PTransport) readMessage(msg *pubsub.Message) {
	// TODO: we would enforce the claims to SourceNodeID here
	// i.e. msg.ReceivedFrom() should match msg.Data.JobEvent.SourceNodeID
	payload := jobEventEnvelope{}
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		log.Error().Msgf("error unmarshalling libp2p event: %v", err)
		return
	}

	now := time.Now()
	then := payload.SentTime
	latency := now.Sub(then)
	latencyMilli := int64(latency / time.Millisecond)
	if latencyMilli > 500 { //nolint:gomnd
		log.Warn().Msgf(
			"[%s=>%s] VERY High message latency: %d ms (%s)",
			payload.JobEvent.SourceNodeID[:8],
			t.host.ID().String()[:8],
			latencyMilli, payload.JobEvent.EventName.String(),
		)
	} else if latencyMilli > 50 { //nolint:gomnd
		log.Warn().Msgf(
			"[%s=>%s] High message latency: %d ms (%s)",
			payload.JobEvent.SourceNodeID[:8],
			t.host.ID().String()[:8],
			latencyMilli, payload.JobEvent.EventName.String(),
		)
	} else {
		log.Trace().Msgf(
			"[%s=>%s] Message latency: %d ms (%s)",
			payload.JobEvent.SourceNodeID[:8],
			t.host.ID().String()[:8],
			latencyMilli, payload.JobEvent.EventName.String(),
		)
	}

	log.Trace().Msgf("Received event %s: %+v", payload.JobEvent.EventName.String(), payload)

	// Notify all the listeners in this process of the event:
	jobCtx := otel.GetTextMapPropagator().Extract(context.Background(), payload.TraceData)

	ev := payload.JobEvent
	// NOTE: Do not use msg.ReceivedFrom as the original sender, it's not. It's
	// the node which gossiped the message to us, which might be different.
	// (was: ev.SourceNodeID = msg.ReceivedFrom.String())
	ev.SenderPublicKey = msg.Key

	var wg realsync.WaitGroup
	func() {
		t.mutex.RLock()
		defer t.mutex.RUnlock()

		for _, fn := range t.subscribeFunctions {
			wg.Add(1)
			go func(f transport.SubscribeFn) {
				defer wg.Done()
				f(jobCtx, ev)
			}(fn)
		}
	}()
	wg.Wait()
}

func (t *LibP2PTransport) listenForEvents(ctx context.Context) {
	for {
		msg, err := t.jobEventSubscription.Next(ctx)
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				log.Trace().Msgf("libp2p transport shutting down: %v", err)
			} else {
				log.Error().Msgf(
					"libp2p encountered an unexpected error, shutting down: %v", err)
			}
			return
		}
		go t.readMessage(msg)
	}
}

// Compile-time interface check:
var _ transport.Transport = (*LibP2PTransport)(nil)
