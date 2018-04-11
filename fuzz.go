// +build gofuzz

package arp

func Fuzz(data []byte) int {
	arp := new(ARPPacket)
	err := arp.UnmarshalARP(data)
	if err != nil {
		return 0
	}
	return 1
}
