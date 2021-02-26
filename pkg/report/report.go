package report

import (
	"BlankZhu/wakizashi/pkg/constant"
	"BlankZhu/wakizashi/pkg/entity"
	"BlankZhu/wakizashi/pkg/transmit"
	"BlankZhu/wakizashi/pkg/types"
	"BlankZhu/wakizashi/pkg/util"
	"bufio"
	"context"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Reporter get send the traffic data to the data backend
type Reporter struct {
	AutoClear   bool            // clear the processed dump file or not
	DumpDir     string          // directory to save dump file
	FileCh      <-chan string   // channel used to communicate between reporter & dumper
	Ifaces      []net.Interface // on which network interface the reporter is working
	RepAddr     string          // address of the center
	RepRetry    int             // retry count to transmit data to center
	RepInterval time.Duration   // retry interval
	repCache    types.ReporterCache
	transCli    transmit.TransmitClient
}

// Init initialize the traffic reporter
func (r *Reporter) Init() {
	r.repCache.Init()
}

// Start starts the reporter process
func (r *Reporter) Start() {
	// todo
}

func (r *Reporter) handleCapturedFile() {
	for {
		select {
		case filename := <-r.FileCh:
			logrus.Infof("processing captured traffic recording file: %s", filename)
			records := r.analyzeCapturedFile(filename)
			r.loadCache(records)

			if r.AutoClear {
				err := os.Remove(path.Join(r.DumpDir, filename))
				if err != nil {
					logrus.Warnf("failed to remove dump file %s, detail: %s", filename, err)
				}
			}
		}
	}
}

func (r *Reporter) analyzeCapturedFile(filename string) []*entity.TrafficRecord {
	var ret []*entity.TrafficRecord
	logrus.Debugf("analyzing file: %s", filename)

	filepath := path.Join(r.DumpDir, filename)
	f, err := os.Open(filepath)
	if err != nil {
		logrus.Errorf("failed to open file %s, detail: %s", filepath, err)
		return ret
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		elems := strings.Split(line, " ")
		if len(elems) != 3 {
			logrus.Warnf("get invalid line in dump file %s, with: %s", filepath, line)
			continue
		}

		srcIP := elems[0]
		dstIP := elems[1]
		sz, err := strconv.ParseUint(elems[2], 10, 64)
		if err != nil {
			logrus.Warnf("failed to extract size in dump file %s, detail; %s", filepath, line)
			continue
		}

		var probeIP string
		ips := util.GetIPSetFromNetworkInterfaces(r.Ifaces)
		if _, b := ips[srcIP]; b {
			probeIP = srcIP
		} else if _, b := ips[dstIP]; b {
			probeIP = dstIP
		} else {
			continue
		}

		r := entity.TrafficRecord{
			SrcIP:   srcIP,
			DstIP:   dstIP,
			Size:    sz,
			ProbeIP: probeIP,
		}
		ret = append(ret, &r)
	}
	if sc.Err() != nil {
		logrus.Errorf("failed to scan dump file %s, detail: %s", filepath, sc.Err())
	}

	return ret
}

func (r *Reporter) loadCache(records []*entity.TrafficRecord) {
	ts := time.Now().UTC().Unix()
	r.repCache.Lock()
	defer r.repCache.Unlock()

	var kb strings.Builder
	for _, record := range records {
		kb.WriteString(record.SrcIP)
		kb.WriteString("_")
		kb.WriteString(record.DstIP)
		key := kb.String()
		kb.Reset()
		_, b := r.repCache.Data[key]
		if b {
			r.repCache.Data[key].Size = r.repCache.Data[key].Size + record.Size
		} else {
			r.repCache.Data[key] = record
			r.repCache.Data[key].Timestamp = ts
		}
	}
}

func (r *Reporter) consume() {
	conn, err := grpc.Dial(r.RepAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second*constant.ProbeTransmitTimeout))
	if err != nil {
		logrus.Errorf("failed to connect to center, detail: %s", err)
		return
	}
	defer conn.Close()
	r.transCli = transmit.NewTransmitClient(conn)
	stream, err := r.transCli.Transmit(context.TODO())
	if err != nil {
		logrus.Errorf("failed to create transmit stream, detail: %s", err)
		return
	}

	ticker := time.NewTicker(r.RepInterval)
	for {
		go func(mtx *sync.RWMutex) {
			mtx.Lock()
			defer mtx.Unlock()
			for k, v := range r.repCache.Data {
				err := stream.Send(&transmit.TransmitRequest{
					Timestamp: uint64(v.Timestamp), // FIXME: potential casting error here
					SrcIP:     v.SrcIP,
					DstIP:     v.DstIP,
					Size:      v.Size,
					PodIP:     v.ProbeIP,
				})
				if err != nil {
					logrus.Warnf("failed to transmit record to center, detail: %s", err)
					continue
				}
				delete(r.repCache.Data, k)
			}
		}(&r.repCache.RWMutex)
		<-ticker.C
	}
}

func (r *Reporter) report() {
	failCnt := 0
	retryFactor := 2

	for {
		if r.RepRetry == 0 {
			r.consume()
			logrus.Warnf("reporter will try consuming cache after %d sec", retryFactor)
			time.Sleep(time.Second * time.Duration(retryFactor))
			continue
		}

		if failCnt < r.RepRetry {
			r.consume()
			logrus.Warnf("reporter will try consuming cache after %s sec", retryFactor)
			failCnt++
			retryFactor = retryFactor * retryFactor
		} else {
			logrus.Errorf("failed to establish connection to center")
			return
		}
	}
}
