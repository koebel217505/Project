// Package clientIP 讀准許連接的Client端IP
package clientIP

import (
	"github.com/koebel217505/Project/projCommon/projTcp"
)

var ClientIPs = make([]projTcp.AddrInfo, 0)

func init() {
	loadClientIP()
}

func loadClientIP() {
	//exePath, _ := os.Executable()
	//path := filepath.Dir(exePath)
	//fileName := filepath.Join(path + "/ClientIPs.csv")
	//fs, err := os.Open(fileName)
	//if err != nil {
	//	return
	//}
	//defer fs.Close()
	//
	//r := csv.NewReader(fs)
	//n := 0
	//for {
	//	row, err := r.Read()
	//	if err != nil && err != io.EOF {
	//		log.Fatalf("can not read, err is %+v", err)
	//	}
	//	if err == io.EOF {
	//		break
	//	}
	//	//fmt.Println(row)
	//	n := n + 1
	//	ClientIPs = append(ClientIPs, projTcp.AddrInfo{
	//		Name: "DB" + projConvert.ConvString(n),
	//		IP:   row[0],
	//	})
	//}
}
