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

func sendThread() (sendCh *projChannel.Channel) {
	sendCh = projChannel.New(10000)
	go func() {
		for {
			if v := sendCh.Pop(); v != nil {
				if v, ok := v.(func()); ok {
					v()
				}
			}
		}
	}()
	return
}

// eventThread bla-bla
func eventThread() (eventCh *projChannel.Channel) {
	eventCh = projChannel.New(10000)
	go func() {
		for {
			if v := eventCh.Pop(); v != nil {
				if v, ok := v.(func()); ok {
					v()
				}
			}
		}
	}()
	return
}

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

	sendCh := sendThread()
	eventCh := eventThread()
	closeCh := projChannel.New(1)
	ctx, cancel := context.WithCancel(context.Background())
	go scannerStart(ctx, closeCh, eventCh)
	go uiStart(ctx, closeCh, eventCh)
	go serverStart(ctx, closeCh, sendCh, eventCh)
	go clientStart(ctx, closeCh, sendCh, eventCh)

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
