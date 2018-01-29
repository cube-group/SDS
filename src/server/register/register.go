package register

import (
    "net/http"
    "github.com/spf13/viper"
    "server/core"
    "fmt"
)

func NewHttpRegisterServer() error {
    http.HandleFunc("/register", onRegister)
    fmt.Println("[Register]Initizlize", viper.GetString("register.address"))
    return http.ListenAndServe(viper.GetString("register.address"), nil)
}

func onRegister(w http.ResponseWriter, req *http.Request) {
    //...todo
    query := req.URL.Query()
    ms := query.Get("ms")
    address := query.Get("address")

    core.Dis().Set(ms, address)
}


