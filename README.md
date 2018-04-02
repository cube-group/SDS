# SDS
* Service Discovery Server

![](https://github.com/cube-group/SDS/blob/master/images/icon.png)

* SDS Framework

![](https://github.com/cube-group/SDS/blob/master/images/framework.png)

* SDS Simple Register Data Stat

![](https://github.com/cube-group/SDS/blob/master/images/manager.png)
* SDS config
```
{
  "proxy": {
    "address": ":3333",
    "maxBytes": 10240
  },
  "register": {
    "address": ":12345",
    "secret": "asdfafdaafaf12132",
    "expire": 30,
  },
  "dis": {
    "type": "memory",
    "address": ""
  }
}
```
* proxy.address - 代理接口地址
* proxy.maxBytes - 代理最大支持的字节数
* register.address - 微服务注册接口
* register.secret - 微服务注册秘钥(暂未启用)
* register.expire - 微服务注册过期时间(单位:秒)
* dis.type - memory(etcd和zookeeper暂未实现)目前仅支持单点注册和代理
* dis.address - If dis.type is the address of zookeeper, etcd (temporarily not supported)


### 解决goole类库
```
git clone https://github.com/golang/net.git $GOPATH/src/golang.org/net

git clone https://github.com/golang/sys.git $GOPATH/src/golang.org/sys

git clone https://github.com/golang/tools.git $GOPATH/src/golang.org/tools

git clone https://github.com/golang/text.git $GOPATH/src/golang.org/x
```

### 运行
```
$cd src/server && ./server
```
### 微服务注册
```
$curl http://127.0.0.1:12345/register?name=ms&address=10.10.2.2:80
```
### 观察当前proxy代理情况
请求:http://127.0.0.1:12345/manager
### 微服务请求模式
```
$curl http://127.0.0.1:3333/ms/controller/action?querystring
```