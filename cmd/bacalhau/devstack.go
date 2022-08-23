package bacalhau

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/filecoin-project/bacalhau/pkg/capacitymanager"
	"github.com/filecoin-project/bacalhau/pkg/computenode"
	"github.com/filecoin-project/bacalhau/pkg/config"
	"github.com/filecoin-project/bacalhau/pkg/controller"
	"github.com/filecoin-project/bacalhau/pkg/devstack"
	"github.com/filecoin-project/bacalhau/pkg/executor"
	noop_executor "github.com/filecoin-project/bacalhau/pkg/executor/noop"
	executor_util "github.com/filecoin-project/bacalhau/pkg/executor/util"
	"github.com/filecoin-project/bacalhau/pkg/publisher"
	publisher_util "github.com/filecoin-project/bacalhau/pkg/publisher/util"
	"github.com/filecoin-project/bacalhau/pkg/storage"
	"github.com/filecoin-project/bacalhau/pkg/system"
	"github.com/filecoin-project/bacalhau/pkg/util/templates"
	"github.com/filecoin-project/bacalhau/pkg/verifier"
	verifier_util "github.com/filecoin-project/bacalhau/pkg/verifier/util"
	"github.com/rs/zerolog/log"
	"k8s.io/kubectl/pkg/util/i18n"

	"github.com/spf13/cobra"
)

var (
	devStackLong = templates.LongDesc(i18n.T(`
		Start a cluster of nodes and run a job on them.
`))

	//nolint:lll // Documentation
	devstackExample = templates.Examples(i18n.T(`
		# Create a devstack cluster.
		bacalhau devstack
`))

	// Set Defaults (probably a better way to do this)
	ODs = NewDevStackOptions()

	// For the -f flag
)

type DevStackOptions struct {
	NumberOfNodes     int    // Number of nodes to start in the cluster
	NumberOfBadActors int    // Number of nodes to be bad actors
	IsNoop            bool   // Noop executor and verifier for all jobs
	Peer              string // Connect node 0 to another network node
}

func NewDevStackOptions() *DevStackOptions {
	return &DevStackOptions{
		NumberOfNodes:     3,
		NumberOfBadActors: 0,
		IsNoop:            false,
		Peer:              "",
	}
}

func init() { //nolint:gochecknoinits // Using init in cobra command is idomatic
	devstackCmd.PersistentFlags().IntVar(
		&ODs.NumberOfNodes, "nodes", ODs.NumberOfNodes,
		`How many nodes should be started in the cluster`,
	)
	devstackCmd.PersistentFlags().IntVar(
		&ODs.NumberOfBadActors, "bad-actors", ODs.NumberOfBadActors,
		`How many nodes should be bad actors`,
	)
	devstackCmd.PersistentFlags().BoolVar(
		&ODs.IsNoop, "noop", false,
		`Use the noop executor and verifier for all jobs`,
	)
	devstackCmd.PersistentFlags().StringVar(
		&ODs.Peer, "peer", ODs.Peer,
		`Connect node 0 to another network node`,
	)

	setupJobSelectionCLIFlags(devstackCmd)
	setupCapacityManagerCLIFlags(devstackCmd)
}

var devstackCmd = &cobra.Command{
	Use:     "devstack",
	Short:   "Start a cluster of bacalhau nodes for testing and development",
	Long:    devStackLong,
	Example: devstackExample,
	RunE: func(cmd *cobra.Command, args []string) error { // nolintunparam // incorrect lint that is not used
		config.DevstackSetShouldPrintInfo()

		if ODs.NumberOfBadActors >= ODs.NumberOfNodes {
			return fmt.Errorf("cannot have more bad actors than there are nodes")
		}

		// Cleanup manager ensures that resources are freed before exiting:
		cm := system.NewCleanupManager()
		cm.RegisterCallback(system.CleanupTracer)
		defer cm.Cleanup()

		// Context ensures main goroutine waits until killed with ctrl+c:
		ctx, cancel := system.WithSignalShutdown(context.Background())
		defer cancel()

		getStorageProviders := func(ipfsMultiAddress string, nodeIndex int) (
			map[storage.StorageSourceType]storage.StorageProvider, error) {

			if ODs.IsNoop {
				return executor_util.NewNoopStorageProviders(cm)
			}

			return executor_util.NewStandardStorageProviders(cm, ipfsMultiAddress)
		}

		getExecutors := func(ipfsMultiAddress string, nodeIndex int, ctrl *controller.Controller) (
			map[executor.EngineType]executor.Executor, error) {

			if ODs.IsNoop {
				return executor_util.NewNoopExecutors(cm, noop_executor.ExecutorConfig{})
			}

			return executor_util.NewStandardExecutors(cm,
				ipfsMultiAddress, fmt.Sprintf("devstacknode%d", nodeIndex))
		}

		getVerifiers := func(
			ipfsMultiAddress string,
			nodeIndex int,
			ctrl *controller.Controller,
		) (map[verifier.VerifierType]verifier.Verifier, error) {
			if ODs.IsNoop {
				return verifier_util.NewNoopVerifiers(cm, ctrl.GetStateResolver())
			}
			return verifier_util.NewNoopVerifiers(cm, ctrl.GetStateResolver())
		}

		getPublishers := func(
			ipfsMultiAddress string,
			nodeIndex int,
			ctrl *controller.Controller,
		) (map[publisher.PublisherType]publisher.Publisher, error) {
			if ODs.IsNoop {
				return publisher_util.NewNoopPublishers(cm, ctrl.GetStateResolver())
			}
			return publisher_util.NewIPFSPublishers(cm, ctrl.GetStateResolver(), ipfsMultiAddress)
		}

		jobSelectionPolicy := getJobSelectionConfig()
		totalResourceLimit, jobResourceLimit := getCapacityManagerConfig()

		computeNodeConfig := computenode.ComputeNodeConfig{
			JobSelectionPolicy: jobSelectionPolicy,
			CapacityManagerConfig: capacitymanager.Config{
				ResourceLimitTotal: totalResourceLimit,
				ResourceLimitJob:   jobResourceLimit,
			},
		}

		portFileName := "/tmp/bacalhau-devstack.port"
		_, err := os.Stat(portFileName)
		if err == nil {
			log.Fatal().Msgf("Found file %s - Devstack likely already running", portFileName)
		}
		pidFileName := "/tmp/bacalhau-devstack.pid"
		_, err = os.Stat(pidFileName)
		if err == nil {
			log.Fatal().Msgf("Found file %s - Devstack likely already running", pidFileName)
		}

		stack, err := devstack.NewDevStack(
			cm,
			ODs.NumberOfNodes,
			ODs.NumberOfBadActors,
			getStorageProviders,
			getExecutors,
			getVerifiers,
			getPublishers,
			computeNodeConfig,
			ODs.Peer,
			false,
		)
		if err != nil {
			return err
		}

		stack.PrintNodeInfo()

		f, err := os.Create(portFileName)
		if err != nil {
			log.Fatal().Msgf("Error writing out port file to %v", portFileName)
		}
		defer os.Remove(portFileName)
		firstNode := stack.Nodes[0]
		_, err = f.WriteString(strconv.Itoa(firstNode.APIServer.Port))
		if err != nil {
			log.Fatal().Msgf("Error writing out port file: %v", portFileName)
		}

		fPid, err := os.Create(pidFileName)
		if err != nil {
			log.Fatal().Msgf("Error writing out pid file to %v", pidFileName)
		}
		defer os.Remove(pidFileName)

		_, err = fPid.WriteString(strconv.Itoa(os.Getpid()))
		if err != nil {
			log.Fatal().Msgf("Error writing out pid file: %v", pidFileName)
		}

		<-ctx.Done() // block until killed
		return nil
	},
}
