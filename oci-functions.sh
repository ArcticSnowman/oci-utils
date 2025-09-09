#!/usr/bin/env bash


# Get the compartment ID from the compartment name
function get_compartment_id() {

  local compartment_name="$1"

  if [ -z "$compartment_name" ]; then
    echo "Compartment Name is required"
    exit 1
  fi

  oci iam compartment list --compartment-id-in-subtree true --all --query "data[?name=='$compartment_name'].id | [0]" --raw-output
}

# Get the cluster ID from the cluster name
function get_cluster_id() {
  local compartment_id="$1"
  local cluster_name="$2"
  oci ce cluster list --compartment-id "$compartment_id" --query "data[?name=='$cluster_name' && .'lifecycle-state' == 'ACTIVE'].id | [0]" --raw-output
}