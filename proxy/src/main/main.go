package main

import (
	"context"
	"fmt"
	"github.com/koebel217505/Project/projCommon/projChannel"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"os/signal"
	"path/filepath"
)

func main() {
	n := 0
	exePath, _ := os.Executable()
	var pids, _ = process.Pids()
	for _, pid := range pids {
		pn, _ := process.NewProcess(pid)
		pName, _ := pn.Name()
		if pName == filepath.Base(exePath) {
			fmt.Println(pName)
			n++
			if n > 1 {
				return
			}
		}
	}

	closeCh := projChannel.NewChannel(1)
	ctx, cancel := context.WithCancel(context.Background())
	go scannerStart(ctx, closeCh)
	go uiStart(ctx, closeCh)
	go serverStart(ctx, closeCh)
	go clientStart(ctx, closeCh)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-closeCh.Self():
		cancel()
	}

	log.Println("exiting...")
}
