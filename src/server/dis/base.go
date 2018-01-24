package dis

//micro-service discovery interface
type IDis interface {
    //get the micro-service address by the micro-service-name
    Get(ms string) string
    //set the micro-service
    Set(ms string, version string, address string)
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
    List map[string]DisMsItem
}

//get micro-service address
//load-balance
func (this *Dis)Get(ms string) string {
    //...todo override
    return ""
}

//set micro-service address
//ms contains name
//if ms has the version, make sure that the name contains ":version", such as "ucenter:1.0.0"
func (this *Dis)Set(ms, address string) bool {
    //...todo override
    return false
}

//proxy success
func (this *Dis)Success(ms string) {
    //...todo override
}

//proxy failed
func (this *Dis)Failed(ms string) {
    //...todo override
}