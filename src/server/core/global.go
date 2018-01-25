package core

import (
    "server/dis"
    "github.com/spf13/viper"
)

const (
    //memory mode
    DIS_MEMORY = "momery"
    //etcd k-v mode
    DIS_ETCD = "etcd"
    //zookeeper k-v mode
    DIS_ZOOKEEPER = "zookeeper"
    //redis mode
    DIS_REDIS = "redis"
)

//global discovery instance
var d dis.IDis

//discovery instance
func Dis() dis.IDis {
    if d == nil {
        switch viper.GetString("dis.type"){
        case DIS_MEMORY:
            d = new(dis.DisMemory)
        default:
            d = new(dis.DisMemory)
        }
        d.Init()
    }
    return d
}