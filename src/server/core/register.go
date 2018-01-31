package core

import (
    "net/http"
    "github.com/spf13/viper"
    "fmt"
    "html/template"
    "alex/auth"
)

var secret string
func NewHttpRegisterServer() error {
    //register secret check
    secret = viper.GetString("register.secret")
    //staticDir := viper.GetString("register.static")
    http.HandleFunc("/register", onRegister)
    http.HandleFunc("/manager", onManager)
    http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
    fmt.Println("[Register]Initizlize", viper.GetString("register.address"))
    return http.ListenAndServe(viper.GetString("register.address"), nil)
}

//ms register
func onRegister(w http.ResponseWriter, req *http.Request) {
    //...todo
    query := req.URL.Query()

    ms := query.Get("ms")
    address := query.Get("address")
    if secret != "" {
        sign := query.Get("sign")
        time := query.Get("time")
        signMap := map[string]interface{}{
            "ms":ms,
            "address":address,
            "time":time,
        }
        if (!auth.BasicAuthCheckSign(sign, signMap)) {
            w.Write([]byte(`{"code":100,"msg":"basic auth error"}`))
            return
        }
    }

    err := Dis().Set(ms, address)
    if err != nil {
        w.Write([]byte(fmt.Sprintf(`{"code":100,"msg":"%v"}`, err.Error())))
    } else {
        w.Write([]byte(`{"code":0,"msg":"register success"}`))
    }
}

//display manager page
func onManager(w http.ResponseWriter, req *http.Request) {
    t, err := template.ParseFiles("view/layout.html", "view/manager.html", "view/footer.html")
    if err != nil {
        fmt.Println(err.Error())
        w.Write([]byte("无法渲染页面"))
    } else {
        t.Execute(w, Dis().GetData())
    }
}


