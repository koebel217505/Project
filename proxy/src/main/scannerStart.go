package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/koebel217505/Project/projCommon/projChannel"
	"os"
	"strings"
)

func scannerStart(c context.Context, closeCh *projChannel.Channel) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		switch {
		case strings.ToLower(scanner.Text()) == strings.ToLower("WindowHide"):
			fmt.Println("WindowHide")
		case strings.ToLower(scanner.Text()) == strings.ToLower("WindowShow"):
			fmt.Println("WindowShow")
		case strings.ToLower(scanner.Text()) == strings.ToLower("send"):
			//data := projTcp.Test{1, 2, 3}
			//dataBytes, _ := restruct.Pack(binary.LittleEndian, &data)
			//projVar.Client.Send(22222, dataBytes)
		}
	}
}
