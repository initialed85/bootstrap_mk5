package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

const (
	IPv6LinkLocalAddrPrefix = "fe80::"
	CohdaIPv6AddrPrefix     = "fe80::6e5:"
)

type PossibleInterface struct {
	Intfc string
	Addrs []string
}

func getPossibleInterfaces() ([]PossibleInterface, error) {
	intfcs, err := net.Interfaces()
	if err != nil {
		return []PossibleInterface{}, nil
	}

	possibleInterfaces := make([]PossibleInterface, 0)
	for _, intfc := range intfcs {
		if intfc.Flags&net.FlagUp == 0 || intfc.Flags&net.FlagBroadcast == 0 {
			continue
		}

		addrs, err := intfc.Addrs()
		if err != nil {
			continue
		}

		possibleAddrs := make([]string, 0)
		for _, addr := range addrs {
			if !strings.HasPrefix(addr.String(), IPv6LinkLocalAddrPrefix) {
				continue
			}

			possibleAddrs = append(possibleAddrs, addr.String())
		}

		if len(possibleAddrs) == 0 {
			continue
		}

		possibleInterfaces = append(
			possibleInterfaces,
			PossibleInterface{
				Intfc: intfc.Name,
				Addrs: possibleAddrs,
			},
		)
	}

	return possibleInterfaces, nil
}

func runCommand(executable string, arguments ...string) (string, string, error) {
	cmd := exec.Command(
		executable,
		arguments...,
	)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func main() {
	possibleInterfaces, err := getPossibleInterfaces()
	if err != nil {
		log.Fatal(err)
	}

	addrs := make([]string, 0)
	for _, possibleInterface := range possibleInterfaces {
		stdout, _, err := runCommand(
			"ping6",
			"-c",
			"4",
			fmt.Sprintf("ff02::1%%%v", possibleInterface.Intfc),
		)
		if err != nil {
			continue
		}

		for _, line := range strings.Split(stdout, "\n") {
			if !strings.Contains(line, CohdaIPv6AddrPrefix) {
				continue
			}

			addr := strings.TrimRight(strings.TrimRight(strings.Split(strings.Split(line, "from ")[1], " ")[0], ":"), ",")

			addrs = append(addrs, addr)
		}
	}

	tempAddrs := make(map[string]struct{})

	for _, addr := range addrs {
		tempAddrs[addr] = struct{}{}
	}

	deduplicatedAddrs := make([]string, 0)

	for addr := range tempAddrs {
		deduplicatedAddrs = append(deduplicatedAddrs, addr)
	}

	if len(deduplicatedAddrs) == 0 {
		fmt.Print("error: failed to find any MK5s")

		os.Exit(1)
	}

	for _, addr := range deduplicatedAddrs {
		fmt.Println(strings.TrimSpace(addr))
	}
}
