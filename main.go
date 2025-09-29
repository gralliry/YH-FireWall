package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
)

func main() {
	// 选择网卡，这里用第一个网卡
	ifs, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	if len(ifs) == 0 {
		log.Fatal("找不到网卡")
	}
	for i, dev := range ifs {
		fmt.Printf("[%d] 名字: %s\n", i, dev.Description)
		fmt.Printf("    系统标识: %s\n", dev.Name)
		for _, addr := range dev.Addresses {
			fmt.Printf("    IP: %s\n", addr.IP)
		}
		fmt.Println()
	}
	device := ifs[0].Name
	fmt.Println("使用网卡:", device)

	// 打开网卡，抓所有包
	handle, err := pcap.OpenLive(device, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
}
