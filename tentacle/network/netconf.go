package network

import (
	"fmt"
	"net"
	"tentacle/config"

	"config"
)

func getIpByDevice() {
	ifs, _ := net.Interfaces()
	for i, f := range ifs {
		if f.Name == config.GlobalConfig.NetDevice {
			
		}

		fmt.Println(f.Name, f.MTU, f.Index, i, f.Flags.String())
		addrs, _ := f.Addrs()
		for _, a := range addrs {
			fmt.Println("--> ", a.String())
		}
	}
}
