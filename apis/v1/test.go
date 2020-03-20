package v1

import (
	"crypto/md5"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/idoubi/goz"
	"hs_pl/lib/redisLib"
	"log"
	"net/url"
	"sync"
	"time"
	"fmt"
)

type TestController struct {}

func (test *TestController) Async(c *gin.Context) {
	paramUri := c.Query("rq_uri")
	themeUri, _:= url.QueryUnescape(paramUri)

	lockPre := "lock_theme_pre"
	cacheKey := strMd5(themeUri)

	redisClient := redisLib.GetClient()
	data, err:= redisClient.Get(cacheKey).Result()

	var dataBody []byte
	var dataMap map[string]interface{}

	dataBody = []byte(data)
	ttl, _ := redisClient.TTL("hello").Result()

	if err != nil {
		ttl = 3600 * 1e9;
		dataBody = fetchGwTheme(lockPre, cacheKey, themeUri)
	}

	if err := json.Unmarshal(dataBody, &dataMap); err != nil {
		log.Print(err)
		log.Print("parse to map 失败")
	}

	if uint64(ttl)/1e9 <= 1000 {
		log.Print("不晓事")
		go fetchGwTheme(lockPre, cacheKey, themeUri)
	}

	c.JSON(200, gin.H{
		"code":      200,
		"data":      dataMap,
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}

func (test *TestController) Lock(c *gin.Context) {
	go func() {
		lockName := "lock:test"
		acquireTimeOut := 9 * time.Minute
		lockTimeOut := 4 * time.Second

		count := 10
		var wg sync.WaitGroup
		for count > 0 {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				s, e := redisLib.GetLock(lockName, acquireTimeOut, lockTimeOut)
				log.Print(id,s,e)
				r := redisLib.ReleaseLock(lockName, s)
				log.Print(r)
			}(count)
			count--
		}
		wg.Wait()
	}()

	c.JSON(200, gin.H{
		"code":      200,
		"data":      "lock",
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}

func fetchGwTheme(lockPre string, cacheKey string, url string) []byte {
	cli := goz.NewClient()
	resp, err := cli.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := resp.GetBody()
	contents := body.GetContents()

	if len(contents) > 0 {
		wLock := lockPre + cacheKey
		success, err := redisLib.GetClient().SetNX(wLock, cacheKey, 5 * 1e9).Result()

		if err != nil {
			log.Print(err)
		}
		if success {
			redisLib.GetClient().Set(cacheKey, contents, 3600 * 1e9).Err()
			log.Print("已更新缓存")
		}
	}

	return body
}

func strMd5(str string) string {
	data := []byte(str)
	hash := md5.Sum(data)
	md5str := fmt.Sprintf("%x", hash)

	return md5str
}


