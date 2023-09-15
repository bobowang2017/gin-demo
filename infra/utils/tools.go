package utils

import (
	"encoding/json"
	"fmt"
	m "gin-demo/core/model"
	"gin-demo/infra/common"
	"gin-demo/infra/dao"
	"gin-demo/infra/model"
	"gin-demo/infra/utils/log"
	"gin-demo/infra/utils/redis"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GetSysCfg 读取系统配置文件并解析成SysCfg对象
func GetSysCfg() *model.SysCfg {
	var (
		content  string
		redisKey = common.SysCfgRedisKey
		expired  = 3600
		cfg      = model.SysCfg{}
		err      error
	)
	if content, err = redis.Get(redisKey); err != nil {
		log.Logger.Error(common.RedisServerError)
	}

	if content == "" {
		sysCfg, err := dao.NewSystemConfigDao().GetUsingCfg()
		if err != nil {
			log.Logger.Error(err.Error())
			panic(common.GetSysCfgError)
		}
		content = sysCfg.Content
		if err := json.Unmarshal([]byte(content), &cfg); err != nil {
			log.Logger.Error(err.Error())
			panic(common.ParseSysCfgError)
		}
		_ = redis.Set(redisKey, &cfg, expired)
		return &cfg
	}

	if err = json.Unmarshal([]byte(content), &cfg); err != nil {
		log.Logger.Error(err.Error())
		panic(common.ParseSysCfgError)
	}
	return &cfg
}

func GetClientIp(c *gin.Context) (string, error) {
	ip := c.Request.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	//从请求头部的X-FORWARDED-FOR获取Ip
	ips := c.Request.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}
	//从请求头部的RemoteAddr获取Ip
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("正确ip获取失败")
}

func MustGetUserId(c *gin.Context) string {
	return c.MustGet("userInfo").(*m.User).UserId
}

func MustGetUser(c *gin.Context) *m.User {
	return c.MustGet("userInfo").(*m.User)
}

type RepeatKey interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | string
}

// RemoveDuplicates 删除排序数组中的重复元素
func RemoveDuplicates[T RepeatKey](a []T) (ret []T) {
	length := len(a)
	for i := 0; i < length; i++ {
		if i > 0 && a[i-1] == a[i] {
			continue
		}
		ret = append(ret, a[i])
	}
	return ret
}

// SafeGo 对于协程内部运行的函数，如果发生panic会导致整个程序崩溃，故需要手动recover
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

// RandomInt 生成指定范围内的随机数
func RandomInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return int(rand.Int63n(int64(max-min))) + min
}

// StringToTime 时间格式字符串转换成go时间
func StringToTime(tmStr string) time.Time {
	parseTime, _ := time.Parse(common.TimeLayout, tmStr)
	return parseTime
}

func StringToDate(tmStr string) time.Time {
	parseTime, _ := time.Parse(common.DateLayout, tmStr)
	return parseTime
}

// TimeToString 将时间转成格式化字符串
func TimeToString(tmTime time.Time) string {
	if tmTime.IsZero() {
		return ""
	}
	return tmTime.Format(common.TimeLayout)
}

// TimeToDateString 将时间转成格式化字符串
func TimeToDateString(tmTime time.Time) string {
	if tmTime.IsZero() {
		return ""
	}
	return tmTime.Format(common.DateLayout)
}

// UnixToTime 13位时间戳转时间
func UnixToTime(e string) (d time.Time, err error) {
	data, err := strconv.ParseInt(e, 10, 64)
	d = time.Unix(data/1000, 0)
	return
}

// InArray 判断整形元素是否在指定的数组中
func InArray[T RepeatKey](target T, source []T) bool {
	for _, element := range source {
		if target == element {
			return true
		}
	}
	return false
}

type ChanKey interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | string
}

type Channel[T ChanKey] struct {
	mut    sync.Mutex
	C      chan T
	closed bool
}

func NewChannel[T ChanKey](size int) *Channel[T] {
	return &Channel[T]{
		C:      make(chan T, size),
		closed: false,
		mut:    sync.Mutex{},
	}
}

func (c *Channel[T]) Close() {
	c.mut.Lock()
	defer c.mut.Unlock()
	if !c.closed {
		close(c.C)
		c.closed = true
	}
}

func (c *Channel[T]) IsClosed() bool {
	c.mut.Lock()
	defer c.mut.Unlock()
	return c.closed
}
