#!/usr/bin/env bash

set -x
set -e

SCRIPT_DIR=$(dirname "$(readlink -f "${0}")")

pushd "$(pwd)"

catch() {
  echo "error: caught a failure; popping original directory and exiting"

  popd

  exit 1
}

trap 'catch' ERR

cd "${SCRIPT_DIR}"

#
# things you may want to change (if you move the wildcard octet, you'll need to adjust the octet settings further down)
#

ETH_INTERFACE=eth0
ETH_BASE_IP=192.168.x.1
ETH_NETMASK=255.255.255.0
ETH_PREFIX=24
ETH_PORT=13337

WAVE_DATA_INTERFACE=wave-data
WAVE_DATA_BASE_IP=172.16.137.x
WAVE_DATA_NETMASK=255.255.255.0
WAVE_DATA_PORT=13338

# MCS rate options (same as 802.11a)
# MK2MCS_R12BPSK  - 6 Mbps
# MK2MCS_R34BPSK  - 9 Mbps
# MK2MCS_R12QPSK  - 12 Mbps
# MK2MCS_R34QPSK  - 18 Mbps
# MK2MCS_R12QAM16 - 24 Mbps
# MK2MCS_R34QAM16 - 36 Mbps
# MK2MCS_R23QAM64 - 48 Mbps
# MK2MCS_R34QAM64 - 54 Mbps
# MK2MCS_DEFAULT  - ?
# MK2MCS_TRC      - ?

DEFAULT_MCS=MK2MCS_R34QAM64 # 54 Mbps
DEFAULT_TX_POWER=5          # 2.5 dBm (5 * 0.5 dBm)
CHANNEL_NUMBER=174          # ETSI = 23 dBm, SCH; IEEE = 33 dBm, SCH

#
# things you probably don't want to change
#

INTERFACE=wave-raw
CHANNEL=SCH
BW=MK2BW_10MHz
RX_ANT=3
RADIO=A

ETH_IP_ADDR=$(./generate_ip -baseIPAddr ${ETH_BASE_IP} -identifierOctet 3 -interfaceName ${ETH_INTERFACE})
ETH_BROADCAST_IP=$(./generate_ip -baseIPAddr "${ETH_IP_ADDR}" -identifierOctet 4 -specificIdentifier 255)
ETH_NETWORK_IP=$(./generate_ip -baseIPAddr "${ETH_IP_ADDR}" -identifierOctet 4 -specificIdentifier 0)

WAVE_DATA_MAC_ADDR=$(./generate_mac -interfaceName ${ETH_INTERFACE})
WAVE_DATA_IP_ADDR=$(./generate_ip -baseIPAddr ${WAVE_DATA_BASE_IP} -identifierOctet 4 -interfaceName ${ETH_INTERFACE})
WAVE_DATA_BROADCAST_IP=$(./generate_ip -baseIPAddr "${WAVE_DATA_IP_ADDR}" -identifierOctet 4 -specificIdentifier 255)
WAVE_DATA_BASE_GW_IP=$(./generate_ip -baseIPAddr "${WAVE_DATA_IP_ADDR}" -identifierOctet 4 -specificIdentifier 1)

#
# things you should not change
#

rmmod ieee1609dot3 || true
rmmod ieee1609dot4 || true
insmod /opt/cohda/drivers/ieee1609dot4.ko mac=1

# ethertype 0x88b6 is one available for development / experimentation (used by Cohda I think?)
# ethertype 0x86dd is IPv6
# ethertype 0x0800 is IPv4
# ethertype 0x0806 is ARP

chconfig \
  --Set \
  --Filter 0x88b6,0x86dd,0x0800,0x0806 \
  --Interface ${INTERFACE} \
  --Channel ${CHANNEL} \
  --DefaultMCS ${DEFAULT_MCS} \
  --DefaultTxPower ${DEFAULT_TX_POWER} \
  --ChannelNumber ${CHANNEL_NUMBER} \
  --DefaultTRC 0 \
  --DefaultTPC 0 \
  --BW ${BW} \
  --DualTx MK2TXC_TXRX \
  --RxAnt ${RX_ANT} \
  --Radio ${RADIO} \
  --MACAddr "${WAVE_DATA_MAC_ADDR}"

ifconfig ${WAVE_DATA_INTERFACE} "${WAVE_DATA_IP_ADDR}" netmask ${WAVE_DATA_NETMASK} broadcast "${WAVE_DATA_BROADCAST_IP}"

ifconfig ${ETH_INTERFACE} "${ETH_IP_ADDR}" netmask ${ETH_NETMASK} broadcast "${ETH_BROADCAST_IP}"

./add_routes \
  -baseDstIPAddr "${ETH_NETWORK_IP}" \
  -dstIdentifierOctet 3 \
  -dstPrefix ${ETH_PREFIX} \
  -baseGwIPAddr "${WAVE_DATA_BASE_GW_IP}" \
  -gwIdentifierOctet 4 \
  -startIdentifier 1 \
  -stopIdentifier 254 \
  -skipDstIPAddr "${ETH_NETWORK_IP}" \
  -skipGwIPAddr "${WAVE_DATA_IP_ADDR}"

pkill -9 -f castinator || true

nohup ./castinator \
  -leftIntfcName "${ETH_INTERFACE}" \
  -leftUDPAddr "239.192.137.1":${ETH_PORT} \
  -rightIntfcName "${WAVE_DATA_INTERFACE}" \
  -rightUDPAddr "${WAVE_DATA_BROADCAST_IP}":${WAVE_DATA_PORT} \
  >/tmp/castinator.log &

popd
