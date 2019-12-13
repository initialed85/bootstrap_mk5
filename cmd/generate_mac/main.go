package main

import (
	"flag"
	"fmt"
	"github.com/initialed85/bootstrap_mk5/pkg/generate"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	interfaceName := flag.String("interfaceName", "", "Name of interface to use as seed for MAC generation")

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
