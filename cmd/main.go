package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bp "github.com/easyselfhost/backpack"
	"github.com/golang/glog"
	_ "github.com/rclone/rclone/backend/all"
	_ "github.com/rclone/rclone/fs/sync"
)

var (
	configPath  string
	tryFirst    bool
	tryOnly     bool
	showVersion bool
	commandLine *flag.FlagSet = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

func init() {
	// log to stderr for glog
	flag.Set("logtostderr", "1")
	commandLine.StringVar(&configPath, "config", "", "config file path (required)")
	commandLine.BoolVar(&tryFirst, "try-first", false, "try backup before running cron")
	commandLine.BoolVar(&tryOnly, "try-only", false, "try backup only without starting cron")
	commandLine.BoolVar(&showVersion, "version", false, "show version")
	commandLine.Parse(os.Args[1:])
}

func main() {
	if showVersion {
		fmt.Println("v" + bp.Version)
		return
	}

	if configPath == "" {
		fmt.Println("config path cannot be empty")
		os.Exit(1)
	}

	config, err := bp.ParseConfigFromFile(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if tryFirst || tryOnly {
		for _, rule := range config.BackupRules {
			err = bp.NewBackpackFlow(rule).Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	if tryOnly {
		return
	}

	cron, err := bp.NewCronFromConfig(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	cron.StartAsync()
	glog.Info("Backup crons started")

	<-c
	glog.Info("Signal recieved, gracefully shutting down backup cron")
	cron.Stop()
	glog.Info("Backup cron stopped")
}
