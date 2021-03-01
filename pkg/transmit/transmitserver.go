package transmit

import (
	"io"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/peer"
)

// CenterServer implements UnimplementedTransmitServer
type CenterServer struct {
	UnimplementedTransmitServer
	IPSet map[string]struct{}
}

// HandleRequest handles the grpc requests from probe
func (cs *CenterServer) HandleRequest(stream Transmit_TransmitServer) error {
	if peer, ok := peer.FromContext(stream.Context()); ok {
		logrus.Infof("receiving traffic data transmit request from: %s", peer.Addr.String())
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&TransmitReply{
				Res: true,
				Detail: "connection close",
			})
		}
		if err != nil {
			logrus.Errorf("transmit request error, detail: %s", err)
			return err
		}

		// filters those traffic between probe and center
		_, isFromCenter := cs.IPSet[req.SrcIP]
		_, isToCenter := cs.IPSet[req.DstIP]
		if isFromCenter || isToCenter {
			continue
		}

		cs.handleTransmitRequest(req)
	}
}

func (cs *CenterServer) handleTransmitRequest(req *TransmitRequest) {
	// recover := recovery.Get()
	// ts := time.Unix(int64(req.Timestamp), 0)

	// TODO: add databack end
}
