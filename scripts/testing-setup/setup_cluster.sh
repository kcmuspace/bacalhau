#!/bin/bash

terraform apply --auto-approve

echo -n "Sleeping for 30s to allow ssh to start..."
sleep 30
echo "Done."

runRemote() {
  local args script

  remote_addr=$1; shift
  script=$1; shift

# generate eval-safe quoted version of current argument list
  printf -v args '%q ' "$@"

# pass that through on the command line to bash -s
 # note that $args is parsed remotely by /bin/sh, not by bash!
  ssh -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" ubuntu@"$remote_addr" "sudo bash -s -- $args" < "$script"
}

all_nodes_public=()
while IFS='' read -r line; do all_nodes_public+=("$line"); done <  <(terraform output -json | jq -r '.instance_public_dns.value | .[] ')

all_nodes_private=()
while IFS='' read -r line; do all_nodes_private+=("$line"); done <  <(terraform output -json | jq -r '.instance_private_ips.value | .[] ')

first_node="${all_nodes_public[1]}"
echo "Connecting to: ubuntu@$first_node"
runRemote $first_node ./setup_node.sh
peer_token=$(ssh -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" ubuntu@"$first_node" "cat /tmp/bacalhau_peer_token")

index=1
len_nodes=${#all_nodes_public[@]}

for i in "${!all_nodes_public[@]}"; do
    if (( index >= len_nodes)); then
      break
    fi

    this_node_public="${all_nodes_public[((i+1))]}"
    last_node_private="${all_nodes_private[((i))]}"

    PEER_STRING="--peer /ip4/$last_node_private/tcp/0/p2p/$peer_token"
    echo "Peer string: $PEER_STRING"

    echo "Connecting to: ubuntu@$this_node_public"
    echo "$PEER_STRING" | ssh -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" ubuntu@$this_node_public -T "cat > /tmp/PEER_STRING"
    runRemote "$this_node_public" ./setup_node.sh
    peer_token=$(ssh -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" ubuntu@"$this_node_public" "cat /tmp/bacalhau_peer_token")
    ((index=index+1))
done