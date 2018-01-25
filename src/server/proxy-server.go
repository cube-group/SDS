//author linyang created on 2018-01
//SDS service discovery server
//build by golang
package main

import (
    "server/proxy"
    "os"
    "github.com/spf13/viper"
    "fmt"
    "alex/log"
)

func main() {
    //config init
    viper.SetConfigName("init") // name of config file (without extension)
    viper.AddConfigPath("conf") // path to look for the config file in
    viper.SetConfigType("json")
    err2 := viper.ReadInConfig()
    if err2 != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err2))
    }

    //log init
    log.NewLogger("proxy-server", viper.GetString("proxy.logPath"), "Asia/Shanghai", true)

    //代理逻辑
    err := proxy.NewHttpProxyServer()
    if err != nil {
        os.Exit(0)
    }
}
