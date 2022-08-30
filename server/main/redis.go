package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

//定义一个全家的pool
var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTime time.Duration) {
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) { //初始化代码连接，连接到那个redis
			return redis.Dial("tcp", address)
		},
		MaxIdle:     maxIdle,   //最大空闲连接数
		MaxActive:   maxActive, //表示和数据库的最大连接数，0表示没有限制
		IdleTimeout: idleTime,  //最大空闲数
	}
}
