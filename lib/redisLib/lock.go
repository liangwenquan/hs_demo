package redisLib

import (
	"errors"
	"time"
	"github.com/satori/go.uuid"
	"github.com/go-redis/redis"
	"fmt"
)

func GetLock(lockName string, acquireTimeout, lockTimeOut time.Duration) (string, error) {
	code := uuid.NewV4().String()
	endTime := time.Now().Add(acquireTimeout).UnixNano()

	for time.Now().UnixNano() <= endTime {
		if success, err := client.SetNX(lockName, code, lockTimeOut).Result(); err !=nil && err != redis.Nil{
			return "", err
		} else if success {
			return code, nil
		} else if client.TTL(lockName).Val() == -1 {
			client.Expire(lockName, lockTimeOut)
		}
	}
	return "", errors.New("lock time out")
}

func ReleaseLock(lockName, code string) bool {
	txf := func(tx *redis.Tx) error {
		if v, err := tx.Get(lockName).Result(); err != nil && err != redis.Nil {
			return err
		} else if v == code {
			_, err := tx.Pipelined(func(pipe redis.Pipeliner) error {
				pipe.Del(lockName)
				return nil
			})
			return err
		}
		return nil
	}

	for {
		if err := client.Watch(txf, lockName); err == nil {
			return true
		} else if err == redis.TxFailedErr {
			fmt.Println("watch key is modified, retry to release lock. err:", err.Error())
		} else {
			fmt.Println("err:", err.Error())
			return false
		}
	}
}