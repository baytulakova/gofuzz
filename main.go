package main

import(
	"github.com/google/gofuzz"
	"net"
	"os"
	"fmt"
)

func main(){
	arp:=new(ARPPacket)
	f:=fuzz.New().NilChance(0.5)

	var a struct {
		Ht   uint16
		Pt   uint16
		Hal  uint8
		Ipl  uint8
		O    Operation
		Shwa net.HardwareAddr
		Sip  net.IP
		Thwa net.HardwareAddr
		Tip  net.IP
	}

	f.Fuzz(&a)

	arp.HardwareType = 2
	arp.ProtocolType = 0x0800
	arp.HardwareAddrLength = 6
	arp.IPLength = 4
	arp.Operation = 2
	arp.SenderHardwareAddr = a.Shwa
	arp.SenderIP = a.Sip
	arp.TargetHardwareAddr = a.Thwa
	arp.TargetIP = a.Tip

	b, _ := arp.MarshalBinary()

	file, e := os.Create("MarshalBinary")
	if e != nil {
		fmt.Errorf("Unable to create file: ", e)
		os.Exit(1)
	}
	defer file.Close()

	file.Write(b)

	err := arp.UnmarshalARP(b)
	if err != nil {
		fmt.Errorf("Failed to unmarshal ARP")
	}
}
