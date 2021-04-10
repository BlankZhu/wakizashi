package constant

const (
	// Recovery Related constants

	// CenterDefaultConfigPath default config path for wakizashi's center
	CenterDefaultConfigPath = "./center-config.yaml"
	// ProbeDefaultConfigPath default config path for wakiazashi's probe
	ProbeDefaultConfigPath = "./probe-config.yaml"
	// RecoveryDefaultFileName default recovery file name
	RecoveryDefaultFileName = "rcv_data"
	// RecoveryDefaultPosName default position file name
	RecoveryDefaultPosName = "pos_data"
	// RecoveryDefaultPosLimit position limitation for recovery data file
	RecoveryDefaultPosLimit = (1 << 20) * 10 //10M
	// RecoveryDefaultCacheSize cache size for recovery
	RecoveryDefaultCacheSize = 128

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

	// BackendInfluxDB backend name of the influxdb
	BackendInfluxDB = "influxdb"
	// BackendMongoDB backend name of the mongodb
	BackendMongoDB = "mongodb"
	// BackendRedis backend name of the reids
	BackendRedis = "redis"

	// WakizashiDefaultDatabase default database name
	WakizashiDefaultDatabase = "wakizashi"
	// WakizashiDefaultTable default table name
	WakizashiDefaultTable = "wakizashi"
)
