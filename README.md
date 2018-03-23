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
    "loadBalance": {
      "refreshInterval": 30
    }
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
* proxy.address - http proxy ip:address
* proxy.loadBalance.refreshInterval - Proxy service regular priority switching frequency (unit: Second)
* register.address - http Micro service registration ip:address
* register.secret - http Micro service registry key
* register.expire - The length of the registration of micro service (unit: Second)
* dis.type - The default discovery mechanism currently only supports the memory mode
* dis.address - If dis.type is the address of zookeeper, etcd, and redis (temporarily not supported)


### 解决goole类库
```
git clone https://github.com/golang/net.git $GOPATH/src/golang.org/net

git clone https://github.com/golang/sys.git $GOPATH/src/golang.org/sys

git clone https://github.com/golang/tools.git $GOPATH/src/golang.org/tools

git clone https://github.com/golang/text.git $GOPATH/src/golang.org/x
```

### 运行
```
cd src/server && ./server
```