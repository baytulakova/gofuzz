package arp

import (
	"encoding/binary"
	"net"
	"io"
)

// Operation determines whether operation is a request or a reply
type Operation uint16

// ARPPacket represents an ARP packet
type ARPPacket struct {
	// HardwareType specifies an IANA-assigned hardware type, as described
	// in RFC 826.
	HardwareType uint16

	// ProtocolType specifies the internetwork protocol for which the ARP
	// request is intended.  Typically, this is the IPv4 EtherType.
	ProtocolType uint16

	// HardwareAddrLength specifies the length of the sender and target
	// hardware addresses included in a Packet.
	HardwareAddrLength uint8

	// IPLength specifies the length of the sender and target IPv4 addresses
	// included in a Packet.
	IPLength uint8

	// Operation specifies the ARP operation being performed, such as request
	// or reply.
	Operation Operation

	// SenderHardwareAddr specifies the hardware address of the sender of this
	// Packet.
	SenderHardwareAddr net.HardwareAddr

	// SenderIP specifies the IPv4 address of the sender of this Packet.
	SenderIP net.IP

	// TargetHardwareAddr specifies the hardware address of the target of this
	// Packet.
	TargetHardwareAddr net.HardwareAddr

	// TargetIP specifies the IPv4 address of the target of this Packet.
	TargetIP net.IP
}

func (p *ARPPacket) MarshalBinary() ([]byte, error) {
// 2 bytes: hardware type
// 2 bytes: protocol type
// 1 byte : hardware address length
// 1 byte : protocol length
// 2 bytes: operation
// N bytes: source hardware address
// N bytes: source protocol address
// N bytes: target hardware address
// N bytes: target protocol address

// Though an IPv4 address should always 4 bytes, go-fuzz
// very quickly created several crasher scenarios which
// indicated that these values can lie.
b := make([]byte, 2+2+1+1+2+(p.IPLength*2)+(p.HardwareAddrLength*2))

// Marshal fixed length data

binary.BigEndian.PutUint16(b[0:2], p.HardwareType)
binary.BigEndian.PutUint16(b[2:4], p.ProtocolType)

b[4] = p.HardwareAddrLength
b[5] = p.IPLength

binary.BigEndian.PutUint16(b[6:8], uint16(p.Operation))

// Marshal variable length data at correct offset using lengths
// defined in p

n := 8
hal := int(p.HardwareAddrLength)
pl := int(p.IPLength)

copy(b[n:n+hal], p.SenderHardwareAddr)
n += hal

copy(b[n:n+pl], p.SenderIP)
n += pl

copy(b[n:n+hal], p.TargetHardwareAddr)
n += hal

copy(b[n:n+pl], p.TargetIP)

return b, nil
}

func (p *ARPPacket) UnmarshalARP(b []byte) error {
	// Must have enough room to retrieve hardware address and IP lengths
	if len(b) < 8 {
		return io.ErrUnexpectedEOF
	}

	// Must have enough room to retrieve both hardware address and IP addresses
	if len(b) < 28 {
		return io.ErrUnexpectedEOF
	}

	// Sender hardware address
	p.SenderHardwareAddr = b[8:14]

	// Sender IP address
	p.SenderIP = b[14:18]

	// Target hardware address
	p.TargetHardwareAddr = b[18:24]

	// Target IP address
	p.TargetIP = b[24:28]

	return nil
}




