package wasm_ipmac

import (
	"fmt"
	"net"
	"strings"
	"syscall/js"
)


const Pair = `{"ip":%s,"mac":%s}`

type IPMac struct {
	IP string
	Mac string
}

func main()  {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("test")
		cb.Release()
		return getIPMac()
	})
	js.Global().Set("getIPMac", cb)
}

func getIPMacs() string {
	return GetIPMac()
}

//GetIPMac get ip:mac pairs
func GetIPMac() string {
	ipMacs := getIPMac()
	data := make([]string, 0)
	for _, im := range ipMacs {
		pair := fmt.Sprintf(Pair, im.IP, im.Mac)
		data = append(data, pair)
	}
	result := strings.Join(data, ",")
	return "[" +result +"]"
}

//getIPMac get mac
func getIPMac() []IPMac {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	ipMacs := make([]IPMac, 0)
	for _, inter := range interfaces {
		addrs, err := inter.Addrs()
		if err != nil {
			return ipMacs
		}

		if inter.HardwareAddr.String() != "" {
			var ipMac IPMac
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
						ipMac.IP = ipnet.IP.String()
					}
				}
			}
			ipMac.Mac = inter.HardwareAddr.String()
			ipMacs = append(ipMacs, ipMac)
		}
	}
	return ipMacs
}

