//author linyang created on 2018-01
//SDS service discovery server
//build by golang
package main

import (
    "os"
    "github.com/spf13/viper"
    "fmt"
    "server/core"
    "runtime"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    //config init
    viper.SetConfigName("init") // name of config file (without extension)
    viper.AddConfigPath("conf") // path to look for the config file in
    viper.SetConfigType("json")
    err2 := viper.ReadInConfig()
    if err2 != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err2))
    }

    //开启代理服务
    err := core.NewHttpProxyServer()
    if err != nil {
        fmt.Println("[Proxy]Init Error", err.Error())
        os.Exit(0)
    }

    //开启注册服务
    err = core.NewHttpRegisterServer()
    if err != nil {
        fmt.Println("[Register]Init Error", err.Error())
        os.Exit(0)
    }
}
