package main

import (
	"bootstrap_mk5/pkg/generate"
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	interfaceName := flag.String("interfaceName", "", "name of interface to use as seed for MAC generation")
	baseIPAddr := flag.String("baseIPAddr", "", "IP address to use as a base (e.g. 192.168.234.1)")
	identifierOctet := flag.Int("identifierOctet", 0, "octet of base IP address to replace with the identifier (1 - 4 inclusive)")

	flag.Parse()

	if *interfaceName == "" {
		log.Fatal("error: -interfaceName flag empty")
	}

	if *baseIPAddr == "" || strings.Count(*baseIPAddr, ".") != 3 {
		log.Fatal("error: -baseIPAddr flag empty or invalid")
	}

	if *identifierOctet < 1 || *identifierOctet > 4 {
		log.Fatal("error: -identifierOctet flag empty or invalid (must be 1 - 4 inclusive)")
	}

	seedMACAddr, err := generate.GetMACAddrOfInterface(*interfaceName)
	if err != nil {
		log.Fatalf("failed to get MAC of %v because: %v", *interfaceName, err)
	}

	identifier, err := generate.GetIdentifierFromMACAddr(seedMACAddr)
	if err != nil {
		log.Fatalf("failed to identifier from MAC of %v because: %v", seedMACAddr, err)
	}

	fmt.Printf("%v\n", generate.BuildIPAddr(*baseIPAddr, *identifierOctet, identifier))
}
