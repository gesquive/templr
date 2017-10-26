package engine

import (
	"net"
	"strings"
	"text/template"
)

type HostInfo struct {
	Type string
	Addr string
	Name string
}

func NetFuncs() template.FuncMap {
	funcMap := template.FuncMap{
		"slice":           slice,
		"list":            list,
		"rpad":            rightPad,
		"ipfmt":           ipFmt,
		"lookupHosts":     lookupHosts,
		"lookupIPv4Host":  lookupIPv4Host,
		"lookupIPv6Host":  lookupIPv6Host,
		"isValidIPv4":     IsValidIPv4,
		"isValidIPv6":     IsValidIPv6,
		"isValidIPv4Addr": IsValidIPv4Addr,
		"isValidIPv6Addr": IsValidIPv6Addr,
		"isValidIPv4CIDR": IsValidIPv4CIDR,
		"isValidIPv6CIDR": IsValidIPv6CIDR,
	}

	return funcMap
}

func slice(args ...interface{}) []interface{} {
	return args
}

func list(arr []interface{}) string {
	items := []string{}
	for _, i := range arr {
		items = append(items, i.(string))
	}

	return strings.Join(items, ", ")
}

func ipFmt(addr string) string {
	if IsValidIPv4Addr(addr) {
		return rightPad(15, addr)
	} else if IsValidIPv6Addr(addr) {
		return rightPad(23, addr)
	} else if IsValidIPv4CIDR(addr) {
		return rightPad(19, addr)
	} else if IsValidIPv6CIDR(addr) {
		return rightPad(26, addr)
	}
	return addr
}

func rightPad(size int, str interface{}) string {
	strSize := len(str.(string))
	padSize := 1 + size - strSize
	if padSize <= 0 {
		return str.(string)
	}

	result := str.(string) + strings.Repeat(" ", int(padSize))
	return result[:size]
}

func lookupHosts(hosts []interface{}) []HostInfo {
	host4Info := []HostInfo{}
	host6Info := []HostInfo{}
	for _, host := range hosts {
		addrs, _ := lookupHost(host)
		for _, addr := range addrs {
			if IsValidIPv4(addr) {
				host4Info = append(host4Info, HostInfo{"4", addr, host.(string)})
			} else if IsValidIPv6(addr) {
				host6Info = append(host6Info, HostInfo{"6", addr, host.(string)})
			}
		}
	}
	return append(host4Info, host6Info...)
}

func lookupHost(host interface{}) ([]string, error) {
	addrs, err := net.LookupHost(host.(string))
	if err != nil {
		return []string{}, err
	}

	return addrs, nil
}

func lookupIPv4Host(host string) ([]string, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return []string{}, err
	}

	ipv4Addrs := []string{}
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip.To4() != nil {
			ipv4Addrs = append(ipv4Addrs, addr)
		}
	}
	return ipv4Addrs, err
}

func lookupIPv6Host(host string) ([]string, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return []string{}, err
	}

	ipv6Addrs := []string{}
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip.To4() == nil {
			ipv6Addrs = append(ipv6Addrs, addr)
		}
	}
	return ipv6Addrs, err
}

// IsValidIPv4 returns true if the given address is a valid IPv4 address or IPv4 CIDR range
func IsValidIPv4(addr string) bool {
	return IsValidIPv4Addr(addr) || IsValidIPv4CIDR(addr)
}

// IsValidIPv6 returns true if the given address is a valid IPv6 address or IPv6 CIDR range
func IsValidIPv6(addr string) bool {
	return IsValidIPv6Addr(addr) || IsValidIPv6CIDR(addr)
}

// IsValidIPv4Addr returns true if the given address is a valid IPv4 address
func IsValidIPv4Addr(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.To4() != nil
}

// IsValidIPv6Addr returns true if the given address is a valid IPv6 address
func IsValidIPv6Addr(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.To16() != nil
}

// IsValidIPv4CIDR returns true if the given address is a valid IPv4 CIDR range
func IsValidIPv4CIDR(addr string) bool {
	ip, _, err := net.ParseCIDR(addr)
	return err == nil && IsValidIPv4Addr(ip.String())
}

// IsValidIPv6CIDR returns true if the given address is a valid IPv6 CIDR range
func IsValidIPv6CIDR(addr string) bool {
	ip, _, err := net.ParseCIDR(addr)
	return err == nil && IsValidIPv6Addr(ip.String())
}
