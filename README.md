## Usage
    ./mdns_reflector_arm64 -ifaces=en1,bridge100

    -ifaces string
        指定需要反射mDNS报文的网络接口，使用逗号分隔，例如：-ifaces=eth0,en0

## FAQ
* 如何知道需要进行反射的ifaces name？
    - orbstack侧
    启动orbstack后，`ifconfig`观察输出接口哪个网段跟docker内部的网段匹配
    or
    `dns-sd -B _hap._tcp`后观察 if 列数值（代表interface index），再通过`ip link show` 最前面的数值就是if index了（你可能需要 `brew install iproute2mac` 来使用ip command）

    - 本地网络侧
    wifi的话，直接按住option键点击wifi图标，出现一个窗口，接口名称字段就是了
