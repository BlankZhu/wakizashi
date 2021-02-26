package dump

import (
	"BlankZhu/wakizashi/pkg/constant"
	"BlankZhu/wakizashi/pkg/entity"
	"BlankZhu/wakizashi/pkg/util"
	"bufio"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/afpacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
)

// Dumper dumps the traffic of a specified network interface device
type Dumper struct {
	Iface          *net.Interface
	SnapLen        uint32
	RotateInterval time.Duration // rotate captured file
	RepAddr        string        // address of the center
	DumpDir        string
	FileCh         chan<- string
	rawDataCh      chan *entity.RawTrafficRecord
}

// Init initializes the dumper
func (d *Dumper) Init() {
	d.rawDataCh = make(chan *entity.RawTrafficRecord, constant.DefaultChanCap)
}

// Start starts the dumping process, generating the afpacket file
func (d *Dumper) Start() {
	go d.genFile()
	d.dump()
}

func (d *Dumper) dump() {
	probeIPs := util.GetIPSetFromNetworkInterface(d.Iface)
	centerIPs := make(map[string]struct{})
	splited := strings.Split(d.RepAddr, ":")
	if len(splited) != 2 {
		logrus.Warnf("might using invalid center address %s, try this format: [hostname]:[port]", d.RepAddr)
	}
	cips, err := net.LookupIP(splited[0])
	if err != nil {
		logrus.Warnf("failed to lookup IP for %s, detail: %s", splited[0], err)
	} else {
		for _, v := range cips {
			centerIPs[v.String()] = struct{}{}
		}
	}

	handle, err := d.getAfpacketHandle()
	if err != nil {
		logrus.Errorf("failed to create afpacket handle, detail: %s", err)
		return
	}
	defer handle.Close()

	var eth layers.Ethernet
	var ipv4 layers.IPv4
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ipv4)
	decoded := []gopacket.LayerType{}

	for {
		d, p, err := handle.ZeroCopyReadPacketData()
		if err != nil {
			logrus.Warnf("failed to zero copy afpacket data, detail: %s", err)
		}
		parser.DecodeLayers(d, &decoded)
		for _, lt := range decoded {
			switch lt {
			case layers.LayerTypeIPv4:
				_, pSrcIPCheck := probeIPs[ipv4.SrcIP.String()]
				_, pDstIPCheck := probeIPs[ipv4.SrcIP.String()]
				if !pSrcIPCheck && !pDstIPCheck {
					continue
				}
				_, cSrcIPCheck := centerIPs[ipv4.SrcIP.String()]
				_, cDstIPCheck := centerIPs[ipv4.SrcIP.String()]
				if cDstIPCheck || cSrcIPCheck {
					continue
				}

				rd := &entity.RawTrafficRecord{
					SrcIP: ipv4.SrcIP.String(),
					DstIP: ipv4.DstIP.String(),
					Size:  p.Length,
				}
				d.rawDataCh <- rd
			}
		}
	}
}

func (d *Dumper) genFile() {
	ticker := time.NewTicker(d.RotateInterval)
	defer ticker.Stop()

	w, fp, cf, err := d.newWriter()
	if err != nil {
		logrus.Errorf("failed to create dump file writer, detail: %s", err)
		return
	}

	for {
		select {
		case <-ticker.C:
			err = w.Flush()
			if err != nil {
				logrus.Errorf("failed to flush afpacket data to dump file %s, detail: %s", fp.Name(), err)
				continue
			}

			err = fp.Close()
			if err != nil {
				logrus.Errorf("failed to close dump file %s, detail: %s")
				continue
			}

			d.FileCh <- cf
			w, fp, cf, err = d.newWriter()
			if err != nil {
				logrus.Errorf("failed to re-create writer on file %s, detail: %s", fp.Name(), err)
				continue
			}
		case rd := <-d.rawDataCh:
			_, err := w.WriteString(rd.ToString())
			if err != nil {
				logrus.Errorf("failed to write string to writer, detail: %s", err)
			}
			w.WriteString("\n")
		}
	}
}

func (d *Dumper) newWriter() (*bufio.Writer, *os.File, string, error) {
	now := time.Now().UTC()
	capFileName := fmt.Sprintf("%s_%s.cap", d.Iface.Name, now.Format(constant.ISO8601CapFileFormat))
	capFilePath := path.Join(d.DumpDir, capFileName)
	f, err := os.Create(capFilePath)
	if err != nil {
		return nil, nil, "", err
	}
	w := bufio.NewWriter(f)
	return w, f, capFileName, nil
}

func (d *Dumper) getAfpacketHandle() (*afpacket.TPacket, error) {
	szFrame, szBlock, numBlocks, err := afpacketComputeSize(16, int(d.SnapLen), os.Getpagesize())
	if err != nil {
		return nil, err
	}

	ret, err := afpacket.NewTPacket(
		afpacket.OptInterface(d.Iface.Name),
		afpacket.OptFrameSize(szFrame),
		afpacket.OptBlockSize(szBlock),
		afpacket.OptNumBlocks(numBlocks),
		afpacket.OptAddVLANHeader(false),
		afpacket.OptPollTimeout(-time.Millisecond*10), // pcap.BlockForever
		afpacket.SocketRaw,
		afpacket.TPacketVersion3)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func afpacketComputeSize(targetSizeMb int, snaplen int, pageSize int) (
	frameSize int, blockSize int, numBlocks int, err error) {

	if snaplen < pageSize {
		frameSize = pageSize / (pageSize / snaplen)
	} else {
		frameSize = (snaplen/pageSize + 1) * pageSize
	}

	// 128 is the default from the gopacket library so just use that
	blockSize = frameSize * 128
	numBlocks = (targetSizeMb * 1024 * 1024) / blockSize

	if numBlocks == 0 {
		return 0, 0, 0, fmt.Errorf("interface buffer size is too small")
	}

	return frameSize, blockSize, numBlocks, nil
}
