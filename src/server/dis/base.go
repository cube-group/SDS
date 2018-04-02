package dis

import (
    "errors"
    "fmt"
    "sort"
    "time"
    "alex/utils"
    "github.com/spf13/viper"
    "runtime"
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
    //get render data
    GetData() map[string]interface{}
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
    //register time
    RegisterTime time.Time
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
        RegisterTime:this.RegisterTime,
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
    timer.Start(uint64(viper.GetInt64("register.expire")))
}

//刷新排序
func (this *Dis)onRefreshSort(args ...interface{}) {
    fmt.Println("onRefreshSort")
    for ms, _ := range this.Current {
        this.refreshMsSort(ms)
    }
}

//具体刷新某一个微服务
func (this *Dis)refreshMsSort(ms string) {
    var msList DisMsItemList = DisMsItemList{}
    for key, item := range this.List {
        if !item.RegisterTime.IsZero() && (time.Now().Unix() - item.RegisterTime.Unix()) > viper.GetInt64("register.expire") {
            delete(this.List, key)
            continue
        }
        if item.Name == ms {
            msList = append(msList, *item)
        }
    }
    if len(msList) > 0 {
        sort.Sort(msList)
        this.Current[ms] = msList[0].Clone()
    } else {
        delete(this.Current, ms)
    }

    fmt.Println("[Register]", this.List)
    fmt.Println("[Register]", this.Current)
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
    this.List[uniqueKey].RegisterTime = time.Now()
    this.refreshMsSort(ms)
    return nil
}

//proxy success
func (this *Dis)Success(item *DisMsItem) {
    uniqueKey := item.UniqueKey
    item, ok := this.List[uniqueKey]
    if ok {
        useTime := time.Now().UnixNano() - this.List[uniqueKey].LatestTime.UnixNano()
        if !item.LatestTime.IsZero() {
            if item.UseTime > 0 {
                this.List[uniqueKey].UseTime = int64((item.UseTime + useTime) / 2)
            } else {
                this.List[uniqueKey].UseTime = useTime
            }
        } else {
            this.List[uniqueKey].UseTime = useTime
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
        this.List[uniqueKey].FailedCount = this.List[uniqueKey].FailedCount + 1
    }
}

//get the more data
func (this *Dis)GetData() map[string]interface{} {
    //统计信息
    newStat := map[string]interface{}{}
    totalRequest := 0
    totalSuccess := 0
    totalFailed := 0
    var maxRequestItem *DisMsItem
    var minRequestItem *DisMsItem
    for _, i := range this.List {
        totalRequest += i.RequestCount
        totalSuccess += i.SuccessCount
        totalFailed += i.FailedCount
        if maxRequestItem == nil {
            maxRequestItem = i
        } else if maxRequestItem.RequestCount < i.RequestCount {
            maxRequestItem = i
        }
        if minRequestItem == nil {
            minRequestItem = i
        } else if minRequestItem.RequestCount > i.RequestCount {
            minRequestItem = i
        }
    }
    if totalRequest > 0 {
        newStat["successPercent"] = int(totalSuccess * 100 / totalRequest)
        if maxRequestItem != nil {
            newStat["maxPercent"] = int(maxRequestItem.RequestCount * 100 / totalRequest)
        }
        if minRequestItem != nil {
            newStat["minPercent"] = int(minRequestItem.RequestCount * 100) / totalRequest
        }
    } else {
        newStat["successPercent"] = 0
        newStat["maxPercent"] = 0
        newStat["minPercent"] = 0
    }
    if maxRequestItem != nil {
        newStat["maxName"] = maxRequestItem.UniqueKey
    } else {
        newStat["maxName"] = "maxRequestMicroServiceName"
    }
    if minRequestItem != nil {
        newStat["minName"] = minRequestItem.UniqueKey
    } else {
        newStat["minName"] = "minRequestMicroServiceName"
    }
    newStat["gonum"] = runtime.NumGoroutine()

    newList := map[string]DisMsItemList{}
    for ms, _ := range this.Current {
        var msList DisMsItemList = DisMsItemList{}
        for _, i := range this.List {
            if i.Name == ms {
                msList = append(msList, *i)
            }
        }
        newList[ms] = msList
    }

    return map[string]interface{}{"list":newList, "stat":newStat}
}