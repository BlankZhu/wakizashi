/*
Package livenessprobe describes the liveness probe used by wakizashi.
This is possible for k8s or other custom monitor to check wakizashi's liveness.
*/
package livenessprobe

import (
	"net/http"
	"strconv"
	"sync"
)

// Probe Interface  probe for k8s on path: /healthz
type Probe interface {
	// Start starts the liveness probe on given prot
	Start(uint16) error
	// Stop will stop the liveness probe
	Stop() error
	// SetLiveness is used to set the liveness
	SetLiveness(bool)
}

// singleton
var probe Probe
var once sync.Once

// GetLivenessProbe get the shipperProbe server
func GetLivenessProbe() Probe {
	once.Do(func() {
		lp := livenessProbe{}
		probe = &lp
	})

	return probe
}

type livenessProbe struct {
	serv            *http.Server
	livenessHandler *livenessHandler
}

// Start will block the process until the inner server is closed, or return err
func (lp *livenessProbe) Start(port uint16) error {
	portStr := strconv.Itoa(int(port))

	lp.livenessHandler = &livenessHandler{
		Liveness: true,
	}

	mux := http.NewServeMux()
	mux.Handle("/healthz", lp.livenessHandler)

	lp.serv = &http.Server{
		Addr:    ":" + portStr,
		Handler: mux,
	}
	err := lp.serv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

// Stop will shutdown the liveness probe
func (lp *livenessProbe) Stop() error {
	err := lp.serv.Shutdown(nil)
	if err != nil {
		return err
	}
	return lp.serv.Close()
}

// SetLiveness can change the liveness of probe, changing its response to 200(true) or 500(false), thread-safe
func (lp *livenessProbe) SetLiveness(liveness bool) {
	lp.livenessHandler.SetLiveness(liveness)
}

type livenessHandler struct {
	Liveness bool
	mtx      sync.Mutex
}

func (lh *livenessHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	lh.mtx.Lock()
	ln := lh.Liveness
	lh.mtx.Unlock()

	if ln {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Liveness probe reporting fine."))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Liveness probe reporting error, wakizashi may be panicing."))
	}
}

func (lh *livenessHandler) SetLiveness(liveness bool) {
	lh.mtx.Lock()
	defer lh.mtx.Unlock()
	lh.Liveness = liveness
}
