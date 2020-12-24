package entity

import "fmt"

// RawTrafficRecord describe the original data collected by wakizashi's traffic probe
type RawTrafficRecord struct {
	SrcIP string
	DstIP string
	Size  uint64
}

// ToString convert the data to string
func (rtr *RawTrafficRecord) ToString() string {
	return fmt.Sprintf("%s %s %d", rtr.SrcIP, rtr.DstIP, rtr.Size)
}
