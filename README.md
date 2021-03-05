# golang-websocket
websocket lib with golang

由于项目需要使用websocket，在网上也没有找到满足需求的开源项目。
索性就自己写了一个，主要需求点就是服务端可以管理多个客户端链接，
同时客户端在断开之后需要能够主动重新链接。

## ***websocket服务端功能***
1、根据client的ip来管理每一个链接；

## ***websocket客户端功能***
2、链接断开后，会周期性发起重连接；
