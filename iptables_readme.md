
### table 规则表：
+ filter表——三个链：INPUT、FORWARD、OUTPUT 
  +  作用：过滤数据包 内核模块：iptables_filter.
+ Nat表——三个链：PREROUTING、POSTROUTING、OUTPUT
    +  作用：用于网络地址转换（IP、端口） 内核模块：iptable_nat
+ Mangle表——五个链：PREROUTING、POSTROUTING、INPUT、OUTPUT、FORWARD
    +  作用：修改数据包的服务类型、TTL、并且可以配置路由实现QOS内核模块：iptable_mangle(别看这个表这么麻烦，咱们设置策略时几乎都不会用到它),
+ Raw表——两个链：OUTPUT、PREROUTING
    +  作用：决定数据包是否被状态跟踪机制处理 内核模块：iptable_raw

### 规则表之间的优先顺序
+ Raw > mangle > nat> > filter

### list
```shell
iptables --list -n

iptables -nL --line-number

```

查看已添加的iptables规则
```shell
iptables -L -n -v
```

删除已添加的iptables规则
#将所有iptables以序号标记显示，执行：
```shell
iptables -L -n --line-numbers
iptables -D INPUT 8
```

### save
```shell
/etc/rc.d/init.d/iptables save

service iptables save
```

iptables常见操作案例
开启Web服务端口
```shell
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A OUTPUT -p tcp --sport 80 -j ACCEPT
```

开启邮件服务的25、110端口
```shell
iptables -A INPUT -p tcp --dport 110 -j ACCEPT
iptables -A INPUT -p tcp --dport 25 -j ACCEPT
```

开启FTP服务的21端口
```shell
iptables -A INPUT -p tcp --dport 21 -j ACCEPT
iptables -A INPUT -p tcp --dport 20 -j ACCEPT
```

开启DNS服务的53端口,假设本机开启了DNS服务
```shell
iptables -A INPUT -p tcp --dport 53 -j ACCEPT
```

允许icmp服务进出
```shell
iptables -A INPUT -p icmp -j ACCEPT
iptables -A OUTPUT -p icmp -j ACCEP
```

假设OUTPUT默认为DROP的话)
允许loopback,不然会导致DNS无法正常关闭等问题，假设默认INPUT DROP，需要开放
```shell
IPTABLES -A INPUT -i lo -p all -j ACCEPT
```

假设默认OUTPUT DROP，需要开放
```shell
IPTABLES -A OUTPUT -o lo -p all -j ACCEPT
```

减少不安全的端口连接
```shell
iptables -A OUTPUT -p tcp --sport 31337 -j DROP
iptables -A OUTPUT -p tcp --dport 31337 -j DROP
```

只允许192.168.0.3的机器进行SSH连接
```shell
iptables -A INPUT -s 192.168.0.3 -p tcp --dport 22 -j ACCEPT
```

如果允许或限制一段IP地址可用192.168.0.0/24表示192.168.0.1-255端的所有IP, 24表示子网掩码数。
```shell
iptables -A INPUT -s 192.168.0.0/24 -p tcp --dport 22 -j ACCEPT

注意：指定某个主机或者某个网段进行SSH连接，需要把iptables配置文件中的-A INPUT -p tcp -m tcp --dport 22 -j ACCEPT删除
因为它表示所有地址都可以登陆.
```

处理IP碎片数量，防止DDOS攻击，允许每秒100个
```shell
iptables -A FORWARD -f -m limit --limit 100/s --limit-burst 100 -j ACCEPT
```

设置ICMP包过滤, 允许每秒1个包, 限制触发条件是10个包
```shell
iptables -A FORWARD -p icmp -m limit --limit 1/s --limit-burst 10 -j ACCEPT
```

DROP非法连接
```shell
iptables -A INPUT   -m state --state INVALID -j DROP
iptables -A OUTPUT  -m state --state INVALID -j DROP
iptables -A FORWARD -m state --state INVALID -j DROP
```

允许所有已经建立的和相关的连接
```shell
iptables-A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables-A OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
```


假设默认为ACCEPT,防止外网用内网IP欺骗
```shell
iptables -t nat -A PREROUTING -i eth0 -s 10.0.0.0/8 -j DROP
iptables -t nat -A PREROUTING -i eth0 -s 172.16.0.0/12 -j DROP
iptables -t nat -A PREROUTING -i eth0 -s 192.168.0.0/16 -j DROP
```

禁用FTP(21)端口
```shell
iptables -t nat -A PREROUTING -p tcp --dport 21 -j DROP
```

配置白名单
```shell
iptables -A INPUT -s 10.0.0.0/8 -j ACCEPT
iptables -A INPUT -s 172.16.0.0/12 -j ACCEPT
iptables -A INPUT -s 192.168.0.0/16 -j ACCEPT
```

屏蔽IP
屏蔽单个IP的命令
```shell
iptables -I INPUT -s 123.45.6.7 -j DROP
```

封整个段即从123.0.0.1到123.255.255.254的命令
```shell
iptables -I INPUT -s 123.0.0.0/8 -j DROP
```

封IP段即从123.45.0.1到123.45.255.254的命令
```shell
iptables -I INPUT -s 124.45.0.0/16 -j DROP
```

封IP段即从123.45.6.1到123.45.6.254的命令是
```shell
iptables -I INPUT -s 123.45.6.0/24 -j DROP
```

屏蔽恶意主机(比如，192.168.0.8)
```shell
iptables -A INPUT -p tcp -m tcp -s 192.168.0.8 -j DROP
```

NAT网络转发规则
内网所有IP:172.16.93.0/24转为10.0.0.1
```shell
iptables -t nat -A POSTROUTING -s 172.16.93.0/24  -j SNAT --to-source 10.0.0.1
```

内网IP192.168.1.100转为10.0.0.1
```shell
iptables -t nat -A POSTROUTING -s 192.168.1.100  -j SNAT --to-source 10.0.0.1
```

假设eth0的IP经常变化,做NAT的方法
```shell
iptables -t nat -A POSTROUTING -s 10.8.0.0/255.255.255.0 -o eth0 -j MASQUERADE
```

ADSL拔号,IP经常变化,做NAT的方法
```shell
iptables -t nat -A POSTROUTING -j MASQUERADE -o pppXXX
```
