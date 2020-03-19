package main

import (
	"fmt"
	"github.com/spf13/viper"
	"hs_pl/lib/redisLib"
	"hs_pl/routers"
	"net/http"
)

func init()  {
	redisLib.NewClient()
	initConfig()
}

func main() {
	// get the global routers from routers.go
	router := routers.Router

	// set up a http server
	server := http.Server{
		Addr: ":9088",
		Handler: router,
		MaxHeaderBytes: 1 << 20,
	}

	// run the server
	server.ListenAndServe()
}

func initConfig()  {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	fmt.Println("config redis:", viper.Get("redis"))
}