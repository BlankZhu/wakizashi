package main

import (
	"BlankZhu/wakizashi/pkg/backend"
	"BlankZhu/wakizashi/pkg/config"
	"BlankZhu/wakizashi/pkg/constant"
	"BlankZhu/wakizashi/pkg/device"
	liveprobe "BlankZhu/wakizashi/pkg/probe"
	"BlankZhu/wakizashi/pkg/recovery"
	"BlankZhu/wakizashi/pkg/transmit"
	"BlankZhu/wakizashi/pkg/util"
	"flag"
	"fmt"
	"net"
	"path"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const title = `
__          __   _    _              _     _ 
\ \        / /  | |  (_)            | |   (_)
 \ \  /\  / /_ _| | ___ ______ _ ___| |__  _ 
  \ \/  \/ / _' | |/ / |_  / _' / __| '_ \| |
   \  /\  / (_| |   <| |/ / (_| \__ \ | | | |
	\/  \/ \__,_|_|\_\_/___\__,_|___/_| |_|_|
                       ======= Center =======
`

var (
	buildTime    string
	buildVersion string
	gitCommitID  string
)

func launchHealthProbe(port uint16, fin chan<- struct{}) {
	p := liveprobe.GetLivenessProbe()
	err := p.Start(port)
	if err != nil {
		logrus.Fatalf("health probe launch error, detail: %s", err)
	}
	p.Stop()
	close(fin)
}

func main() {
	cfgPathPtr := flag.String("c", constant.CenterDefaultConfigPath, "path to center's config yaml file")
	verPtr := flag.Bool("v", false, "print version info")
	flag.Parse()

	fmt.Println(title)
	fmt.Printf("Build time: %s\nBuild version: %s\nGit commit ID: %s\n", buildTime, buildVersion, gitCommitID)
	if *verPtr {
		return
	}

	// load config
	conf := config.CenterConfig{}
	if err := conf.LoadConfigFromYAML(*cfgPathPtr); err != nil {
		logrus.Fatalf("failed to load config from %s, detail: %s", *cfgPathPtr, err)
	}
	logrus.SetLevel(logrus.Level(conf.LogLev))
	logrus.Infof("config loaded: %s", conf.ToString())
	if err := conf.CreateRecoveryDir(); err != nil {
		logrus.Fatalf("failed to create directory for recovery on path %s, detail: %s", conf.RecovDir, err)
	}

	// get center's possible IP address
	devs, err := device.GetAllNetworkDevices()
	if err != nil {
		logrus.Fatalf("failed to get network device, detail: %s", err)
	}
	if len(devs) == 0 {
		logrus.Fatalf("no network device dectected, check the network environment")
	}
	ips := util.GetIPSetFromNetworkInterfaces(devs)

	// initialized data backend client
	backend.Init(conf.BackendConfig)
	cli := backend.Get()
	if cli == nil {
		logrus.Fatalf("failed to get data backend client (nil)")
	}
	err = (*cli).Connect()
	if err != nil {
		logrus.Fatalf("failed to connect to database, detail: %s", err)
	}
	defer (*cli).Close()

	// setup recovery
	r := recovery.Get()
	recoveryPath := path.Join(conf.RecovDir, constant.RecoveryDefaultFileName)
	positionPath := path.Join(conf.RecovDir, constant.RecoveryDefaultPosName)
	r.Init(recoveryPath,
		positionPath,
		constant.RecoveryDefaultPosLimit,
		constant.RecoveryDefaultCacheSize,
		(*cli).Write,
	)
	go func() {
		for {
			<-time.After(time.Duration(conf.RecovInterval) * time.Second)
			r.RepostRecord()
		}
	}()

	// setup health probe
	fin := make(chan struct{}, 1)
	go launchHealthProbe(conf.HealthPort, fin)

	// start listening requests from probe side
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(int(conf.Port)))
	if err != nil {
		logrus.Fatalf("failed to listen on port %d, detail: %s", conf.Port, err)
	}
	logrus.Infof("wakizashi center listening on port %d", conf.Port)
	serv := grpc.NewServer()
	transmit.RegisterTransmitServer(serv, &transmit.CenterServer{IPSet: ips})
	if err := serv.Serve(lis); err != nil {
		logrus.Fatalf("failed to start grpc transmit server, detail: %s", err)
	}

	p := liveprobe.GetLivenessProbe()
	p.SetLiveness(false)
	// should never reaches here if probe meets no error
	<-fin
	logrus.Warn("wakizashi center exit after health probe returned")
}
