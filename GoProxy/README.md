# GoProxy

Go语言写的http代理Demo



## HTTP代理协议流程

1. 客户端给服务器发送个 CONNECT 指令，服务器回复 `HTTP/1.1 200 Connection Established` ，同时建立与远程的连接
2. 服务器无脑转发

> 不知道为什么HTTP协议就直接GET了，也不CONNECT

**练手作品，不喜勿喷**