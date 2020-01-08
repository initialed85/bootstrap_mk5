#!/usr/bin/env bash

set -x
set -e

pushd "$(pwd)"

catch() {
  echo "error: caught a failure; popping original directory and exiting"

  popd

  exit 1
}

trap 'catch' ERR

cd ./deploy

echo "cleaning..."
rm -fr ./hosts 2>&1 || true
echo ""

echo "building hosts file"
cp -frv ./hosts-base ./hosts
./find_mk5s >>./hosts
echo ""

echo "hosts are"
cat ./hosts
echo ""

echo "stamping..."
date >../dist/deploy_date.txt
echo ""

echo "deploying"
ansible-playbook -i ./hosts ./deploy.yml
echo ""

if [[ "${1}" == "reboot" ]]; then
  echo "rebooting"
  ansible-playbook -i ./hosts ./reboot.yml
  echo ""
fi

popd
