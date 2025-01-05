package main

import (
	"flag"
	"log"
	"net"
	"strings"
)

func main() {
	log.Println("mDNS reflector started")
	var ifaceNames string
	flag.StringVar(&ifaceNames, "ifaces", "", "指定需要反射mDNS报文的网络接口，使用逗号分隔，例如：-ifaces=eth0,en0")
	flag.Parse()

	if ifaceNames == "" {
		log.Fatal("必须指定至少一个网络接口")
	}
	log.Println("ifaces params parsed")
	ifaceNameList := strings.Split(ifaceNames, ",")
	ifaces := make([]*net.Interface, 0, len(ifaceNameList))

	for _, name := range ifaceNameList {
		iface, err := net.InterfaceByName(strings.TrimSpace(name))
		if err != nil {
			log.Fatalf("无法获取接口%s：%v", name, err)
		}
		ifaces = append(ifaces, iface)
	}

	// 存储每个接口对应的连接
	conns := make(map[string]*net.UDPConn)

	// mDNS多播地址和端口
	mdnsAddr := net.UDPAddr{
		IP:   net.ParseIP("224.0.0.251"),
		Port: 5353,
	}

	for _, iface := range ifaces {
		// 在指定接口上监听UDP多播地址
		conn, err := net.ListenMulticastUDP("udp4", iface, &mdnsAddr)
		if err != nil {
			log.Fatalf("无法在接口%s上监听mDNS报文：%v", iface.Name, err)
		}

		// 设置读取缓冲区大小
		err = conn.SetReadBuffer(4096)
		if err != nil {
			log.Printf("无法设置读取缓冲区：%v", err)
		}

		conns[iface.Name] = conn

		// 启动goroutine读取报文
		go readPackets(iface, conn, conns)
	}

	// 阻塞主线程
	select {}
}

func readPackets(srcIface *net.Interface, srcConn *net.UDPConn, conns map[string]*net.UDPConn) {
	buf := make([]byte, 65535)
	for {
		n, _, err := srcConn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("读取报文错误：%v", err)
			continue
		}

		// 打印收到的报文信息（可选）
		// fmt.Printf("收到来自接口%s的报文，源地址：%v，长度：%d字节\n", srcIface.Name, srcAddr, n)

		// 转发报文到其他接口
		forwardPacket(srcIface, buf[:n], conns)
	}
}

func forwardPacket(srcIface *net.Interface, packet []byte, conns map[string]*net.UDPConn) {
	for name, conn := range conns {
		if name == srcIface.Name {
			// 不将报文发送回其来源接口
			continue
		}

		// 发送报文到多播地址
		_, err := conn.WriteToUDP(packet, &net.UDPAddr{
			IP:   net.ParseIP("224.0.0.251"),
			Port: 5353,
		})
		if err != nil {
			log.Printf("在接口%s上发送报文错误：%v", name, err)
		} else {
			// 打印转发信息（可选）
			// fmt.Printf("将在接口%s上转发报文，长度：%d字节\n", name, len(packet))
		}
	}
}
