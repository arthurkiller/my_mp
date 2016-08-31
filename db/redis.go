package redism

import (
	"hash/crc32"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Redism interface {
	Get(name, opt, key string) redis.Conn
	Close() error
}

type redism struct {
	redisPoolMaster map[string]([]*redis.Pool)
	redisPoolSlave  map[string]([][]*redis.Pool)
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
			return c, err
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
	for name, _ := range conf.Masters { //this for-loop get the server name to identifi different redis
		rm.redisPoolMaster[name] = make([]*redis.Pool, len(conf.Masters[name]))
		for i, v := range conf.Masters[name] { // this for-loop
			rm.redisPoolMaster[name][i] = builder(v, conf)
			rm.redisPoolSlave[name][i] = make([]*redis.Pool, len(conf.Slaves[name][i]))
			for j, vv := range conf.Slaves[name][i] {
				rm.redisPoolSlave[name][i][j] = builder(vv, conf)
			}
		}
	}
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
