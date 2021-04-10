package main

import (
	"BlankZhu/wakizashi/pkg/config"
	"BlankZhu/wakizashi/pkg/constant"
	"BlankZhu/wakizashi/pkg/device"
	"BlankZhu/wakizashi/pkg/dump"
	"BlankZhu/wakizashi/pkg/report"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const title = `
__          __   _    _              _     _ 
\ \        / /  | |  (_)            | |   (_)
 \ \  /\  / /_ _| | ___ ______ _ ___| |__  _ 
  \ \/  \/ / _' | |/ / |_  / _' / __| '_ \| |
   \  /\  / (_| |   <| |/ / (_| \__ \ | | | |
	\/  \/ \__,_|_|\_\_/___\__,_|___/_| |_|_|
                        ======= Probe =======
`

var (
	buildTime    string
	buildVersion string
	gitCommitID  string
)

func launchDumping(devs []net.Interface, conf *config.ProbeConfig, fileCh chan<- string) {
	var wg sync.WaitGroup
	for _, dev := range devs {
		wg.Add(1)
		dumper := dump.Dumper{
			DumpDir:        conf.DumpDir,
			FileCh:         fileCh,
			Iface:          &dev,
			RepAddr:        conf.CenterAddr,
			RotateInterval: time.Duration(conf.CapInterval) * time.Second,
			SnapLen:        uint32(256),
		}
		dumper.Init()

		go func() {
			defer wg.Done()
			dumper.Start()
		}()
	}
	wg.Wait()
}

func launchReporting(devs []net.Interface, conf *config.ProbeConfig, fileCh <-chan string) {
	var wg sync.WaitGroup
	reporter := report.Reporter{
		AutoClear:   conf.AutoClear,
		DumpDir:     conf.DumpDir,
		FileCh:      fileCh,
		Ifaces:      devs,
		RepAddr:     conf.CenterAddr,
		RepInterval: time.Duration(conf.CapInterval/2) * time.Second,
		RepRetry:    conf.UploadRetry,
	}
	reporter.Init()
	wg.Add(1)
	go func() {
		defer wg.Done()
		reporter.Start()
	}()
	wg.Wait()
}

func main() {
	cfgPathPtr := flag.String("c", constant.ProbeDefaultConfigPath, "path to probe's config yaml file")
	verPtr := flag.Bool("v", false, "print version info")
	flag.Parse()

	fmt.Println(title)
	fmt.Printf("Build time: %s\nBuild version: %s\nGit commit ID: %s\n", buildTime, buildVersion, gitCommitID)
	if *verPtr {
		return
	}

	// load config
	conf := config.ProbeConfig{}
	if err := conf.LoadConfigFromYAML(*cfgPathPtr); err != nil {
		logrus.Fatalf("failed to load config from %s, detail: %s", *cfgPathPtr, err)
	}
	logrus.SetLevel(logrus.Level(conf.LogLev))
	logrus.Infof("config loaded: %s", conf.ToString())
	if err := conf.CreateDumpDir(); err != nil {
		logrus.Fatalf("failed to create dump directory on path %s, detail: %s", conf.DumpDir, err)
	}

	// get network devices
	devs := make([]net.Interface, 0)
	for _, regex := range conf.NetworkDevs {
		tmp, err := device.GetNetworkDevices(regex)
		if err != nil {
			logrus.Fatalf("failed to get device on regex %s, detail: %s", regex, err)
		}
		devs = append(devs, tmp...)
	}
	if len(devs) == 0 {
		logrus.Fatalf("no network device dectected, check the network environment")
	}

	fileCh := make(chan string, constant.DefaultChanCap)
	var wg sync.WaitGroup
	go func() {
		launchDumping(devs, &conf, fileCh)
		defer wg.Done()
	}()
	go func() {
		launchReporting(devs, &conf, fileCh)
		defer wg.Done()
	}()
	wg.Wait()
	logrus.Warn("wakizashi probe exit after dumper & repoter returned")
}
