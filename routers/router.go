package routers

import (
	"hs_pl/apis/v1"
	"github.com/gin-gonic/gin"
)

var testCtl = new(v1.TestController)

// a global routers
var Router *gin.Engine

// Init routers, adding paths to it.
func init() {
	Router = gin.Default()
	// api group for v1
	v1Group := Router.Group("/api")
	{
		testGroup :=v1Group.Group("data")
		{
			testGroup.GET("/theme-list", testCtl.Async)
			testGroup.GET("/lock", testCtl.Lock)
		}
	}
}