package main

import (
	"bootstrap_mk5/pkg/generate"
	"flag"
	"fmt"
	"log"
)

func main() {
	interfaceName := flag.String("interfaceName", "", "name of interface to use as seed for MAC generation")

	flag.Parse()

	if *interfaceName == "" {
		log.Fatal("error: -interfaceName flag empty")
	}

	seedMACAddr, err := generate.GetMACAddrOfInterface(*interfaceName)
	if err != nil {
		log.Fatalf("failed to get MAC of %v because: %v", *interfaceName, err)
	}

	identifier, err := generate.GetIdentifierFromMACAddr(seedMACAddr)
	if err != nil {
		log.Fatalf("failed to identifier from MAC of %v because: %v", seedMACAddr, err)
	}

	fmt.Printf("%v\n", generate.BuildMACAddr(seedMACAddr, identifier))
}
