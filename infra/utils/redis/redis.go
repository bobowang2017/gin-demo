package redis

import (
	"encoding/json"
	"gin-demo/infra/utils/config"
	"gin-demo/infra/utils/log"
	"github.com/gomodule/redigo/redis"
	"time"
)

var pool *redis.Pool

func SetUp(redisCfg config.Redis) {
	pool = &redis.Pool{
		MaxIdle:     redisCfg.MaxIdle,
		MaxActive:   redisCfg.MaxActive,
		IdleTimeout: redisCfg.IdleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisCfg.Host,
				redis.DialConnectTimeout(redisCfg.ConnectTimeout*time.Millisecond),
				redis.DialReadTimeout(redisCfg.ReadTimeout*time.Millisecond),
				redis.DialWriteTimeout(redisCfg.WriteTimeout*time.Millisecond))
			if err != nil {
				log.Logger.Error(err)
				return nil, err
			}
			// 选择db
			_, _ = c.Do("AUTH", redisCfg.Password)
			return c, nil
		},
	}
}

func closeCon(con redis.Conn) {
	_ = con.Close()
}

func Exists(key string) bool {
	conn := pool.Get()
	defer closeCon(conn)
	exist, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis EXIST Error")
		return false
	}
	return exist
}

func Set(key string, data interface{}, time int) error {
	/**
	Redis String set操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	value, err := json.Marshal(data)
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis Marshal Error")
		return err
	}
	if time == -1 {
		_, err = conn.Do("SET", key, value)
	} else {
		_, err = conn.Do("SET", key, value, "EX", time)
	}
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis Set Error")
		return err
	}
	return nil
}

func SetNx(key string, data interface{}, time int) error {
	/**
	Redis String setNx操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = redis.String(conn.Do("SET", key, value, "EX", time, "NX"))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis SetNx Error")
		return err
	}
	return nil
}

func Expire(key string, time int) error {
	/**
	Redis Expire操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	_, err := conn.Do("EXPIRE", key, time)
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis Expire Error")
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	/**
	Redis String get操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	res, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis Get Error")
	}
	return res, nil
}

func Delete(key string) (bool, error) {
	/**
	Redis String delete操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	res, err := redis.Bool(conn.Do("DEL", key))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis Delete Error")
	}
	return res, err
}

func LPush(key string, data []string) error {
	/**
	Redis String LPush操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	_, err := conn.Do("LPush", redis.Args{}.Add(key).AddFlat(data)...)
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis LPush Error")
		return err
	}
	return nil
}

func RPush(key string, data []string) error {
	/**
	Redis String LPush操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	_, err := conn.Do("RPush", redis.Args{}.Add(key).AddFlat(data)...)
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis RPush Error")
		return err
	}
	return nil
}

func LRange(key string, start, end int) ([]string, error) {
	/**
	Redis List LRange操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	var result []string
	values, err := redis.Values(conn.Do("LRange", key, start, end))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis LRange Error")
		return result, err
	}
	for _, v := range values {
		result = append(result, string(v.([]byte)))
	}
	return result, nil
}

func LTrim(key string, start, end int) error {
	/**
	Redis List LTrim操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	//这里加不加redis.String都可以，如果操作成功返回的都是字符串OK
	if _, err := redis.String(conn.Do("ltrim", key, start, end)); err != nil {
		log.Logger.Error(err.Error(), "| Redis LTrim Error")
		return err
	}
	return nil
}

func LLen(key string) (int, error) {
	/**
	Redis List Len操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	length, err := redis.Int(conn.Do("LLen", key))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis LLen Error")
		return 0, err
	}
	return length, nil
}

func LPop(key string) (string, error) {
	/**
	Redis List LPop操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	data, err := redis.String(conn.Do("LPop", key))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis LPop Error")
		return "", err
	}
	return data, nil
}

func RPop(key string) (string, error) {
	/**
	Redis List RPop操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	data, err := redis.String(conn.Do("RPop", key))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis RPop Error")
		return "", err
	}
	return data, nil
}

func HSet(key, field, value string) error {
	/**
	Redis String HSet操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	_, err := conn.Do("HSet", key, field, value)
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis HSet Error")
		return err
	}
	return nil
}

func HGetAll(key string) (map[string]string, error) {
	conn := pool.Get()
	defer closeCon(conn)
	result, err := redis.Values(conn.Do("HGetAll", key))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis HGetAll Error")
		return nil, err
	}
	data := make(map[string]string)
	for i := 0; i < len(result); i += 2 {
		key := string(result[i].([]byte))
		value := string(result[i+1].([]byte))
		data[key] = value
	}
	return data, nil
}

func HGet(key, field string) (string, error) {
	/**
	Redis String HGet操作
	*/
	conn := pool.Get()
	defer closeCon(conn)
	res, err := redis.String(conn.Do("HGet", key, field))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis HGet Error")
	}
	return res, err
}

func IncrBy(key string, cnt int) (int64, error) {
	/**
	Redis String IncrBy
	*/
	conn := pool.Get()
	defer closeCon(conn)
	res, err := redis.Int64(conn.Do("IncrBy", key, cnt))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis IncrBy Error")
	}
	return res, nil
}

func DecrBy(key string, cnt int) (int64, error) {
	/**
	Redis String IncrBy
	*/
	conn := pool.Get()
	defer closeCon(conn)
	res, err := redis.Int64(conn.Do("DecrBy", key, cnt))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis DecrBy Error")
	}
	return res, nil
}

func MDelete(keys ...string) error {
	/**
	Redis String MDelete
	*/
	conn := pool.Get()
	defer closeCon(conn)
	_ = conn.Send("MULTI")
	for _, v := range keys {
		if len(v) == 0 {
			continue
		}
		_ = conn.Send("DEL", v)
	}
	_, err := conn.Do("EXEC")
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis MDelete Error")
	}
	return err
}

func SAdd(key string, val ...interface{}) error {
	/**
	Redis SAdd
	*/
	conn := pool.Get()
	defer closeCon(conn)
	args := []interface{}{key}
	for _, v := range val {
		args = append(args, v)
	}
	_, err := redis.Int64(conn.Do("SADD", args...))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis SADD Error")
		return err
	}
	return nil
}

func SCard(key string) (int64, error) {
	/**
	Redis SCard
	*/
	conn := pool.Get()
	defer closeCon(conn)
	res, err := redis.Int64(conn.Do("SCARD", key))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis SCARD Error")
		return res, err
	}
	return res, nil
}

func SRem(key string, val ...interface{}) error {
	/**
	Redis SRem
	*/
	conn := pool.Get()
	defer closeCon(conn)
	args := []interface{}{key}
	for _, v := range val {
		args = append(args, v)
	}
	_, err := redis.Int64(conn.Do("SREM", args...))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis SREM Error")
		return err
	}
	return nil
}

func SIsMember(key, val string) (bool, error) {
	/**
	Redis SIsMember
	*/
	conn := pool.Get()
	defer closeCon(conn)
	res, err := redis.Bool(conn.Do("SISMEMBER", key, val))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis SISMEMBER Error")
		return false, err
	}
	return res, nil
}

func SMembers(key string) ([]string, error) {
	/**
	Redis SMembers
	*/
	conn := pool.Get()
	defer closeCon(conn)
	res, err := redis.Values(conn.Do("SMEMBERS", key))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis SMEMBERS Error")
		return nil, err
	}
	rows := make([]string, 0)
	for _, v := range res {
		rows = append(rows, string(v.([]byte)))
	}
	return rows, nil
}

func ZAdd(key string, data map[string]float64) error {
	conn := pool.Get()
	defer closeCon(conn)
	if data == nil {
		return nil
	}
	params := []interface{}{key}
	for k, v := range data {
		params = append(params, v, k)
	}
	_, err := conn.Do("ZADD", params...)
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis ZAdd Error")
		return err
	}
	return nil
}

func ZRem(key string, members ...interface{}) error {
	conn := pool.Get()
	defer closeCon(conn)
	_, err := conn.Do("ZREM", members...)
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis ZREM Error")
		return err
	}
	return nil
}

func ZRevRange(key string, start, end int) ([]string, error) {
	conn := pool.Get()
	defer closeCon(conn)
	rows, err := redis.Strings(conn.Do("ZREVRANGE", key, start, end))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis ZREVRANGE Error")
		return rows, err
	}
	return rows, nil
}

func ZRange(key string, start, end int) ([]string, error) {
	conn := pool.Get()
	defer closeCon(conn)
	rows, err := redis.Strings(conn.Do("ZRange", key, start, end))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis ZRange Error")
		return rows, err
	}
	return rows, nil
}

func ZCard(key string) (int64, error) {
	conn := pool.Get()
	defer closeCon(conn)
	res, err := redis.Int64(conn.Do("ZCard", key))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis ZCard Error")
		return res, err
	}
	return res, nil
}

func ExecLua(paramCnt int, luaScript string, keys, values []string) (int, error) {
	conn := pool.Get()
	defer closeCon(conn)
	lua := redis.NewScript(paramCnt, luaScript)
	args := redis.Args{}.AddFlat(keys).AddFlat(values)
	res, err := redis.Int(lua.Do(conn, args...))
	if err != nil {
		log.Logger.Error(err.Error(), "| Redis ExecLua Error")
		return 0, err
	}
	return res, nil
}
