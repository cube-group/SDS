package dis

import (
    "errors"
    "fmt"
    "sort"
    "math"
    "time"
    "alex/utils"
    "github.com/spf13/viper"
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
    //micro-service name
    Name         string
    //micro-service address(such as 127.0.0.1:80)
    Address      string
    //request total count
    RequestCount int
    //request success count
    SuccessCount int
    //request failed count
    FailedCount  int
    //request average time
    AverageTime  int
    //latest time
    LatestTime   time.Time
}

//DisMsItem clone
func (this *DisMsItem)Clone() DisMsItem {
    return DisMsItem{
        Name:this.Name,
        Address:this.Address,
        RequestCount:this.RequestCount,
        SuccessCount:this.SuccessCount,
        FailedCount:this.FailedCount,
        AverageTime:this.AverageTime,
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
    if math.Abs(p[i].RequestCount - p[j].RequestCount) > 100000 {
        return p[i].RequestCount > p[j].RequestCount
    } else {
        return p[i].SuccessCount / p[i].RequestCount > p[j].SuccessCount / p[j].RequestCount
    }
}

//micro-service discovery basic class
//List Structure {"app":{"10.0.0.1:9090":{"RequestCount":...}}}
type Dis struct {
    //key is the ms-address:DisMsItem
    List    map[string]*DisMsItem
    //key is the ms:address
    Current map[string]string
}

//dis initialize
func (this *Dis)Init() {
    this.List = map[string]*DisMsItem{}
    this.Current = map[string]string{}

    timer := new(utils.Timer)
    timer.Handler = this.onRefreshSort
    timer.Start(viper.GetInt64("proxy.loadBalance.refreshInterval"))
}

func (this *Dis)onRefreshSort() {
    for _, list := range this.List {
        sort.Sort(list)
    }
}

//get micro-service address
//load-balance
func (this *Dis)Get(ms string) (string, error) {
    //...todo coming soon
    address, ok := this.Current[ms]
    if !ok {
        return errors.New(fmt.Sprintf("no ms [%v]", ms))
    }
    return address, nil
}

//set micro-service address
//ms contains name
//if ms has the version, make sure that the name contains ":version", such as "ucenter:1.0.0"
func (this *Dis)Set(ms, address string) error {
    key := fmt.Sprintf("%v-%v", ms, address)
    _, ok := this.List[key]
    if !ok {
        this.List[key] = &DisMsItem{Name:ms, Address:address}
    }

    var msList DisMsItemList = DisMsItemList{}
    for _, item := range this.List {
        if item.Name == ms {
            msList = append(msList, *item)
        }
    }
    if len(msList) == 0 {
        return errors.New("exception for msList length")
    }
    sort.Sort(msList)
    this.Current[ms] = &msList[0].Clone()
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