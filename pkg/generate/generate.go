package generate

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const octets4And5 = "62:91"

func GetMACAddrOfInterface(interfaceName string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	interfaceName = strings.TrimSpace(strings.ToLower(interfaceName))

	compareInterfaceNames := make([]string, 0)
	for _, i := range interfaces {
		compareInterfaceName := strings.TrimSpace(strings.ToLower(i.Name))

		compareInterfaceNames = append(compareInterfaceNames, compareInterfaceName)

		if compareInterfaceName != interfaceName {
			continue
		}

		macAddr := i.HardwareAddr.String()
		if strings.Count(macAddr, ":") != 5 {
			return "", fmt.Errorf("found interface %v but MAC of '%v' is not valid", interfaceName, macAddr)
		}

		return i.HardwareAddr.String(), nil
	}

	return "", fmt.Errorf("failed to find interface %v; options are %v", interfaceName, compareInterfaceNames)
}

func GetIdentifierFromMACAddr(MACAddr string) (uint8, error) {
	identifier, err := strconv.ParseUint(strings.Split(MACAddr, ":")[5], 16, 8)
	if err != nil {
		return 0, err
	}

	return uint8(identifier), nil
}

func BuildMACAddr(seedMACAddr string, identifier uint8) string {
	vendorPart := strings.Join(strings.Split(seedMACAddr, ":")[0:3], ":")

	return strings.TrimSpace(strings.Join([]string{vendorPart, strings.Trim(octets4And5, ":"), fmt.Sprintf("%02x", identifier)}, ":"))
}

func BuildIPAddr(baseIPAddr string, identifierOctet int, identifier uint8) string {
	parts := strings.Split(baseIPAddr, ".")
	leftParts := parts[0 : identifierOctet-1]
	rightParts := parts[identifierOctet:]

	finalParts := append(leftParts, []string{fmt.Sprintf("%v", identifier)}...)

	return strings.Join(append(finalParts, rightParts...), ".")
}
