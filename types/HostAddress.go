package types

// Reference: https://www.ietf.org/rfc/rfc4120.txt
// Section: 5.2.5

import (
	"bytes"
	"fmt"
	"github.com/jcmturner/asn1"
	"net"
)

/*
HostAddress and HostAddresses

HostAddress     ::= SEQUENCE  {
	addr-type       [0] Int32,
	address         [1] OCTET STRING
}

-- NOTE: HostAddresses is always used as an OPTIONAL field and
-- should not be empty.
HostAddresses   -- NOTE: subtly different from rfc1510,
		-- but has a value mapping and encodes the same
	::= SEQUENCE OF HostAddress

The host address encodings consist of two fields:

addr-type
	This field specifies the type of address that follows.  Pre-
	defined values for this field are specified in Section 7.5.3.

address
	This field encodes a single address of type addr-type.
*/

const (
	AddrType_IPv4            = 2
	AddrType_Directional     = 3
	AddrType_ChaosNet        = 5
	AddrType_XNS             = 6
	AddrType_ISO             = 7
	AddrType_DECNET_Phase_IV = 12
	AddrType_AppleTalk_DDP   = 16
	AddrType_NetBios         = 20
	AddrType_IPv6            = 24
)

type HostAddresses []HostAddress

type HostAddress struct {
	AddrType int    `asn1:"explicit,tag:0"`
	Address  []byte `asn1:"explicit,tag:1"`
}

func GetHostAddress(s string) (HostAddress, error) {
	var h HostAddress
	cAddr, _, err := net.SplitHostPort(s)
	if err != nil {
		return h, fmt.Errorf("Invalid format of client address: %v", err)
	}
	ip := net.ParseIP(cAddr)
	hb, err := ip.MarshalText()
	if err != nil {
		return h, fmt.Errorf("Could not marshal client's address into bytes: %v", err)
	}
	var ht int
	if ip.To4() != nil {
		ht = AddrType_IPv4
	} else if ip.To16() != nil {
		ht = AddrType_IPv6
	} else {
		return h, fmt.Errorf("Could not determine client's address types: %v", err)
	}
	h = HostAddress{
		AddrType: ht,
		Address:  hb,
	}
	return h, nil
}

func (h *HostAddress) GetAddress() (string, error) {
	var b []byte
	_, err := asn1.Unmarshal(h.Address, &b)
	return string(b), err
}

func HostAddressesEqual(h, a []HostAddress) bool {
	if len(h) != len(a) {
		return false
	}
	for _, e := range a {
		var found bool
		found = false
		for _, i := range h {
			if e.Equal(i) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func HostAddressesContains(h []HostAddress, a HostAddress) bool {
	for _, e := range h {
		if e.Equal(a) {
			return true
		}
	}
	return false
}

func (h *HostAddress) Equal(a HostAddress) bool {
	if h.AddrType != a.AddrType {
		return false
	}
	return bytes.Equal(h.Address, a.Address)
}

func (h *HostAddresses) Contains(a HostAddress) bool {
	for _, e := range *h {
		if e.Equal(a) {
			return true
		}
	}
	return false
}

func (h *HostAddresses) Equal(a []HostAddress) bool {
	if len(*h) != len(a) {
		return false
	}
	for _, e := range a {
		if !h.Contains(e) {
			return false
		}
	}
	return true
}
