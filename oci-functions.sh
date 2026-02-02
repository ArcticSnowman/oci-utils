#!/usr/bin/env bash


# Get the compartment ID from the compartment name
function get_compartment_id() {

  local compartment_name="$1"

  if [ -z "$compartment_name" ]; then
    printf "Compartment Name is required\n" > /dev/stderr
    usage > /dev/stderr
    return ""
  fi

  oci iam compartment list --compartment-id-in-subtree true --include-root --all --query "data[?name=='$compartment_name'].id | [0]" --raw-output
}

# Get the cluster ID from the cluster name
function get_cluster_id() {
  local compartment_id="$1"
  local cluster_name="$2"

  if [ -z "$cluster_name" ]; then
    printf "Cluster Name is required\n" > /dev/stderr
    usage  > /dev/stderr
    return ""
  fi

  oci ce cluster list --compartment-id "$compartment_id" --lifecycle-state "ACTIVE" --query "data[?name=='$cluster_name'].id | [0]" --raw-output
}

# Get the node pool ID from the node pool name
function get_nodepool_id() {
  local compartment_id="$1"
  local cluster_id="$2"
  local nodepool_name="$3"
  oci ce node-pool list --compartment-id "$compartment_id" --cluster-id "$cluster_id" --lifecycle-state "ACTIVE" --lifecycle-state "UPDATING" --query "data[?name=='$nodepool_name'].id | [0]" --raw-output
}

function get_image_id() {
  local image_name="$1"
  local compartment_id="$2"

  if [ -z "$image_name" ]; then
    printf "Image Name is required\n" > /dev/stderr
    usage > /dev/stderr
    return ""
  fi

  oci compute image list --compartment-id "$compartment_id" --all | jq -r ".data[] | select(.\"display-name\"==\"$image_name\").id"
}


function get_compartment_name() {
  local cmpid="$1"

  oci iam compartment get --compartment-id=$cmpid | jq -r '.data.name'
}

function get_compartment_parent() {
  local cmpid="$1"

  oci iam compartment get --compartment-id=$cmpid | jq -r '.data."compartment-id"'
}

