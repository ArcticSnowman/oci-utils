# oci-utils

A collection of shell scripts for interacting with Oracle Cloud Infrastructure (OCI) resources.

## Scripts Included

- oci-functions.sh: Common functions used by other scripts.
- oci-list-custom-images: List available images in OCI.
- oci-list-compartments: List compartments in your OCI tenancy.
- oci-list-oke-clusters: List Oracle Kubernetes Engine (OKE) clusters.
- oci-list-oke-node-pools: List node pools in OKE clusters.
- oci-list-oke-nodes: List nodes in OKE clusters.
- oci-list-subnets: List subnets in your OCI tenancy.
- oci-recycle-node-pool: Recycle nodes in an OKE node pool.

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
oci-list-custom-images
```

## License

MIT License

