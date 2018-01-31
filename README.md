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