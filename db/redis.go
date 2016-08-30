package redism

import (
	"hash/crc32"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Redism interface {
	Get(opt, key string) redis.Conn
	Close() error
}

type redism struct {
	redisPoolMaster []*redis.Pool
	redisPoolSlave  [][]*redis.Pool
}
type RedismConf struct {
	Maxactive   int
	Maxidle     int
	Idletimeout int
	//storge the address in the array
	Masters []string
	Slaves  [][]string
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
	redism := new(redism)
	redism.redisPoolMaster = make([]*redis.Pool, len(conf.Masters))
	redism.redisPoolSlave = make([]([]*redis.Pool), len(conf.Masters))

	for i, v := range conf.Masters {
		redism.redisPoolMaster[i] = builder(v, conf)
		redism.redisPoolSlave[i] = make([]*redis.Pool, len(conf.Slaves[i]))
		for j, vv := range conf.Slaves[i] {
			redism.redisPoolSlave[i][j] = builder(vv, conf)
		}
	}

	return redism
}

func (r *redism) Get(opt, key string) redis.Conn {
	if len(r.redisPoolSlave) == 0 {
		opt = "W"
	}
	if opt == "R" {
		hash := crc32.ChecksumIEEE([]byte(key))
		i := hash % uint32(len(r.redisPoolMaster))
		index := hash % uint32(len(r.redisPoolSlave[i]))
		c := r.redisPoolSlave[i][index].Get()
		return c
	}

	hash := crc32.ChecksumIEEE([]byte(key))
	i := hash % uint32(len(r.redisPoolMaster))
	c := r.redisPoolMaster[i].Get()
	return c
}

func (r *redism) Close() error {
	var err error = nil
	for _, v := range r.redisPoolMaster {
		err = v.Close()
	}
	for _, v := range r.redisPoolSlave {
		for _, vv := range v {
			err = vv.Close()
		}
	}
	return err
}
