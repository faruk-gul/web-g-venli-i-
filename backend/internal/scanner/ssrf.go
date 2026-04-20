package scanner

import (
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strings"
)

var blockedCIDRs = []string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"169.254.0.0/16",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
}

func ValidateTarget(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid target: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("only http and https targets are allowed")
	}

	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("target host is required")
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("cannot resolve host: %w", err)
	}

	for _, ip := range ips {
		addr, parseErr := netip.ParseAddr(strings.TrimSpace(ip.String()))
		if parseErr != nil {
			continue
		}
		if isBlockedAddress(addr) {
			return fmt.Errorf("target resolves to private or restricted network address: %s", ip.String())
		}
	}

	return nil
}

func isBlockedAddress(ip netip.Addr) bool {
	for _, cidr := range blockedCIDRs {
		prefix, err := netip.ParsePrefix(cidr)
		if err == nil && prefix.Contains(ip) {
			return true
		}
	}
	return false
}

