package util

import (
	"net"

	"golang.org/x/sys/unix"
)

// UDPAddr -> RawSockaddr conversion
func NetAddrToSockaddr(addr *net.UDPAddr) unix.Sockaddr {
	if ip := addr.IP.To4(); ip != nil {
		return &unix.SockaddrInet4{
			Port: addr.Port,
			Addr: [4]byte(ip),
		}
	}

	if ip := addr.IP.To16(); ip != nil {
		return &unix.SockaddrInet6{
			Addr: [16]byte(ip),
			Port: addr.Port,
		}
	}

	return nil
}

// RawSockaddr -> UDPAddr conversion
func SockaddrToUDPAddr(sockAddr unix.Sockaddr) *net.UDPAddr {
	switch typedSocketAddr := sockAddr.(type) {

	case *unix.SockaddrInet4:
		return &net.UDPAddr{
			IP:   typedSocketAddr.Addr[:],
			Port: typedSocketAddr.Port,
			Zone: "",
		}

	case *unix.SockaddrInet6:
		if !IsValidIPv6(string(typedSocketAddr.Addr[:])) {
			return nil
		}

		return &net.UDPAddr{
			IP:   typedSocketAddr.Addr[:],
			Port: typedSocketAddr.Port,
			Zone: GetZone(typedSocketAddr.ZoneId),
		}
	default:
		return nil
	}
}

func IsValidIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil || parsedIP.To4() != nil {
		return false
	}
	return true
}

func GetZone(zoneID uint32) (zone string) {
	ifi, err := net.InterfaceByIndex(int(zoneID))
	if err == nil {
		zone = ifi.Name
	}
	return
}
