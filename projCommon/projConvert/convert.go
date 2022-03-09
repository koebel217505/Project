package projConvert

import (
	"fmt"
	"github.com/cstockton/go-conv"
	"net"
)

func ConvString(form any) string {
	v, _ := conv.String(form)
	return v
}

func ConvUint8(form any) uint8 {
	s, _ := conv.Uint8(form)
	return s
}

func ConvUint16(form any) uint16 {
	v, _ := conv.Uint16(form)
	return v
}

func ConvUint32(form any) uint32 {
	v, _ := conv.Uint32(form)
	return v
}

func ConvUint64(form any) uint64 {
	s, _ := conv.Uint64(form)
	return s
}

func ConvInt8(form any) int8 {
	v, _ := conv.Int8(form)
	return v
}

func ConvInt16(form any) int16 {
	v, _ := conv.Int16(form)
	return v
}

func ConvInt(form any) int {
	v, _ := conv.Int(form)
	return v
}

func ConvInt32(form any) int32 {
	v, _ := conv.Int32(form)
	return v
}

func ConvInt64(form any) int64 {
	v, _ := conv.Int64(form)
	return v
}

func ConvFloat32(form any) float32 {
	v, _ := conv.Float32(form)
	return v
}

func ConvFloat64(form any) float64 {
	v, _ := conv.Float64(form)
	return v
}

func ConvBool(form any) bool {
	v, _ := conv.Bool(form)
	return v
}

func GetIPs() (ips []string) {

	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}
