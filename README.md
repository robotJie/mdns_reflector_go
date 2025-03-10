# mDNS 反射器

这是一个用 Go 语言编写的 mDNS 反射器，可以在不同网络接口之间转发 mDNS 报文。

## 功能

- 在指定的网络接口之间转发 mDNS 报文
- 支持通过命令行参数或配置文件指定网络接口
- 支持通过命令行命令更新配置文件

## 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/mdns-reflector.git
cd mdns-reflector

# 编译
go build -o mdns-go
```

## 使用方法

### 通过命令行参数指定接口

```bash
# 使用命令行参数指定接口
./mdns-go --ifaces=eth0,en0
```

### 通过配置文件指定接口

首先创建或更新配置文件：

```bash
# 创建/更新配置文件
./mdns-go config ifaces eth0,en0
```

然后直接运行程序，它会自动从配置文件中读取接口列表：

```bash
# 从配置文件读取接口
./mdns-go
```

## 配置文件

配置文件 `config.yml` 的格式如下：

```yaml
ifaces:
  - eth0
  - en0
```

## 注意事项

- 程序需要以 root 权限运行，因为它需要监听网络接口上的多播地址
- 确保指定的网络接口存在且可用

## 如何查找网络接口名称

* macOS:
  - 对于 Wi-Fi 接口，按住 Option 键点击菜单栏上的 Wi-Fi 图标，接口名称会显示在详细信息中
  - 或者使用命令 `ifconfig` 查看所有网络接口

* Linux:
  - 使用命令 `ip link show` 或 `ifconfig` 查看所有网络接口
