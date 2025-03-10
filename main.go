package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 表示配置文件结构
type Config struct {
	Ifaces []string `yaml:"ifaces"`
}

func main() {


	// 定义命令行参数
	var ifaceNames string
	var configIfaces string

	flag.StringVar(&ifaceNames, "ifaces", "", "指定需要反射mDNS报文的网络接口，使用逗号分隔，例如：-ifaces=eth0,en0")
	flag.StringVar(&configIfaces, "config-ifaces", "", "持久化需要反射mDNS报文的网络接口，使用逗号分隔，例如：-config-ifaces=eth0,en0")
	flag.Parse()

	if configIfaces != "" {
		saveConfig(configIfaces)
		return
	}

	// 如果没有通过命令行指定接口，则尝试从配置文件读取
	var ifaceNameList []string
	if ifaceNames == "" {
		config, err := loadConfig()
		if err != nil {
			log.Fatalf("无法加载配置文件: %v", err)
		}

		if len(config.Ifaces) == 0 {
			log.Fatal("必须指定至少一个网络接口，可以通过--config-ifaces参数或config.yml配置文件")
		}

		ifaceNameList = config.Ifaces
		log.Printf("从配置文件加载接口: %v", ifaceNameList)
	} else {
		ifaceNameList = strings.Split(ifaceNames, ",")
		log.Println("从命令行参数加载接口")
	}

	ifaces := make([]*net.Interface, 0, len(ifaceNameList))

	for _, name := range ifaceNameList {
		iface, err := net.InterfaceByName(strings.TrimSpace(name))
		if err != nil {
			log.Fatalf("无法获取接口%s：%v", name, err)
		}
		ifaces = append(ifaces, iface)
	}

	log.Println("mDNS reflector started")

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

// 保存配置到文件
func saveConfig(ifacesStr string) {
	ifacesList := strings.Split(ifacesStr, ",")
	config := Config{
		Ifaces: ifacesList,
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("无法序列化配置: %v", err)
	}

	err = os.WriteFile("config.yml", data, 0644)
	if err != nil {
		log.Fatalf("无法写入配置文件: %v", err)
	}

	fmt.Println("配置已保存到 config.yml")
}

// 从文件加载配置
func loadConfig() (Config, error) {
	config := Config{}

	// 检查配置文件是否存在
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		return config, fmt.Errorf("配置文件 config.yml 不存在")
	}

	data, err := os.ReadFile("config.yml")
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
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
