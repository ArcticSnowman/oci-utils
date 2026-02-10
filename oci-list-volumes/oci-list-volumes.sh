#!/usr/bin/env bash

set -euo pipefail

COMPARTMENT=""
CMPID=""
BOOTVOL=false
UNATTACHED=false
DEBUG=false

function usage() {
  echo "Usage: $0 [--compartment <compartment-name>] [--boot] [--unattached] [--debug]"
  echo ""
  echo "Options:"
  echo "  --compartment   Name of the compartment (required)"
  echo "  --boot          List only boot volumes (optional)"
  echo "  --unattached    List only unattached volumes (optional)"
  echo "  --debug         Enable debug mode (optional)"
  echo "  --help|-h       Show this help message"
  exit 1
}

# Get the direcory of the current script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"  
. "$SCRIPT_DIR/oci-functions.sh"

while [[ $# -gt 0 ]]; do
  case $1 in
    --compartment)
      COMPARTMENT="$2"
      shift # past argument
      shift # past value
      ;;
    --boot)
      BOOTVOL=true
      shift # past argument
      ;;
    --unattached)
      UNATTACHED=true
      shift # past argument
      ;;
    --debug)
      set -x
      DEBUG=true
      shift # past argument
      ;;
    --help|-h)
      usage
      ;;
    *) 
      echo "Unknown option: $1"
      usage
      ;;      
  esac
done

CMPID=$(get_compartment_id "$COMPARTMENT")

if [ "$BOOTVOL" = true ]; then
    OCIBVCMD="oci bv boot-volume"
    ARGS=""
    TITLE="Boot Volumes"
    OCIATTCHCMD="oci compute boot-volume-attachment"
    VOLUMEARGS="--boot-volume-id"
else
    OCIBVCMD="oci bv volume"
    ARGS="--lifecycle-state AVAILABLE"
    TITLE="Block Volumes"
    OCIATTCHCMD="oci compute volume-attachment"
    VOLUMEARGS="--volume-id"
fi

echo "Listing $TITLE in compartment $COMPARTMENT"

volidlist=$($OCIBVCMD list --compartment-id $CMPID --all $ARGS | jq -r '.data[] | .id')

if [ -z "$volidlist" ]; then
  echo "No volumes found"
  exit 0
fi

echo "Found $(echo "$volidlist" | wc -l) volumes"

echo "Getting volume details..."

header=0
for volid in $volidlist; do
    # if [ $header -eq 0 ]; then
    #     echo -e "ID\tDisplay-Name\tState\tSize(GB)\tCreated\tAttached"
    #     header=1
    # fi
    voljson=$($OCIBVCMD get $VOLUMEARGS "$volid" )
    volad=$(echo "$voljson" | jq -r '.data."availability-domain"')
    attchcount=$($OCIATTCHCMD list --compartment-id $CMPID --availability-domain $volad $VOLUMEARGS "$volid"| jq -r '.data | length')
    if [ -z "$attchcount" ]; then
        attached="No"
    else
        if [ "$UNATTACHED" == true ]; then
            continue
        fi
        attached="Yes"
    fi

    echo "$voljson" | jq --arg attached "$attached" -r '.data | [.id, ."display-name", ."size-in-gbs", ."time-created", $attached] | @tsv' 

    if [ "$DEBUG" == "true" ]; then
      break
    fi
done | column -t -N "ID,Display Name,Size(GB),Created,Attached" -s $'\t'


