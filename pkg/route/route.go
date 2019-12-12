package route

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
)

func AddRoute(destinationIPAddr string, destinationPrefix int, gateway string) error {
	dst, err := netlink.ParseIPNet(fmt.Sprintf("%v/%v", destinationIPAddr, destinationPrefix))
	if err != nil {
		return err
	}

	gw := net.ParseIP(gateway)

	return netlink.RouteAdd(&netlink.Route{
		Dst: dst,
		Gw:  gw,
	})
}
