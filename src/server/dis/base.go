package dis

import (
    "errors"
    "fmt"
)

//micro-service discovery interface
type IDis interface {
    Init()
    //get the micro-service address by the micro-service-name
    Get(ms string) (string, error)
    //set the micro-service
    Set(ms string, version string, address string) error
}

//discovery micro-service item
type DisMsItem struct {
    //request total count
    RequestCount int
    //request success count
    SuccessCount int
    //request failed count
    FailedCount  int
    //request average time
    AverageTime  int
}

//micro-service discovery basic class
type Dis struct {
    List map[string][]DisMsItem
}

func (this *Dis)Init() {
    this.List = map[string]DisMsItem{}
}

//get micro-service address
//load-balance
func (this *Dis)Get(ms string) (string, error) {
    //...todo coming soon
    _, ok := this.List[ms]
    if !ok {
        return errors.New(fmt.Sprintf("no ms [%v]", ms))
    }
    return "127.0.0.1:80", nil
}

//set micro-service address
//ms contains name
//if ms has the version, make sure that the name contains ":version", such as "ucenter:1.0.0"
func (this *Dis)Set(ms, address string) error {
    //...todo override
    return nil
}

//proxy success
func (this *Dis)Success(ms string) {
    //...todo override
}

//proxy failed
func (this *Dis)Failed(ms string) {
    //...todo override
}