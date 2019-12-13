package main

import (
	"bootstrap_mk5/pkg/generate"
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	baseIPAddr := flag.String("baseIPAddr", "", "IP address to use as a base (e.g. 192.168.234.1)")
	identifierOctet := flag.Int("identifierOctet", 0, "Octet of base IP address to replace with the identifier (1 - 4 inclusive)")
	specificIdentifier := flag.Int("specificIdentifier", -1, "Specific identifier to use (cannot be used with -interfaceName)")
	interfaceName := flag.String("interfaceName", "", "Name of interface to use as seed for MAC generation (cannot be used with -specificIdentifier)")

	flag.Parse()

	var identifier uint8

	if *baseIPAddr == "" || strings.Count(*baseIPAddr, ".") != 3 {
		log.Fatal("error: -baseIPAddr flag empty or invalid")
	}

	if *identifierOctet < 1 || *identifierOctet > 4 {
		log.Fatal("error: -identifierOctet flag empty or invalid (must be 1 - 4 inclusive)")
	}

	if *specificIdentifier == -1 && *interfaceName == "" {
		log.Fatal("error: need -specificIdentifier or -interfaceName")
	}

	if *specificIdentifier != -1 {
		if *specificIdentifier < 0 || *specificIdentifier > 255 {
			log.Fatal("error: -specificIdentifier flag invalid (must be 1 - 255 inclusive)")
		}

		if *interfaceName != "" {
			log.Fatal("error: -interfaceName cannot be if -specificIdentifier is set")
		}

		identifier = uint8(*specificIdentifier)
	} else {
		if *interfaceName == "" {
			log.Fatal("error: -interfaceName flag empty")
		}

		seedMACAddr, err := generate.GetMACAddrOfInterface(*interfaceName)
		if err != nil {
			log.Fatalf("failed to get MAC of %v because: %v", *interfaceName, err)
		}

		identifier, err = generate.GetIdentifierFromMACAddr(seedMACAddr)
		if err != nil {
			log.Fatalf("failed to identifier from MAC of %v because: %v", seedMACAddr, err)
		}
	}

	fmt.Printf("%v\n", generate.BuildIPAddr(*baseIPAddr, *identifierOctet, identifier))
}
