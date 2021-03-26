package main

import (
	"github.com/gin-gonic/gin"
	"goroutine"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func init() {

}
func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {

		c.JSON(200, "ok")
		go func() {
			for {

				// Most cost of heap is timer.
				time.Sleep(1 * time.Nanosecond)
				goroutine.ProtectedGo(func() {

				}, goroutine.GoParam{
					ExpectedExpireSecond: 10,
					ShouldProtected:      true,
				})
			}
		}()

	})
	go func() {
		if e := http.ListenAndServe(":6060", nil); e != nil {
			panic(e)
		}
	}()

	if e := r.Run(":8080"); e != nil {
		panic(e)
	}
}
