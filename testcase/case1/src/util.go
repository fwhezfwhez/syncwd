package src

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

func NewPool(server, password string, db int) *redis.Pool {
	//if password == "" {
	//	panic("pw empty")
	//}
	return &redis.Pool{
		MaxIdle:     200,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				fmt.Printf("occur error at newPool Dial: %v\n", err)
				return nil, err
			}
			_, e := c.Do("ping")
			if e != nil {
				if password != "" {
					if _, err := c.Do("AUTH", password); err != nil {
						c.Close()
						fmt.Printf("occur error at newPool Do Auth: %v\n", err)
						return nil, err
					}
					if _, e := c.Do("ping"); e != nil {
						return nil, fmt.Errorf("ping twice err: %v", e)
					}
				}
			}

			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				fmt.Printf("occur error at newPool Do SELECT: %v\n", err)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}