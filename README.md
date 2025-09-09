# oci-utils

A collection of shell scripts for interacting with Oracle Cloud Infrastructure (OCI) resources.

## Scripts Included

- oci-functions.sh: Common functions used by other scripts.
- oci-list-ads-images: List available images in OCI.
- oci-list-compartments: List compartments in your OCI tenancy.
- oci-list-oke-clusters: List Oracle Kubernetes Engine (OKE) clusters.
- oci-list-oke-node-pools: List node pools in OKE clusters.
- oci-list-oke-nodes: List nodes in OKE clusters.
- oci-list-subnets: List subnets in your OCI tenancy.

## Prerequisites

- Bash shell
- OCI CLI installed and configured ([OCI CLI documentation](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm))

## Usage

1. Clone this repository:
   ```bash
   git clone https://github.com/ArcticSnowman/oci-utils.git
   cd oci-utils
   ```

2. Source the functions (if needed):
   ```bash
   source oci-functions.sh
   ```

3. Run any script:
   ```bash
   ./oci-list-compartments
   ```

## License

MIT License

