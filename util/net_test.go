package util_test

import (
	"math/rand"
	"net"
	"reflect"
	"testing"

	"github.com/tgragnato/magnetico/util"
	"golang.org/x/sys/unix"
)

func TestNetAddrToSockaddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		addr *net.UDPAddr
		want unix.Sockaddr
	}{
		{
			name: "IPv4",
			addr: &net.UDPAddr{
				IP:   net.ParseIP("192.0.2.1"),
				Port: 8080,
			},
			want: &unix.SockaddrInet4{
				Port: 8080,
				Addr: [4]byte{192, 0, 2, 1},
			},
		},
		{
			name: "IPv6",
			addr: &net.UDPAddr{
				IP:   net.ParseIP("2001:4860:4860::8888"),
				Port: 8080,
			},
			want: &unix.SockaddrInet6{
				Port: 8080,
				Addr: [16]byte{32, 1, 72, 96, 72, 96, 0, 0, 0, 0, 0, 0, 0, 0, 136, 136},
			},
		},
		{
			name: "Invalid IP",
			addr: &net.UDPAddr{
				IP:   net.ParseIP("invalid"),
				Port: 8080,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := util.NetAddrToSockaddr(test.addr)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("NetAddrToSockaddr() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestSockaddrToUDPAddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		sockAddr unix.Sockaddr
		want     *net.UDPAddr
	}{
		{
			name: "IPv4 with valid IP:Port",
			sockAddr: &unix.SockaddrInet4{
				Addr: [4]byte{192, 0, 2, 1},
				Port: 8080,
			},
			want: &net.UDPAddr{
				IP:   []byte{192, 0, 2, 1},
				Port: 8080,
			},
		},
		{
			name: "IPv6 with invalid IP",
			sockAddr: &unix.SockaddrInet6{
				Addr: [16]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
				Port: 9090,
			},
			want: nil,
		},
		{
			name: "IPv6 with invalid ZoneId",
			sockAddr: &unix.SockaddrInet6{
				Addr:   [16]byte{32, 1, 72, 96, 72, 96, 0, 0, 0, 0, 0, 0, 0, 0, 136, 136},
				Port:   8080,
				ZoneId: 12345,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := util.SockaddrToUDPAddr(test.sockAddr)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("SockaddrToUDPAddr() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsValidIPv6(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		ip   string
		want bool
	}{
		{
			ip:   "2001:4860:4860::8888",
			want: true,
		},
		{
			ip:   "192.0.2.1",
			want: false,
		},
		{
			ip:   "299.0.2.1",
			want: false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.ip, func(t *testing.T) {
			t.Parallel()
			if got := util.IsValidIPv6(test.ip); got != test.want {
				t.Errorf("IsValidIPv6() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_GetZone(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name   string
		zoneID uint32
		want   string
	}{
		{
			name:   "ZoneZero",
			zoneID: 0,
			want:   "",
		},
		{
			name:   "ZoneRandom",
			zoneID: uint32(rand.Intn(900000) + 100000),
			want:   "",
		},
	}

	// Avoid issues in sandboxes with limited network permissions
	loopbackIface := "lo"
	lo, err := net.InterfaceByName(loopbackIface)
	if err != nil {
		loopbackIface = "lo0"
		lo, err = net.InterfaceByName(loopbackIface)
	}
	if err == nil {
		tests = append(tests, struct {
			name   string
			zoneID uint32
			want   string
		}{
			name:   "ZoneLoopback",
			zoneID: uint32(lo.Index),
			want:   loopbackIface,
		})
	}

	for _, tt := range tests {
		test := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := util.GetZone(test.zoneID)
			if got != test.want {
				t.Errorf("getZone() = %v, want %v", got, test.want)
			}
		})
	}
}
