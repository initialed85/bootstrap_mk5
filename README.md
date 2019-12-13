# bootstrap_mk5

This repo contains some tools and scripts to bootstrap a Cohda MK5 for 802.11p operation.

## What does it do?

- `bootstrap.sh`
    - Set the radio interface to operate in 802.11p
    - Using the last octet of the `eth0` MAC
        - Set the `wave-data` MAC
        - Set the `wave-data` IP
        - Set the `eth0` IP
        - Create routes for all 253 other possible nodes
        - relay broadcasts on port 6291 between eth0 and wave-data (avoiding loops)

## What are the components?

- `pkg`
    - `generate`
        - Functions used to generate MACs and IPs
    - `route`
        - Functions used to add routes
- `cmd`
    - `generate_mac`
        - Command line tool to generate a MAC address (from a MAC address)
    - `generate_ip`
        - Command line tool to generate an IP address (from a MAC address)
    - `add_routes`
        - Command line tool to add routes (within some guidelines)
    - `castinator`
        - Command line tool to relay UDP packets from one address to another
- `scripts`
    - `bootstrap.sh`
        - The script that pulls it all together

## What needs to be deployed?

- `/home/user`
    - `bootstrap_mk5`
        - `generate_mac`
        - `generate_ip`
        - `add_routes`
        - `castinator`
        - `bootstrap.sh`

Additionally, `/mnt/ubi/rc.local` must exist and must call to `/home/user/bootstrap_mk5/bootstrap.sh`

## How do I build it?

For your native platform:

    ./build.sh

For a Linux, ARM-based platform:

    GOOS=linux GOARCH=arm ./build.sh

## How do I test it?

Haha.
