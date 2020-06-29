package redis

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

func Init(addr, db, password string, timeout int) {
	if pool == nil {
		pool = newRedisPool(addr, db, password, timeout)
	}
}

// STRING
// --------------------------------------------------------------------------------

func Del(key string) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("DEL", key)
	return err
}

func Expire(key string, expire int) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("EXPIRE", key, expire)
	return err
}

func Persist(key string) error {
	c := pool.Get()
	defer c.Close()
	_, err := c.Do("PERSIST", key)
	return err
}

func Get(key string) (string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.String(c.Do("GET", key))
}
func Incr(key string) (int64, error) {
	c := pool.Get()
	defer c.Close()

	return redis.Int64(c.Do("INCRBY", key, 1))
}

func IncrBy(key string, val int64) (int64, error) {
	c := pool.Get()
	defer c.Close()
	return redis.Int64(c.Do("INCRBY", key, val))
}

func Keys(pattern string) ([]string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.Strings(c.Do("KEYS", pattern))
}

func Set(key string, val []byte) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("SET", key, string(val))
	return err
}

func Setex(key string, expire int, val interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("SETEX", key, expire, val)
	return err
}

func Exist(key string) (bool, error) {
	c := pool.Get()
	if i, err := redis.Int(c.Do("EXISTS", key)); err != nil {
		return false, err
	} else {
		if i > 0 {
			return true, nil
		}
		return false, nil
	}

}

// HASH
// --------------------------------------------------------------------------------

func Hgetall(key string) (map[string]string, error) {
	c := pool.Get()
	defer c.Close()

	res, err := bytesSlice(c.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	writeToContainer(res, reflect.ValueOf(result))

	return result, err
}

func Hmset(key string, mapping interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("HMSET", redis.Args{}.Add(key).AddFlat(mapping)...)
	return err
}

func Hget(key string, filed interface{}) (string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.String(c.Do("HGET", key, filed))
}

func Hset(key string, filed interface{}, val interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("HSET", key, filed, val)
	return err
}

func Hdel(key string, filed interface{}) (interface{}, error) {
	c := pool.Get()
	defer c.Close()

	return c.Do("HDEL", key, filed)
}

// SETS
// --------------------------------------------------------------------------------

func Smembers(key string) ([]string, error) {
	c := pool.Get()
	defer c.Close()
	return redis.Strings(c.Do("SMEMBERS", key))
}

func Sadd(key string, val interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("SADD", key, val)
	return err
}

func Srem(key string, val interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("SREM", key, val)
	return err
}

// ZSETS
// --------------------------------------------------------------------------------

// 有序集合成员设置
func Zadd(key string, score interface{}, member interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("ZADD", key, score, member)
	return err
}

// 有序集合增量修改
func Zincrby(key string, increment int, member interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("ZINCRBY", key, increment, member)
	return err
}

// 删除有序集合成员
func Zrem(key string, member interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("ZREM", key, member)
	return err
}

// 获取集合成员数
func Zcard(key string) (int, error) {
	c := pool.Get()
	defer c.Close()

	return redis.Int(c.Do("ZCARD", key))
}

func Zadds(key string, data ...interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("ZADD", key, data)
	return err
}

func Zrange(key string, withScore bool) ([]string, error) {
	c := pool.Get()
	defer c.Close()

	if withScore {
		return redis.Strings(c.Do("ZRANGEBYSCORE", key, "-INF", "+INF", "WITHSCORES"))
	}
	return redis.Strings(c.Do("ZRANGEBYSCORE", key, "-INF", "+INF"))
}

// 上面的引用当前名称的功能，所以暂时加s区分
func Zranges(key string, start int, stop int) ([]string, error) {
	c := pool.Get()
	defer c.Close()
	return redis.Strings(c.Do("ZRANGE", key, start, stop))
}

func ZrangeByScore(key string, start string, end string, offset int, count int) ([]string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.Strings(c.Do("ZRANGEBYSCORE", key, start, end, "WITHSCORES", "LIMIT", offset, count))
}

func Zrevrangebyscore(key string, offset int, count int) ([]string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.Strings(c.Do("ZREVRANGEBYSCORE", key, "+INF", "-INF", "WITHSCORES", "LIMIT", offset, count))
}

func Srandmember(key string, length int) ([]string, error) {
	c := pool.Get()
	defer c.Close()
	return redis.Strings(c.Do("SRANDMEMBER", key, length))
}

func Sismember(key string, member interface{}) (bool, error) {
	c := pool.Get()
	defer c.Close()
	return redis.Bool(c.Do("SISMEMBER", key, member))
}

// LIST
// --------------------------------------------------------------------------------

func Brpoplpush(src string, dest string, timeout int) (string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.String(c.Do("BRPOPLPUSH", src, dest, timeout))
}

func LRange(key string, start int, end int) ([]string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.Strings(c.Do("LRANGE", key, start, end))
}

func Lrem(key string, count int, val interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("LREM", key, count, val)
	return err
}

func Rpush(key string, val interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("RPUSH", key, val)
	return err
}

func Lpush(key string, val interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("LPUSH", key, val)
	return err
}

func Llen(key string) (int, error) {
	c := pool.Get()
	defer c.Close()

	res, err := c.Do("LLEN", key)
	if err != nil {
		return -1, err
	}
	return int(res.(int64)), nil
}

func Lpop(key string) ([]byte, error) {
	c := pool.Get()
	defer c.Close()

	res, err := c.Do("LPOP", key)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("EOF")
	}
	return res.([]byte), nil
}

func Rpop(key string) ([]byte, error) {
	c := pool.Get()
	defer c.Close()

	res, err := c.Do("RPOP", key)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("EOF")
	}
	return res.([]byte), nil
}

func Blpop(key string, timeout int) (interface{}, error) {
	return bpop("BLPOP", key, timeout)
}

func Brpop(key string, timeout int) (interface{}, error) {
	return bpop("BRPOP", key, timeout)
}

func bpop(cmd string, key string, timeout int) (interface{}, error) {
	c := pool.Get()
	defer c.Close()

	res, err := c.Do(cmd, key, timeout)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("EOF")
	}

	// Get value from list
	if list, ok := res.([]interface{}); ok {
		for i, value := range list {
			if i == 1 {
				return value, nil
			}
		}
	}
	return nil, errors.New("EOF")
}

// General Commands

func Do(cmd string, args ...interface{}) (reply interface{}, err error) {
	c := pool.Get()
	defer c.Close()

	return c.Do(cmd, args...)
}

// ------------------------------------------------------------------------

func newRedisPool(addr string, db string, password string, timeout int) *redis.Pool {

	// Set dial options
	// Specifies the timeout for connecting to the Redis server
	dialOptions := make([]redis.DialOption, 1)
	dialOptions[0] = redis.DialConnectTimeout(time.Duration(timeout) * time.Second)

	return &redis.Pool{
		MaxIdle:     80,
		MaxActive:   10000,
		IdleTimeout: 600 * time.Second,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", addr, dialOptions...)
			if err != nil {
				return nil, err
			}
			_, err = con.Do("AUTH", password)
			if err == nil {
				con.Do("SELECT", db)
			}
			return con, err
		},
	}
}

func bytesSlice(reply interface{}, err error) ([][]byte, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []interface{}:
		result := make([][]byte, len(reply))
		for i := range reply {
			if reply[i] == nil {
				continue
			}
			p, ok := reply[i].([]byte)
			if !ok {
				return nil, fmt.Errorf("redigo: Unexpected element type for []byte, got type %T", reply[i])
			}
			result[i] = p
		}
		return result, nil
	case nil:
		return nil, redis.ErrNil
	case redis.Error:
		return nil, reply
	}
	return nil, fmt.Errorf("redigo: Unexpected type for []byte, got type %T", reply)
}

func writeTo(data []byte, val reflect.Value) error {
	s := string(data)
	switch v := val; v.Kind() {

	// if we're writing to an interace value, just set the byte data
	// TODO: should we support writing to a pointer?
	case reflect.Interface:
		v.Set(reflect.ValueOf(data))

	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		ui, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(ui)

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)

	case reflect.String:
		v.SetString(s)

	case reflect.Slice:
		typ := v.Type()
		if typ.Elem().Kind() == reflect.Uint || typ.Elem().Kind() == reflect.Uint8 || typ.Elem().Kind() == reflect.Uint16 || typ.Elem().Kind() == reflect.Uint32 || typ.Elem().Kind() == reflect.Uint64 || typ.Elem().Kind() == reflect.Uintptr {
			v.Set(reflect.ValueOf(data))
		}
	}
	return nil
}

func writeToContainer(data [][]byte, val reflect.Value) error {
	switch v := val; v.Kind() {
	case reflect.Ptr:
		return writeToContainer(data, reflect.Indirect(v))
	case reflect.Interface:
		return writeToContainer(data, v.Elem())
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return errors.New("redigo: Invalid map type")
		}
		elemtype := v.Type().Elem()
		for i := 0; i < len(data)/2; i++ {
			mk := reflect.ValueOf(string(data[i*2]))
			mv := reflect.New(elemtype).Elem()
			writeTo(data[i*2+1], mv)
			v.SetMapIndex(mk, mv)
		}
	case reflect.Struct:
		for i := 0; i < len(data)/2; i++ {
			name := string(data[i*2])
			field := v.FieldByName(name)
			if !field.IsValid() {
				continue
			}
			writeTo(data[i*2+1], field)
		}
	default:
		return errors.New("redigo: Invalid container type")
	}
	return nil
}
