package portcheck

import (
	"fmt"
	"net"
)

func IsPortFree(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	
	ln.Close()
	return true
}

func FindFreePort(port int) int {
	port += 1
	for IsPortFree(port) != true {
		port += 1
	}

	return port
}
