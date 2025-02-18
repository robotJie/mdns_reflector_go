## 使用方法

```bash
./mdns_reflector_arm64 -ifaces=en1,bridge100
```

### 参数说明

- `-ifaces string` 指定需要反射 mDNS 报文的网络接口，使用逗号分隔，例如：`-ifaces=eth0,en0`

## 开机自启

1. 将 `com.user.homeassistant.mdns.plist` 文件添加到 `/Library/LaunchDaemons` 路径下。
2. 确保文件权限正确，通常需要设置为 `root:wheel` 并赋予 `644` 权限。
3. 使用以下命令加载并启动服务：

```bash
sudo launchctl load /Library/LaunchDaemons/com.user.homeassistant.mdns.plist
sudo chmod 644 /Library/LaunchDaemons/com.user.homeassistant.mdns.plist
sudo chown root:wheel /Library/LaunchDaemons/com.user.homeassistant.mdns.plist
sudo chmod +x /path/to/mdns_reflector_arm64
```

### 注意事项
- 确保 `mdns_reflector_arm64` 文件具有可执行权限。
- 如果修改了 `.plist` 文件，需重新加载服务以应用更改。

## FAQ
* 如何知道需要进行反射的ifaces name？
    - orbstack侧
    启动orbstack后，`ifconfig`观察输出接口哪个网段跟docker内部的网段匹配
    or
    `dns-sd -B _hap._tcp`后观察 if 列数值（代表interface index），再通过`ip link show` 最前面的数值就是if index了（你可能需要 `brew install iproute2mac` 来使用ip command）

    - 本地网络侧
    wifi的话，直接按住option键点击wifi图标，出现一个窗口，接口名称字段就是了
