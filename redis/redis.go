package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type Config struct {
	Protocol        string
	Host            string
	MaxActive       int
	MaxIdle         int
	IdleTimeout     time.Duration
	MaxConnLifetime time.Duration
}

var pool *redis.Pool

func Init(config *Config) {
	pool = &redis.Pool{
		MaxIdle:         config.MaxIdle,
		MaxActive:       config.MaxActive,
		IdleTimeout:     config.IdleTimeout,
		MaxConnLifetime: config.MaxConnLifetime,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(config.Protocol, config.Host)
		},
	}
}

func Close() error {
	return pool.Close()
}

type (
	Redis struct {
		conn redis.Conn
	}
)

var ErrNil = redis.ErrNil

func GetRedis() (*Redis, error) {
	conn := pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	return &Redis{
		conn: conn,
	}, nil
}

func (r *Redis) Close() error {
	return r.conn.Close()
}

func (r *Redis) Err() error {
	return r.conn.Err()
}

func (r *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return r.conn.Do(commandName, args...)
}

func (r *Redis) Send(commandName string, args ...interface{}) error {
	return r.conn.Send(commandName, args...)
}

func (r *Redis) Flush() error {
	return r.conn.Flush()
}

func (r *Redis) Receive() (reply interface{}, err error) {
	return r.conn.Receive()
}

func Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

func Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

func Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

func Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

func String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

func Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}

func Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

func Values(reply interface{}, err error) ([]interface{}, error) {
	return redis.Values(reply, err)
}

func Float64s(reply interface{}, err error) ([]float64, error) {
	return redis.Float64s(reply, err)
}

func Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

func ByteSlices(reply interface{}, err error) ([][]byte, error) {
	return redis.ByteSlices(reply, err)
}

func Ints(reply interface{}, err error) ([]int, error) {
	return redis.Ints(reply, err)
}

func StringMap(result interface{}, err error) (map[string]string, error) {
	return redis.StringMap(result, err)
}

func IntMap(result interface{}, err error) (map[string]int, error) {
	return redis.IntMap(result, err)
}

func Int64Map(result interface{}, err error) (map[string]int64, error) {
	return redis.Int64Map(result, err)
}

func Positions(result interface{}, err error) ([]*[2]float64, error) {
	return redis.Positions(result, err)
}
