package constant

const (
	// Recovery Related constants

	// RecoveryFlushInterval flush interval of recovery in min
	RecoveryFlushInterval = 1
	// DefaultChanCap default capacity of the channel used by wakizashi
	DefaultChanCap = 256

	// ISO8601BasicFormat ISO-8601 basic time format, for time.format
	ISO8601BasicFormat = "20060102T150405Z"
	// ISO8601BasicFormatShort ISO-8601 basic time format in YYYYMMDD, for time.format
	ISO8601BasicFormatShort = "20060102"
	// ISO8601CapFileFormat specialized for wakizashi
	ISO8601CapFileFormat = "20060102150405"

	// AfpacketTargetSizeMB afpacket target size in MB
	AfpacketTargetSizeMB = 16

	// ProbeTransmitTimeout timeout for probe to transmit data to center, in sec
	ProbeTransmitTimeout = 60
)
