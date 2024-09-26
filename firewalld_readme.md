
#### reload
```shell
firewall-cmd --reload
```

#### list
```shell
firewall-cmd --zone=public --list-port
firewall-cmd --zone=public --list-services
firewall-cmd --zone=public --list-rich-rules 

firewall-cmd --list-all 只显示public域下的策略
firewall-cmd --list-all-zones 显示所有的策略
```

#### ip
```shell
firewall-cmd --permanent --add-rich-rule='rule family=ipv4 source address="111.225.149.121" drop' 
firewall-cmd --permanent --remove-rich-rule='rule family=ipv4 source address="111.225.149.121" drop'
firewall-cmd --permanent --add-rich-rule='rule family=ipv4 source address="111.225.0.0/16" drop'
firewall-cmd --permanent --remove-rich-rule='rule family=ipv4 source address="111.225.0.0/16" drop'
```

#### port
```shell
firewall-cmd --zone=public --add-port=80/tcp --permanent
firewall-cmd --zone=public --remove-port=80/tcp --permanent

firewall-cmd --permanent --add-rich-rule="rule family="ipv4" source address="192.168.3.12" port protocol="tcp" port="9191" reject"
firewall-cmd --permanent --remove-rich-rule="rule family="ipv4" source address="192.168.3.12" port protocol="tcp" port="9191" reject"

firewall-cmd --permanent --add-rich-rule="rule family="ipv4" source address="192.168.1.0/16" port protocol="tcp" port="9191" reject"
firewall-cmd --permanent --remove-rich-rule="rule family="ipv4" source address="192.168.1.0/16" port protocol="tcp" port="9191" reject"

firewall-cmd --permanent --add-rich-rule="rule family="ipv4" source address="172.16.1.0/24" port protocol="tcp" port="30001-30030" reject"
firewall-cmd --permanent --remove-rich-rule="rule family="ipv4" source address="172.16.1.0/24" port protocol="tcp" port="30001-30030" reject"
```

#### ip + port
```shell
firewall-cmd --permanent --add-rich-rule="rule family="ipv4" source address="192.168.1.48" port protocol="tcp" port="8600" accept"
firewall-cmd --permanent --remove-rich-rule="rule family="ipv4" source address="192.168.1.48" port protocol="tcp" port="8600" accept"
```

#### services
```shell
firewall-cmd --permanent --zone=public --add-service=https
firewall-cmd --permanent --zone=public --remove-service=https

firewall-cmd --zone=public --add-service=ssh/tcp --permanent
firewall-cmd --zone=public --remove-service=ssh --permanent
```

#### rich rule 优先级最高，可支持系统服务、端口号、源地址和目标地址等诸多信息进行更有针对性的策略配置，推荐使用
```shell
firewall-cmd --permanent --add-rich-rule="rule family="ipv4" source address="192.168.1.48" port protocol="tcp" port="8600" accept"
firewall-cmd --permanent --remove-rich-rule="rule family="ipv4" source address="192.168.1.48" port protocol="tcp" port="8600" accept"

firewall-cmd --permanent --zone=public --add-rich-rule="rule family="ipv4" source address="10.52.2.0.0/24" service name="ssh" reject"
firewall-cmd --permanent --zone=public --remove-rich-rule="rule family="ipv4" source address="192.168.10.0/24" service name="ssh" reject"
```

#### reject 和 drop 区别和联系
+ drop 动作只是简单的直接丢弃数据，并不反馈任何回应。
+ reject 动作则会更为礼貌的返回一个拒绝(终止)数据包(TCP FIN或UDP-ICMP-PORT-UNREACHABLE)，明确的拒绝对方的连接动作。
+ 通常，你的局域网内的所有连接规则都应该使用 reject。对于互联网，除了某些可信任的服务器外，来自互联网的连接通常是 drop。 
+ drop 不符合 TCP 连接规范，可能对你的网络造成不可预期或难以诊断的问题，所以在可信任的局域网内使用 reject 无疑更好！
