package main

import (
	"flag"
	"fmt"
	"github.com/initialed85/castinator/pkg/handler"
	"github.com/initialed85/castinator/pkg/interfaces"
	"github.com/initialed85/castinator/pkg/listener"
	"github.com/initialed85/castinator/pkg/sender"
	"log"
	"sync"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	leftIntfcName := flag.String("leftIntfcName", "", "left interface name")
	leftUDPAddr := flag.String("leftUDPAddr", "", "left UDPv4/v6 address")
	rightIntfcName := flag.String("rightIntfcName", "", "right interface name")
	rightUDPAddr := flag.String("rightUDPAddr", "", "right UDPv4/v6 address")

	flag.Parse()

	if *leftIntfcName == "" {
		log.Fatal("-leftIntfcName empty")
	}

	if *leftUDPAddr == "" {
		log.Fatal("-leftUDPAddr empty")
	}

	if *rightIntfcName == "" {
		log.Fatal("-rightIntfcName empty")
	}

	if *rightUDPAddr == "" {
		log.Fatal("-rightUDPAddr empty")
	}

	leftAddr, leftIntfc, leftSrcAddr, err := interfaces.GetAddressesAndInterfaces(*leftIntfcName, *leftUDPAddr)
	if err != nil {
		log.Fatal(err)
	}

	leftListener, err := listener.GetListener(leftAddr, leftIntfc)
	if err != nil {
		log.Fatal(err)
	}

	leftSender, err := sender.GetSender(leftAddr, leftSrcAddr)
	if err != nil {
		log.Fatal(err)
	}

	rightAddr, rightIntfc, rightSrcAddr, err := interfaces.GetAddressesAndInterfaces(*rightIntfcName, *rightUDPAddr)
	if err != nil {
		log.Fatal(err)
	}

	rightListener, err := listener.GetListener(rightAddr, rightIntfc)
	if err != nil {
		log.Fatal(err)
	}

	rightSender, err := sender.GetSender(rightAddr, rightSrcAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("leftAddr = %+v\n", leftAddr)
	log.Printf("leftIntfc = %+v\n", leftIntfc)
	log.Printf("leftSrcAddr = %+v\n", leftSrcAddr)
	log.Printf("leftListener = %+v\n", leftListener)
	log.Printf("leftSender = %+v\n", leftSender)

	fmt.Println("")

	log.Printf("rightAddr = %+v\n", rightAddr)
	log.Printf("rightIntfc = %+v\n", rightIntfc)
	log.Printf("rightSrcAddr = %+v\n", rightSrcAddr)
	log.Printf("rightListener = %+v\n", rightListener)
	log.Printf("rightSender = %+v\n", rightSender)

	wg := sync.WaitGroup{}

	wg.Add(2)

	go handler.Handle(leftListener, rightSender, leftSender)
	go handler.Handle(rightListener, leftSender, rightSender)

	wg.Wait()
}
