# oci-utils

A collection of scripts and tools for interacting with Oracle Cloud Infrastructure (OCI) resources.

## Scripts Included

- install.sh: Script to install the utilities.
- oci-functions.sh: Common functions used by other scripts.
- oci-list-buckets: List buckets in a specified compartment.
- oci-list-compute-images: List available compute images in OCI.
- oci-list-compartments: List compartments in your OCI tenancy.
- oci-list-log-groups: List log groups in a specified compartment.
- oci-list-logs: List logs in a specified log group.
- oci-list-oke-clusters: List Oracle Kubernetes Engine (OKE) clusters.
- oci-list-oke-node-pools: List node pools in OKE clusters.
- oci-list-oke-nodes: List nodes in OKE clusters.
- oci-list-subnets: List subnets in your OCI tenancy.
- oci-list-volumes: List block/boot volumes in OCI compartment.
- oci-oke-create-kubeconfig: Create kubeconfig file for accessing OKE clusters.
- oci-policy-count-summary: Summarize policy counts by compartment.
- oci-scan-policies: List all policies that apply to a given compartment.
- oci-scan-logs: Search and display logs from OCI logging service.
- oci-transfer-custom-images: Transfer custom images between compartments.
- oci-compute-capacity-report: Report available compute capacity by shape, availability domain, and fault domain.
- oci-compute-reservation-report: Report reserved compute capacity by shape, availability domain, and fault domain.


## Prerequisites

- Bash shell
- OCI CLI installed and configured ([OCI CLI documentation](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm))


## Installation

Clone the repository:

```bash
git clone https://github.com/ArcticSnowman/oci-utils.git
cd oci-utils

```

Install dependencies:

```bash
sudo apt-get install jq column

```
Run Installation script:

```bash
./install.sh
```

## Usage

Each script can be executed directly from the command line. For example:

```bash
oci-list-compute-images
```

## Compute Reservation Report

`oci-compute-reservation-report` lists compute capacity reservations in a compartment and reports the reserved instance count by shape, availability domain, and fault domain. The compartment is specified by name and is required. OCI credentials are read from the selected profile in `~/.oci/config`; the profile defaults to `DEFAULT`.

Build the executable from the repository:

```bash
cd oci-compute-reservation-report
make compile
```

Run the report by specifying a compartment and, optionally, a profile:

```bash
./bin/oci-compute-reservation-report --compartment "engineering"
./bin/oci-compute-reservation-report --compartment "engineering" --profile "PRODUCTION"
```

Use `--debug` (or `-d`) for debug logging. Run `./bin/oci-compute-reservation-report --help` to see all options.

## Compute Capacity Report

`oci-compute-capacity-report` queries OCI for available compute capacity and displays the available count by availability domain and fault domain. Specify either a Karpenter node pool with `--nodepool`, or all of `--shape`, `--ocpu`, and `--memory`. The availability domain can optionally be filtered with `--availabilitydomain`. OCI credentials are read from the selected profile in `~/.oci/config`; the profile defaults to `DEFAULT`.

Build the executable from the repository:

```bash
cd oci-compute-capacity-report
make compile
```

Examples:

```bash
./bin/oci-compute-capacity-report --shape VM.Standard.E5.Flex --ocpu 4 --memory 64
./bin/oci-compute-capacity-report --shape VM.Standard.E5.Flex --ocpu 4 --memory 64 --availabilitydomain "Uocm:PHX-AD-1" --profile "PRODUCTION"
./bin/oci-compute-capacity-report --nodepool "general-purpose"
```

Use `--debug` (or `-d`) for debug logging. Run `./bin/oci-compute-capacity-report --help` to see all options.

## License

MIT License

