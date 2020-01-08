package gps

import (
	"encoding/json"
	"fmt"
	"github.com/initialed85/castinator/pkg/interfaces"
	"github.com/initialed85/castinator/pkg/sender"
	"github.com/initialed85/drive_test/pkg/gps_dumper"
	"log"
	"net"
)

func getJSON(object interface{}) ([]byte, error) {
	objectJSON, err := json.Marshal(object)
	if err != nil {
		return []byte{}, err
	}

	log.Println(string(objectJSON) + "\n")

	return objectJSON, nil
}

type Forwarder struct {
	dumper gps_dumper.Dumper
	sender *net.UDPConn
}

func New(gpsHost string, gpsPort int, sendHost string, sendPort int) (Forwarder, error) {
	f := Forwarder{}

	var err error

	f.dumper, err = gps_dumper.New(gpsHost, gpsPort, f.callback)
	if err != nil {
		return Forwarder{}, err
	}

	addr, err := interfaces.GetAddress(fmt.Sprintf("%v:%v", sendHost, sendPort))
	if err != nil {
		return Forwarder{}, err
	}

	f.sender, err = sender.GetSender(addr, nil)
	if err != nil {
		return Forwarder{}, err
	}

	return f, nil
}

func (f *Forwarder) callback(output gps_dumper.Output) error {
	b, err := getJSON(output)
	if err != nil {
		return err
	}

	_, err = f.sender.Write(b)

	return err
}

func (f *Forwarder) Watch() (done chan bool) {
	return f.dumper.Watch()
}
