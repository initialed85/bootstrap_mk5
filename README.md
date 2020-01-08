# bootstrap_mk5

This repo contains some tools and scripts to bootstrap a Cohda MK5 for 802.11p operation.

## Goal

The code in this repo seeks to set up a Cohda MK5 in 802.11p mode, operating as part of a /24 network (on the wireless side) and tying in a /24 network (on the wired side).

So, if we look at a network of 3 radios:

- cohda-mk5-1
    - eth0 = 192.168.1.1/24
    - wave-data = 172.16.137.1/24
    - routes
        - 192.168.2.0/24 via 172.16.137.2
        - 192.168.3.0/24 via 172.16.137.3
- cohda-mk5-2
    - eth0 = 192.168.2.1/24
    - wave-data = 172.16.137.2/24
    - routes
        - 192.168.1.0/24 via 172.16.137.1
        - 192.168.3.0/24 via 172.16.137.3
- cohda-mk5-3
    - eth0 = 192.168.3.1/24
    - wave-data = 172.16.137.3/24
    - routes
        - 192.168.1.0/24 via 172.16.137.1
        - 192.168.2.0/24 via 172.16.137.2
        
As you can see, there'll be a lot of static routes, so this is handled by some Go command line apps and some scripting.

The problem is probably not a bad argument for a dynamic routing protocol, but then you're at the mercy of that routing protocol's discovery intervals. 

## Prerequisites

- Linux or MacOS
- go1.12 (or newer, probably)
- ping6
- sshpass
    - see [https://gist.github.com/arunoda/7790979](https://gist.github.com/arunoda/7790979)
- ansible

## What does it do?

- `bootstrap.sh`
    - Set the radio interface to operate in 802.11p
    - Using the last octet of the `eth0` MAC
        - Set the `wave-data` MAC
        - Set the `wave-data` IP
        - Set the `eth0` IP
        - Create routes for all 253 other possible nodes
        - relay multicasts on port 13337 on eth0 to broadcasts on port 13338 on wave-data (and vice versa)

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

Heh, TBD.

## How do I deploy it?

Ensure you've built for the ARM-based platform:
    
    GOOS=linux GOARCH=arm ./build.sh
    
Then run the deploy script

    ./deploy.sh
