# bootstrap_mk5

This repo contains some tools and scripts to bootstrap a Cohda MK5 for 802.11p operation.

## Context

- Cohda MK5
    - A V2X hardware platform developed by Cohda Wireless
- 802.11p
    - An 802.11a-based standard that provides fast peering by doing away with BSSes (everyone is always on 00:00:00:00:00:00)
- eth0
    - The wired interface on a Cohda MK5
- wave-data
    - The wireless interface on a Cohda MK5

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
        
As you can see, there'll be a lot of static routes, so this is handled by some Go command line tools (that get cross-compiled for ARM, i.e. the Cohda MK5) and some scripting.

The problem is probably not a bad argument for a dynamic routing protocol, but then you're at the mercy of that routing protocol's discovery intervals. 

## Prerequisites

- Linux or MacOS
- go1.12 (or newer, probably)
- ping6
- sshpass
    - see [https://gist.github.com/arunoda/7790979](https://gist.github.com/arunoda/7790979)
- ansible

## What does it do?

- With `rc.local` as the entrypoint
    - Call `bootstrap.sh`
        - Using the Cohda `chconfig` tool
            - Set the radio interface to operate in 802.11p
        - Using the `generate_mac` tool and the last octet of the `eth0` MAC
            - Set the `wave-data` MAC
        - Using the `generate_ip` tool and the last octet of the `eth0` MAC
            - Set the `wave-data` IP
            - Set the `eth0` IP
        - Using the `add_routes` tool
            - Create routes for all 253 other possible nodes
        - Using the `castinator` tool
            - Provide a bidirectional relay between UDP broadcasts on port 13337 on `eth0` and UDP broadcasts on port 13338 on `wave-data`
        - Using the `locator` tool
            - Forward `TPV` packets from `gpsd` as UDP broadcasts port 13339 on `eth0`

## What are the components?

This project is laid out largely in the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) format; here is an breakdown of the folders:

- `cmd` (Go code that gets built to a command line executable)
    - `add_routes`
        - Tool to add routes (within some guidelines)
    - `castinator`
        - Tool to relay UDP packets from one address to another
    - `find_mk5s`
        - Tool used to discover MK5s plugged in locally via Ethernet
    - `generate_ip`
        - Tool to generate an IP address (from a MAC address, within some guidelines)
    - `generate_mac`
        - Tool to generate a MAC address (from a MAC address, within some guidelines)
    - `locator`
        - Tool to receive `TPV` packets from `gpsd` and send them a UDP address
- `deploy` (Files and executables that relate to deployment)
    - `ansible.cfg`
        - Ansible config
    - `deploy-no-reboot.yml`
        - Ansible Playbook
    - `hosts-base`
        - Base template for Ansible hosts file
- `dist` (Where artifacts to-be-deployed are built to)
- `internal` (Go code used only by this repo)
- `pkg` (Go code that can be used externally)
    - `generate`
        - Functions used to generate MACs and IPs
    - `route`
        - Functions used to add routes
- `scripts`
    - `bootstrap.sh`
        - The script that pulls everything together
    - `rc.local`
        - The entrypoint script

## What gets deployed?

    ./dist/generate_mac -> /home/user/bootstrap_mk5/generate_mac
    ./dist/castinator -> /home/user/bootstrap_mk5/castinator
    ./dist/add_routes -> /home/user/bootstrap_mk5/add_routes
    ./dist/generate_ip -> /home/user/bootstrap_mk5/generate_ip
    
    ./scripts/bootstrap.sh -> /home/user/bootstrap_mk5/bootstrap.sh

    ./dist/build_hash.txt -> /home/user/bootstrap_mk5/build_hash.txt
    ./dist/build_date.txt -> /home/user/bootstrap_mk5/build_date.txt
    
    ./scripts/rc.local -> /mnt/ubi/rc.local

    ./deploy/deploy_date.txt -> /home/user/bootstrap_mk5/deploy_date.txt
    
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
    
NOTE: You may pass the optional `reboot` flag to `deploy.sh`; e.g.: `./deploy.sh reboot`

## What might I want to change?

- The Cohda IPv6 link local address prefix for discovery (e.g. if the hardware revision changes)
    - `cmd/find_mk5s/main.go`
        - `CohdaIPv6AddrPrefix = "fe80::6e5:"`
- The Cohda MK5 credentials for deployment
    - `deploy/hosts-base`
        - `ansible_ssh_user=user`
        - `ansible_ssh_pass=user`
- The IP and MAC address specifics in `scripts/bootstrap.sh`
    - `ETH_INTERFACE=eth0`
    - `ETH_BASE_IP=192.168.x.1`
    - `ETH_NETMASK=255.255.255.0`
    - `ETH_PREFIX=24`
    - `WAVE_DATA_INTERFACE=wave-data`
    - `WAVE_DATA_BASE_IP=172.16.137.x`
    - `WAVE_DATA_NETMASK=255.255.255.0`
- The rate and channel specifics in `scripts/bootstrap.sh`
    - `DEFAULT_MCS=MK2MCS_R34QAM64 # 54 Mbps`
    - `DEFAULT_TX_POWER=5          # 2.5 dBm (5 * 0.5 dBm)`
    - `CHANNEL_NUMBER=174          # ETSI = 23 dBm, SCH; IEEE = 33 dBm, SCH`
- The multicast-to-broadcast relay specifics in `scripts/bootstrap.sh`
    - `ETH_PORT=13337`
    - `WAVE_DATA_PORT=13338`
- The GPS forwarder specifics in `scripts/bootstrap.sh`
    - `GPS_PORT=13339`
