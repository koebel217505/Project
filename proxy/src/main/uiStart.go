package main

import (
	"context"
	"fmt"
	"github.com/koebel217505/Project/projCommon/projChannel"
	"github.com/koebel217505/Project/projCommon/projVar"
	"github.com/zserge/lorca"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func uiStart(c context.Context, closeCh *projChannel.Channel) {
	var args []string
	var err error
	projVar.UI, err = lorca.New("", "", 750, 750, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer projVar.UI.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	projVar.UI.Bind("start", func() {
		log.Println("UI is ready")

		//v := proVar.Server.GetClientsAddrArray()
		//for key, value := range v {
		//	proVar.UI.Eval(fmt.Sprintf(`app.$data.items.push({ID:%d,IP:"%s",Status:"未連線"})`, key, value.IP))
		//}

		//proVar.UI.Eval(fmt.Sprintf(`app.SQLIP="` + loadConfig.Config.SQLIP + ":" + loadConfig.Config.SQLPort + `"`))
		//proVar.Server.SetUIStatus(projType.UIStatus_Start)
	})

	//Create and bind Go object to the UI
	projVar.UI.Bind("UIMsg", func(message string) string {
		fmt.Println(message)
		if message == "clearPlayer" {
			//DataBase.DrapAccountTable()
			//DataBase.CreateAccountTable()
			//proVar.UI.Eval(`
			//	app.alertfuc("清除完成");
			//`)
		}

		return message
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	exePath, _ := os.Executable()
	path := filepath.Dir(exePath)
	go http.Serve(ln, http.FileServer(http.Dir(path+"/ui/proxy")))
	projVar.UI.Load(fmt.Sprintf("http://%s", ln.Addr()))

	for {
		select {
		case <-c.Done():
		case <-projVar.UI.Done():
			closeCh.Push(nil)
			return
		}
	}
}
