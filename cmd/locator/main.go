package main

import (
	"flag"
	"github.com/initialed85/bootstrap_mk5/pkg/gps"
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

	gpsHost := flag.String("gpsHost", "localhost", "GPS host")
	gpsPort := flag.Int("gpsPort", 2947, "GPS port")
	sendHost := flag.String("sendHost", "", "UDP host to send to")
	sendPort := flag.Int("sendPort", 0, "UDP port to send to")

	flag.Parse()

	if *gpsHost == "" {
		log.Fatal("-gpsHost empty")
	}

	if *gpsPort <= 0 {
		log.Fatal("-gpsPort empty or less than zero")
	}

	if *sendHost == "" {
		log.Fatal("-sendHost empty")
	}

	if *sendPort <= 0 {
		log.Fatal("-sendPort empty or less than zero")
	}

	f, err := gps.New(*gpsHost, *gpsPort, *sendHost, *sendPort)
	if err != nil {
		log.Fatal(err)
	}

	go f.Watch()

	waitForCtrlC()
}
