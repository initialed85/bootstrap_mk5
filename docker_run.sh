#!/usr/bin/env bash

set -x
set -e

image_name='cohda-mk5-80211p'

docker run --rm -it --net=host ${image_name} "${@}"
