package dis

import (
    "errors"
    "fmt"
    "sort"
    "time"
    "alex/utils"
    "github.com/spf13/viper"
)

//micro-service discovery interface
type IDis interface {
    Init()
    //get the micro-service address by the micro-service-name
    Get(ms string) (*DisMsItem, error)
    //set the micro-service
    Set(ms, address string) error
    //success
    Success(item *DisMsItem)
    //failed
    Failed(item *DisMsItem)
}

//discovery micro-service item
type DisMsItem struct {
    //micro-service name
    Name         string
    //micro-service address(such as 127.0.0.1:80)
    Address      string
    //name-address
    UniqueKey    string
    //request total count
    RequestCount int
    //request success count
    SuccessCount int
    //request failed count
    FailedCount  int
    //proxy average use time
    UseTime      int64
    //latest time
    LatestTime   time.Time
}

//DisMsItem clone
func (this *DisMsItem)Clone() *DisMsItem {
    return &DisMsItem{
        Name:this.Name,
        Address:this.Address,
        UniqueKey:this.UniqueKey,
        RequestCount:this.RequestCount,
        SuccessCount:this.SuccessCount,
        FailedCount:this.FailedCount,
        UseTime:this.UseTime,
        LatestTime:this.LatestTime,
    }
}

//实现排序接口
type DisMsItemList []DisMsItem

func (p DisMsItemList) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

func (p DisMsItemList) Len() int {
    return len(p)
}

func (p DisMsItemList) Less(i, j int) bool {
    if p[i].UseTime < p[j].UseTime {
        return false
    }
    return true
}

//micro-service discovery basic class
//List Structure {"app":{"10.0.0.1:9090":{"RequestCount":...}}}
type Dis struct {
    //key is the ms-address:DisMsItem
    List    map[string]*DisMsItem
    //key is the ms:address
    Current map[string]*DisMsItem
}

//dis initialize
func (this *Dis)Init() {
    this.List = map[string]*DisMsItem{}
    this.Current = map[string]*DisMsItem{}

    timer := new(utils.Timer)
    timer.Handler = this.onRefreshSort
    timer.Start(uint64(viper.GetInt64("proxy.loadBalance.refreshInterval")))
}

//刷新排序
func (this *Dis)onRefreshSort(args ...interface{}) {
    for _, item := range this.Current {
        this.refreshMsSort(item.Name)
    }
}

//具体刷新某一个微服务
func (this *Dis)refreshMsSort(ms string) {
    var msList DisMsItemList = DisMsItemList{}
    for _, item := range this.List {
        if item.Name == ms {
            msList = append(msList, *item)
        }
    }
    if len(msList) > 0 {
        sort.Sort(msList)
        this.Current[msList[0].Name] = msList[0].Clone()
    }
}

//get micro-service address
//load-balance
func (this *Dis)Get(ms string) (*DisMsItem, error) {
    item, ok := this.Current[ms]
    if !ok {
        return nil, errors.New(fmt.Sprintf("no ms [%v]", ms))
    }

    this.List[item.UniqueKey].LatestTime = time.Now()
    return item, nil
}

//set micro-service address
//ms contains name
//if ms has the version, make sure that the name contains ":version", such as "ucenter:1.0.0"
func (this *Dis)Set(ms, address string) error {
    uniqueKey := fmt.Sprintf("%v-%v", ms, address)
    _, ok := this.List[uniqueKey]
    if !ok {
        this.List[uniqueKey] = &DisMsItem{Name:ms, Address:address, UniqueKey:uniqueKey}
    }
    this.refreshMsSort(ms)
    fmt.Println("[Register]", this.List)
    fmt.Println("[Register]", this.Current)
    return nil
}

//proxy success
func (this *Dis)Success(item *DisMsItem) {
    uniqueKey := item.UniqueKey
    item, ok := this.List[uniqueKey]
    if ok {
        if !item.LatestTime.IsZero() {
            useTime := time.Now().UnixNano() - this.List[uniqueKey].LatestTime.UnixNano()
            if item.UseTime > 0 {
                this.List[uniqueKey].UseTime = int64((item.UseTime + useTime) / 2)
            } else {
                this.List[uniqueKey].UseTime = useTime
            }
        }
        this.List[uniqueKey].RequestCount = this.List[uniqueKey].RequestCount + 1
        this.List[uniqueKey].SuccessCount = this.List[uniqueKey].SuccessCount + 1
    }
}

//proxy failed
func (this *Dis)Failed(item *DisMsItem) {
    uniqueKey := item.UniqueKey
    _, ok := this.List[uniqueKey]
    if ok {
        this.List[uniqueKey].LatestTime = time.Now()
        this.List[uniqueKey].RequestCount = this.List[uniqueKey].RequestCount + 1
        this.List[uniqueKey].SuccessCount = this.List[uniqueKey].SuccessCount + 1
    }
}