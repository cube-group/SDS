//标准http协议
//GET /favicon.ico HTTP/1.1
//Host: proxy.my.com:3333
//Connection: keep-alive
//User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36
//Accept: image/webp,image/apng,image/*,*/*;q=0.8
//Referer: http://proxy.my.com:3333/ucenter/index.html?aaa=1231&b=1231
//Accept-Encoding: gzip, deflate
//Accept-Language: zh-CN,zh;q=0.9,en;q=0.8
//Cookie: p=XgAAAIvtagAA; s=dpr=1; __utma=144340137.591520436.1504438227.1504438227.1512552533.2; __utmz=144340137.1504438227.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(non
package proxy

import (
    "github.com/spf13/viper"
    "net"
    "fmt"
    "bytes"
    "strings"
    "io"
    "runtime"
)

//error string
var errResponse string = "HTTP/1.1 500 %v\n" +
    "Content-Type: text/html; charset=utf-8\n" +
    "Server: Golang-SDS\n" +
    "X-Powered-By: %v\n"

//http proxy server
func NewHttpProxyServer() error {
    l, err := net.Listen("tcp", viper.GetString("http.address"))
    if err != nil {
        fmt.Println("[err]", "listener", err.Error())
    }
    fmt.Println("NewHttpProxyServer", l.Addr())

    for {
        client, err := l.Accept()
        if err != nil {
            fmt.Println("[err]", "client-accpet", err.Error())
            continue
        }

        go handleClientRequest(client)
    }
}

//获取标准错误内容
func getErrString(err string) []byte {
    return []byte(fmt.Sprintf(errResponse, err, runtime.Version()))
}

//代理用户请求
func handleClientRequest(client net.Conn) {
    globalErr := ""
    fmt.Println("handleClientRequest", client.RemoteAddr())

    defer func() {
        if globalErr != "" {
            client.Write(getErrString(globalErr))
        }
        client.Close()
    }()

    original := make([]byte, 1024);
    _, err := client.Read(original)

    if err != nil {
        globalErr = err.Error()
        return
    }

    var method, query, version string
    firstNPosition := bytes.IndexByte(original, '\n')
    firstLine := string(original[:firstNPosition])
    fmt.Sscanf(firstLine, "%s %s %s", &method, &query, &version)

    leaveBytes := original[firstNPosition + 1:]
    secondNPosition := bytes.IndexByte(leaveBytes, '\n')
    //secondLine := string(leaveBytes[:secondNPosition])
    leaveBytes = leaveBytes[secondNPosition + 1:]

    //解析微服务名称
    queryArr := strings.Split(query, "/")
    if len(queryArr) < 3 {
        globalErr = "query no micro service name"
        return
    }
    ms := queryArr[1]
    //...todo 通过微服务名称获取其相应的ip:port
    msAddress := "127.0.0.1:80"

    newFirstLine := fmt.Sprintf(
        "%v %v %v\n",
        method,
        fmt.Sprintf("/%v", strings.Join(queryArr[2:], "/")),
        version,
    )
    newSecondLine := fmt.Sprintf(
        "Host: %v\n",
        msAddress,
    )
    newBytesGroup := [][]byte{
        []byte(newFirstLine),
        []byte(newSecondLine),
        []byte(fmt.Sprintf("SDS-MS: %v\n", ms)),
        []byte("Server: Golang-SDS\n"),
        []byte(fmt.Sprintf("X-Powered-By: %v\n", runtime.Version())),
        leaveBytes,
    }
    newBytes := bytes.Join(newBytesGroup, []byte(""))

    //获得了请求的host和port，就开始拨号吧
    realClient, err := net.Dial("tcp", msAddress)
    defer func() {
        if realClient != nil {
            realClient.Close()
        }
    }()
    if err != nil {
        globalErr = err.Error()
        return
    }

    realClient.Write(newBytes)
    io.Copy(client, realClient)
}