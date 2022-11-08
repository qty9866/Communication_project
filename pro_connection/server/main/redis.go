package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// 定义一个全局的pool
var pool *redis.Pool

// 初始化
func InitPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,     //最大空闲链接数
		MaxActive:   maxActive,   // 表示和数据库的最大连接数，0表示没有限制
		IdleTimeout: idleTimeout, //最大空闲时间
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}
