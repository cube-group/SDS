package register

import (
    "net/http"
    "github.com/spf13/viper"
    "server/core"
)

func NewHttpRegisterServer() {
    http.Handle("/register", onRegister)
    http.ListenAndServe(viper.GetString("register.address"), nil)
}

func onInit() {
    //...todo
}

func onRegister(rw http.ResponseWriter, req *http.Request) {
    //...todo
    core.Dis().Set()
}


