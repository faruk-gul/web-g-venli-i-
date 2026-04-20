package scanner

import (
	"fmt"
	"net"
	"time"
)

func hostPort(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func netDialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, address, timeout)
}
