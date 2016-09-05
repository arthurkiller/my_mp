package redism

import (
	"fmt"
	"hash/crc32"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Redism interface {
	Get(name, opt, key string) redis.Conn
	GetScript(key string) *redis.Script
	Close() error
}

type redism struct {
	redisPoolMaster map[string]([]*redis.Pool)
	redisPoolSlave  map[string]([][]*redis.Pool)
	scripts         map[string]*redis.Script
}
type RedismConf struct {
	Maxactive   int
	Maxidle     int
	Idletimeout int
	//storge the address in the array
	Masters map[string]([]string)
	Slaves  map[string]([][]string)
}

func builder(address string, conf RedismConf) *redis.Pool {
	return &redis.Pool{
		MaxActive:   conf.Maxactive,
		MaxIdle:     conf.Maxidle,
		IdleTimeout: time.Duration(conf.Idletimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func NewRedism(conf RedismConf) Redism {
	rm := new(redism)
	rm.redisPoolMaster = make(map[string]([]*redis.Pool))
	rm.redisPoolSlave = make(map[string]([][]*redis.Pool))

	for name := range conf.Masters { //this for-loop get the server name to identifi different redis
		rm.redisPoolMaster[name] = make([]*redis.Pool, len(conf.Masters[name]))

		for i, v := range conf.Masters[name] { // this for-loop used to build the master redis pool
			rm.redisPoolMaster[name][i] = builder(v, conf)
			rm.redisPoolSlave[name] = make([][]*redis.Pool, len(conf.Slaves[name]))

			for j, vv := range conf.Slaves[name][i] {
				rm.redisPoolSlave[name][i] = make([]*redis.Pool, len(conf.Slaves[name][i]))
				rm.redisPoolSlave[name][i][j] = builder(vv, conf)
			}
		}
	}

	rm.scripts = make(map[string]*redis.Script, 4)
	//TODO: add the lua script here
	rm.scripts["uid-check-fans"] = redis.NewScript(2, fmt.Sprintf(`
		local count = 1
		local length = tonumber(KEYS[2])
		local result = {}
		for i = 1, length do
		local val = tonumber(redis.call("SISMEMBER", tostring(KEYS[1]),tostring(ARGV[i])))
			if val == 1 then
			result[count] = tostring(ARGV[i])
			count = count + 1
			end
		end
		return result
	`))

	rm.scripts["newsid-get-news"] = redis.NewScript(1, fmt.Sprintf(`
		local length = tonumber(KEYS[1])
		local result = {}
			for i = 1, length do
			local val = redis.call("HGETALL",tostring(ARGV[i]))
			result[i] = val
			end
		return result
	`))

	rm.scripts["uid-put-news"] = redis.NewScript(2, fmt.Sprintf(`
		local count = 1
		local length = tonumber(KEYS[2])
		locla result = {}
		for i = 1, length do
			redis.call("LPUSH",ARGV[i],KEYS[1])
		end
	`))

	return rm
}

func (r *redism) Get(name, opt, key string) redis.Conn {
	if len(r.redisPoolSlave[name]) == 0 {
		opt = "W"
	}
	if opt == "R" {
		hash := crc32.ChecksumIEEE([]byte(key))
		i := hash % uint32(len(r.redisPoolMaster[name]))
		index := hash % uint32(len(r.redisPoolSlave[name][i]))
		c := r.redisPoolSlave[name][i][index].Get()
		return c
	}

	hash := crc32.ChecksumIEEE([]byte(key))
	i := hash % uint32(len(r.redisPoolMaster[name]))
	c := r.redisPoolMaster[name][i].Get()
	return c
}

func (r *redism) GetScript(key string) *redis.Script {
	return r.scripts[key]
}

func (r *redism) Close() error {
	var err error
	for _, v := range r.redisPoolMaster {
		for _, vv := range v {
			err = vv.Close()
		}
	}
	for _, v := range r.redisPoolSlave {
		for _, vv := range v {
			for _, vvv := range vv {
				err = vvv.Close()
			}
		}
	}
	return err
}
