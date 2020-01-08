package main

import (
	"flag"
	"github.com/initialed85/castinator/pkg/castinator"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func waitForCtrlC() {
	sig := make(chan os.Signal, 2)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	log.Printf("waiting for CTRL + C")

	<-sig

	log.Printf("CTRL + C caught")
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	leftIntfcName := flag.String("leftIntfcName", "", "left interface name")
	leftUDPListenAddr := flag.String("leftUDPListenAddr", "", "left UDP listen addr")
	leftUDPSendAddr := flag.String("leftUDPSendAddr", "", "left UDP send addr")
	rightIntfcName := flag.String("rightIntfcName", "", "right interface name")
	rightUDPListenAddr := flag.String("rightUDPListenAddr", "", "right UDP listen addr")
	rightUDPSendAddr := flag.String("rightUDPSendAddr", "", "right UDP send addr")

	flag.Parse()

	if *leftIntfcName == "" {
		log.Fatal("-leftIntfcName empty")
	}

	if *leftUDPListenAddr == "" {
		log.Fatal("-leftUDPListenAddr empty")
	}

	if *leftUDPSendAddr == "" {
		log.Fatal("-leftUDPSendAddr empty")
	}

	if *rightIntfcName == "" {
		log.Fatal("-rightIntfcName empty")
	}

	if *rightUDPListenAddr == "" {
		log.Fatal("-rightUDPListenAddr empty")
	}

	if *rightUDPSendAddr == "" {
		log.Fatal("-rightUDPSendAddr empty")
	}

	c, err := castinator.New(*leftIntfcName, *leftUDPListenAddr, *leftUDPSendAddr, *rightIntfcName, *rightUDPListenAddr, *rightUDPSendAddr)
	if err != nil {
		log.Fatal(err)
	}

	c.Start()

	waitForCtrlC()

	c.Stop()
}
