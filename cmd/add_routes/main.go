package main

import (
	"bootstrap_mk5/pkg/generate"
	"bootstrap_mk5/pkg/route"
	"flag"
	"log"
	"strings"
)

func main() {
	baseDstIPAddr := flag.String("baseDstIPAddr", "", "IP address to use as a base for the destination (e.g. 192.168.2.0)")
	dstIdentifierOctet := flag.Int("dstIdentifierOctet", 0, "octet of base destination IP address to replace with the identifier (1 - 4 inclusive)")
	dstPrefix := flag.Int("dstPrefix", -1, "IP prefix to use with the destination (e.g. 24)")
	baseGwIPAddr := flag.String("baseGwIPAddr", "", "IP address to use as a base for the gateway (e.g. 192.168.234.0)")
	gwIdentifierOctet := flag.Int("gwIdentifierOctet", 0, "octet of base gateway IP address to replace with the identifier (1 - 4 inclusive)")
	startIdentifier := flag.Int("startIdentifier", -1, "identifier to start at (0 - 255 inclusive)")
	stopIdentifier := flag.Int("stopIdentifier", -1, "identifier to stop at (0 - 255 inclusive)")

	flag.Parse()

	if *baseDstIPAddr == "" || strings.Count(*baseDstIPAddr, ".") != 3 {
		log.Fatal("error: -baseDstIPAddr flag empty or invalid")
	}

	if *dstIdentifierOctet < 1 || *dstIdentifierOctet > 4 {
		log.Fatal("error: -dstIdentifierOctet flag empty or invalid (must be 1 - 4 inclusive)")
	}

	if *dstPrefix < 0 || *dstPrefix > 32 {
		log.Fatal("error: -dstPrefix flag empty or invalid (must be 0 - 32 inclusive)")
	}

	if *baseGwIPAddr == "" || strings.Count(*baseGwIPAddr, ".") != 3 {
		log.Fatal("error: -baseGwIPAddr flag empty or invalid")
	}

	if *gwIdentifierOctet < 1 || *gwIdentifierOctet > 4 {
		log.Fatal("error: -gwIdentifierOctet flag empty or invalid (must be 1 - 4 inclusive)")
	}

	if *startIdentifier < 0 || *startIdentifier > 255 {
		log.Fatal("error: -startIdentifier flag empty or invalid (must be 0 - 255 inclusive)")
	}

	if *stopIdentifier < 0 || *stopIdentifier > 255 || *stopIdentifier < *startIdentifier {
		log.Fatal("error: -stopIdentifier flag empty or invalid (must be 0 - 255 inclusive and greater than startIdentifier)")
	}

	for i := *startIdentifier; i <= *stopIdentifier; i++ {
		dstIPAddr := generate.BuildIPAddr(*baseDstIPAddr, *dstIdentifierOctet, uint8(i))

		gwIPAddr := generate.BuildIPAddr(*baseGwIPAddr, *gwIdentifierOctet, uint8(i))

		err := route.AddRoute(dstIPAddr, *dstPrefix, gwIPAddr)
		if err != nil {
			log.Printf("failed to add route %v/%v via %v because: %v", dstIPAddr, *dstPrefix, gwIPAddr, err)

			continue
		}

		log.Printf("added route %v/%v via %v", dstIPAddr, *dstPrefix, gwIPAddr)
	}
}
