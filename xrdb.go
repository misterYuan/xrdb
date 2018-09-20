/*redis-db*/
package xrdb

import (
	"encoding/json"
	"log"

	"github.com/gomodule/redigo/redis"
)

var err_nil_value = "redigo: nil returned"

func WithConn(task func(conn redis.Conn) error) error {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	defer conn.Close()
	return task(conn)
}

func Set(key, value, expire interface{}) error {
	return WithConn(func(conn redis.Conn) error {
		_, err := conn.Do("SET", key, value, "EX", expire)
		return err
	})
}

func Get(key string) (value string, ok bool) {
	var err error
	err = WithConn(func(conn redis.Conn) error {
		value, err = redis.String(conn.Do("GET", key))
		return err
	})
	if err != nil && err.Error() == err_nil_value {
		return "", false
	}
	if err != nil && err.Error() != err_nil_value {
		log.Println(err)
		panic(err.Error())
	}
	return value, true
}

func Exists(key string) bool {
	var ok bool
	var err error
	err = WithConn(func(conn redis.Conn) error {
		ok, err = redis.Bool(conn.Do("EXISTS", key))
		return err
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return ok
}

func Del(key string) {
	err := WithConn(func(conn redis.Conn) error {
		_, err := conn.Do("DEL", key)
		return err
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

/*
列表
*/
func LPush(lName string, value interface{}, expire string) error {
	if err := WithConn(func(conn redis.Conn) error {
		_, err := conn.Do("LPUSH", lName, value)
		return err
	}); err != nil {
		return err
	}
	return WithConn(func(conn redis.Conn) error {
		_, err := conn.Do("EXPIRE", lName, expire)
		return err
	})
}

func LRangeAll(lName string) [][]byte {
	var bs [][]byte
	var err error
	err = WithConn(func(conn redis.Conn) error {
		bs, err = redis.ByteSlices(conn.Do("lrange", lName, 0, -1))
		return err
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return bs
}

/*
将对象序列化存储
*/
func SetI(key, i, expire interface{}) error {
	bs, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = WithConn(func(conn redis.Conn) error {
		_, err := conn.Do("SET", key, bs, "EX", expire)
		return err
	})
	return err
}

func GetI(key, i interface{}) bool {
	var err error
	var bs []byte
	err = WithConn(func(conn redis.Conn) error {
		bs, err = redis.Bytes(conn.Do("GET", key))
		return err
	})

	if err != nil && err.Error() == err_nil_value {
		return false
	}
	if err != nil && err.Error() != err_nil_value {
		log.Println(err)
		panic(err.Error())
	}
	if err := json.Unmarshal(bs, i); err != nil {
		log.Println(err)
		panic(err.Error())
	}
	return true
}

/*
hashmap
*/
func HMSet(key, field, value, expire string) error {
	if err := WithConn(func(conn redis.Conn) error {
		_, err := conn.Do("HMSET", key, field, value)
		return err
	}); err != nil {
		return err
	}
	return WithConn(func(conn redis.Conn) error {
		_, err := conn.Do("EXPIRE", key, expire)
		return err
	})
}

func HMGet(key, field string) (value string, ok bool) {
	var err error
	err = WithConn(func(conn redis.Conn) error {
		value, err = redis.String(conn.Do("HGET", key, field))
		return err
	})
	if err != nil && err.Error() == err_nil_value {
		return "", false
	}
	if err != nil && err.Error() != err_nil_value {
		log.Println(err)
		panic(err.Error())
	}
	return value, true
}

func HDel(key, field string) {
	WithConn(func(conn redis.Conn) error {
		conn.Do("HDEL", key, field)
		return nil
	})
}
