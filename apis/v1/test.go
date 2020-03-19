package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/idoubi/goz"
	"hs_pl/lib/redisLib"
	"log"
	"strconv"
	"time"
)

const (
	// 每分钟刷新主题缓存
	FLUSHTIME = 1 * 60
)

type TestController struct {}

func (test *TestController) Async(c *gin.Context) {

	input := c.DefaultQuery("input", "")
	interval := c.DefaultQuery("interval", "")
	size := c.DefaultQuery("size", "")

	redisClient := redisLib.GetClient()
	data, err:= redisClient.HGet("h:test", "content").Result()

	if err != nil {
		fmt.Print("缓存呢？？？")
	}

	lastCacheTime, _ := redisClient.HGet("h:test", "time").Result()
	lastTime, _ := strconv.Atoi(lastCacheTime)
	currentTime := uint64(time.Now().Unix())

	//如果请求时间超出上一次缓存时间1min，则刷新缓存内容
	if currentTime - FLUSHTIME >= uint64(lastTime) {
		go func() {
			cli := goz.NewClient()
			resp, err := cli.Get("https://gw.datayes.com/rrp_adventure/mobile/whitelist/theme", goz.Options{
				Query: map[string]interface{}{
					"input": input,
					"interval": interval,
					"size": size,
				},
			})
			if err != nil {
				log.Fatalln(err)
			}

			body, err := resp.GetBody()
			contents := body.GetContents()

			if len(contents) > 0 {
				log.Print("打印当前时间")
				log.Print(currentTime)
				redisLib.GetClient().HSet("h:test", "content", contents).Err()
				redisLib.GetClient().HSet("h:test", "time", currentTime).Err()
			}

			log.Print("已更新缓存")
		}()
	}

	c.JSON(200, gin.H{
		"code":      200,
		"data":      data,
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}


