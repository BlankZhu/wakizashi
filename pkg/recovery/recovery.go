// Package recovery describe the recovery mechanism used by wakizashi for error posting handling.
// Example:
//  r := recovery.Get()
//  err := r.Init(recovFile, posFile, recovLimit, cacheSize, postFunc)
//  if err != nil {
//  ...
//  }
//  go func() {
//  	for {
// 			select {
// 				case <-time.After(60 * time.Second):
// 				r.RepostRecord()
// 			}
// 		}
//  }()
//  ...
//  r.Add2Recovery(record)
package recovery

import (
	wconst "BlankZhu/wakizashi/pkg/const"
	"BlankZhu/wakizashi/pkg/entity"
	wlogger "BlankZhu/wakizashi/pkg/log"
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

// RecoverPostFunc define the actual re-post behaviour of the recovery
type RecoverPostFunc func(record *entity.TrafficRecord) error

// DefaultPosLimit the maximun size of pos, implying when to rotate to the next file
const DefaultPosLimit = (1 << 20) * 10 // currently, 10M
// DefaultPosFileName filename of position file, saving the offset of recovery file
const DefaultPosFileName = "pos"

// DefaultRecoveryFileName the filename of recovery file, holding all those traffic record failed to post
const DefaultRecoveryFileName = "data"

// FailureRecover help wakizashi
var failureRecovery *recovery
var newOnce sync.Once
var initOnce sync.Once

type recovery struct {
	recoveryFilePath    string     // recording traffic data
	positionFilePath    string     // recording the current offset of recovery file
	recoveryLengthLimit int64      // size limit of recovery data recoding file
	recoveryMutex       sync.Mutex // mutex of the recovery file
	positionMutex       sync.Mutex // mutex of the position file

	cache      []*entity.TrafficRecord // cache to hold the incoming traffic record
	cacheSize  int
	cacheMutex sync.Mutex
	recordChan chan *entity.TrafficRecord // cache channel for incoming writing
	postFunc   RecoverPostFunc            // function used for posting to data storage backend
}

// Get return the singleton recovery
func Get() *recovery {
	newOnce.Do(func() {
		failureRecovery = &recovery{}
	})
	return failureRecovery
}

// Init initialize the recovery and launch it, and only the first call will do the job
func (r *recovery) Init(recoveryFilePath, positionFilePath string, recoveryLengthLimit int64, cacheSize int, pfunc RecoverPostFunc) error {
	var retErr error = nil

	initOnce.Do(func() {
		r.recordChan = make(chan *entity.TrafficRecord, cacheSize)

		r.recoveryMutex.Lock()
		defer r.recoveryMutex.Unlock()
		r.positionMutex.Lock()
		defer r.positionMutex.Unlock()

		r.recoveryFilePath = recoveryFilePath
		r.positionFilePath = positionFilePath
		r.recoveryLengthLimit = recoveryLengthLimit
		r.cacheSize = cacheSize
		r.cache = make([]*entity.TrafficRecord, 0, r.cacheSize)
		r.postFunc = pfunc

		rf, err := os.OpenFile(r.recoveryFilePath, os.O_CREATE, 0666)
		if err != nil {
			retErr = err
		}
		rf.Close()

		pf, err := os.OpenFile(r.positionFilePath, os.O_CREATE, 0666)
		if err != nil {
			retErr = err
		}
		pf.Close()

		go func() {
			for {
				select {
				case record := <-r.recordChan:
					r.cacheMutex.Lock()
					r.cache = append(r.cache, record)
					if len(r.cache) > r.cacheSize {
						r.FlushRecords()
					}
					r.cacheMutex.Unlock()
				case <-time.After(60 * time.Second * wconst.RecoveryFlushInterval / 6):
					r.cacheMutex.Lock()
					r.FlushRecords()
					r.cacheMutex.Unlock()
				}
			}
		}()
	})
	return retErr
}

// RepostRecord read recovery file, then do the recovery reposting
func (r *recovery) RepostRecord() {
	wlog := wlogger.Get()
	r.recoveryMutex.Lock()
	defer r.recoveryMutex.Unlock()

	// get position
	pos := r.getPosition()

	// set seeking-position
	rf, err := os.OpenFile(r.recoveryFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	defer rf.Close()
	if err != nil {
		wlog.Errorf("Failed to open recovery file, detail: %s", err.Error())
		return
	}
	_, err = rf.Seek(pos, 0)
	if err != nil {
		wlog.Errorf("Failed to seek recovery file, detail: %s", err.Error())
		return
	}
	// read history records
	failedRecords := make([]entity.TrafficRecord, 0, r.cacheSize)
	lineReader := bufio.NewScanner(rf)
	for lineReader.Scan() {
		var record entity.TrafficRecord
		line := lineReader.Text()
		pos += int64(len(line) + 1)

		err := json.Unmarshal([]byte(line), &record)
		if err != nil {
			continue
		}

		err = r.post(&record)
		if err != nil {
			wlog.Warningf("Failed to post record, detail: %d", err.Error())
			failedRecords = append(failedRecords, record)
		}
	}

	// check if current pos is out of PosLimit
	if pos > r.recoveryLengthLimit {
		pos = r.clearRecord(pos)
	}

	for _, record := range failedRecords {
		str, err := record.ToJSONString()
		if err != nil {
			wlog.Warningf("Failed to parse to JSON string: %v", record)
			continue
		}
		_, err = rf.WriteString(str + "\n")
		if err != nil {
			wlog.Errorf("Failed to write recovery file, with reocrd: %s", str)
			continue
		}
	}

	err = r.writePosition(pos)
	if err != nil {
		wlog.Errorf("Failed to write position file, detail: %s", err.Error())
	}
	return
}

// FlushRecords flush all the traffic record in cache to recovery file
func (r *recovery) FlushRecords() {
	if len(r.cache) == 0 {
		return
	}

	wlog := wlogger.Get()

	r.recoveryMutex.Lock()
	defer r.recoveryMutex.Unlock()

	f, err := os.OpenFile(r.recoveryFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		wlog.Errorf("Failed to open recovery file, detail: %s", err.Error())
		return
	}
	w := bufio.NewWriter(f)

	for _, record := range r.cache {
		str, err := record.ToJSONString()
		if err != nil {
			wlog.Warningf("Failed to parse to JSON string: %v", record)
			continue
		}
		_, err = w.WriteString(str + "\n")
		if err != nil {
			wlog.Warningf("Failed to write string to recovery file, detail: %s", err.Error())
		}
	}

	r.cache = nil
	r.cache = make([]*entity.TrafficRecord, 0, r.cacheSize)

	err = w.Flush()
	if err != nil {
		wlog.Warningf("Flush error, detail: %s", err.Error())
	}
}

// Add2Recovery is used to add a post-failed record to recovery by the caller
func (r *recovery) Add2Recovery(record *entity.TrafficRecord) {
	r.recordChan <- record
}

func (r *recovery) post(record *entity.TrafficRecord) error {
	return r.postFunc(record)
}

func (r *recovery) clearRecord(pos int64) int64 {
	wlog := wlogger.Get()

	err := os.Remove(r.recoveryFilePath)
	if err != nil {
		wlog.Errorf("Failed to remove recovery file, detail: %s", err.Error())
		return pos
	}
	rf, err := os.OpenFile(r.recoveryFilePath, os.O_CREATE, 0666)
	if err != nil {
		wlog.Errorf("Failed to re-create recovery file, detail: %s", err.Error())
		return 0
	}
	rf.Close()

	err = os.Remove(r.positionFilePath)
	if err != nil {
		wlog.Errorf("Failed to remove position file, detail: %s", err.Error())
		return 0
	}
	pf, err := os.OpenFile(r.recoveryFilePath, os.O_CREATE, 0666)
	if err != nil {
		wlog.Errorf("Failed to re-create position file, detail: %s", err.Error())
		return 0
	}
	pf.Close()

	return 0
}

func (r *recovery) getPosition() int64 {
	wlog := wlogger.Get()
	r.positionMutex.Lock()
	defer r.positionMutex.Unlock()
	var ret int64 = 0

	posData, err := ioutil.ReadFile(r.positionFilePath)
	if err != nil {
		wlog.Warningf("Failed to open position file, detail: %s", err.Error())
		return ret
	}

	if len(posData) != 0 {
		ret, err = strconv.ParseInt(string(posData), 10, 64)
		if err != nil {
			wlog.Warningf("Cannot parse position data, detail: %s", err.Error())
		}
	}
	return ret
}

func (r *recovery) writePosition(pos int64) error {
	r.positionMutex.Lock()
	defer r.positionMutex.Unlock()

	buf := bytes.NewBufferString(strconv.FormatInt(pos, 10))
	pf, err := os.OpenFile(r.positionFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	_, err = pf.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
