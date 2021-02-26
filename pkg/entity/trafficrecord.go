package entity

import "encoding/json"

// TrafficRecord the record of the traffic detected
type TrafficRecord struct {
	Timestamp int64  `json:"timestamp"` // Timestamp when the traffic record is generated
	ProbeIP   string `json:"probeIP"`   // ProbeIP where is probe is collecting traffic data
	SrcIP     string `json:"srcIP"`     // SrcIP source IP of the traffic
	DstIP     string `json:"dstIP"`     // DstIP destination IP of the traffic
	Size      uint64 `json:"size"`      // Size size of the traffic
}

// ToJSONString convert the TrafficRecord to JSON string if not error
func (t TrafficRecord) ToJSONString() (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(b), err
}
