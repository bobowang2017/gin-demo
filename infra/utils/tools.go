package utils

import (
	"encoding/json"
	"gin-demo/infra/common"
	"gin-demo/infra/dao"
	"gin-demo/infra/model"
	"gin-demo/infra/utils/log"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"time"
)

//读取系统配置文件并解析成SysCfg对象
func GetSysCfg() (*model.SysCfg, error) {
	var (
		content string
		cfg     = model.SysCfg{}
		err     error
	)

	sysCfg, err := dao.NewSystemConfigDao().GetUsingCfg()
	if err != nil {
		log.Logger.Error(err.Error())
		return nil, errors.New(common.GetSysCfgError)
	}
	content = sysCfg.Content
	if err := json.Unmarshal([]byte(content), &cfg); err != nil {
		log.Logger.Error(err.Error())
		return nil, errors.New(common.ParseSysCfgError)
	}
	return &cfg, nil
}

//删除数组中的重复元素
func RemoveDuplicatesInt(a []int) (ret []int) {
	length := len(a)
	for i := 0; i < length; i++ {
		if i > 0 && a[i-1] == a[i] {
			continue
		}
		ret = append(ret, a[i])
	}
	return ret
}

//对于协程内部运行的函数，如果发生panic会导致整个程序崩溃，故需要手动recover
func SafeGo(do func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Logger.Error(err)
			}
		}()
		do()
	}()
}

//生成指定范围内的随机数
func RandomInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return int(rand.Int63n(int64(max-min))) + min
}

//时间格式字符串转换成go时间
func StringToTime(tmStr string) time.Time {
	parseTime, _ := time.Parse(common.TimeLayout, tmStr)
	return parseTime
}

func StringToDate(tmStr string) time.Time {
	parseTime, _ := time.Parse(common.DateLayout, tmStr)
	return parseTime
}

//将时间转成格式化字符串
func TimeToString(tmTime time.Time) string {
	if tmTime.IsZero() {
		return ""
	}
	return tmTime.Format(common.TimeLayout)
}

//将时间转成格式化字符串
func TimeToDateString(tmTime time.Time) string {
	if tmTime.IsZero() {
		return ""
	}
	return tmTime.Format(common.DateLayout)
}

//13位时间戳转时间
func UnixToTime(e string) (d time.Time, err error) {
	data, err := strconv.ParseInt(e, 10, 64)
	d = time.Unix(data/1000, 0)
	return
}

//判断整形元素是否在指定的数组中
func InArrayInt(target int, source []int) bool {
	for _, element := range source {
		if target == element {
			return true
		}
	}
	return false
}
