## Prerequisites

You have to generate at least two files:

- with networks (one per line)
- with ips that belong to these networks (one per line)

You can use this generator

    https://github.com/moon-dragon-dev/go-gen-ip-networks

## How to build

    $ go build -o bin/load_networks cmd/load_networks/load_networks.go
    $ go build -o bin/check_ips cmd/check_ips/check_ips.go

## How to run

    $ bin/load_networks --user ${TNT_USER} --pass ${TNT_PASS} --host ${TNT_HOST} --port ${TNT_PORT} --networks-file ${NETWORKS_FILE}
    $ bin/check_ips --user ${TNT_USER} --pass ${TNT_PASS} --host ${TNT_HOST} --port ${TNT_PORT} --ips-file ${IPS_FILE}
