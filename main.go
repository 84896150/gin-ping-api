package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ping/ping"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	router := gin.Default()
	s := &http.Server{
		Addr:           ":8888",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	v1 := router.Group("/ping")
	{
		//v1.POST("/", CreateUser)
		//v1.GET("/", FetchAllUsers)
		v1.GET("/:ip", FetchSingleUser)
		//v1.PUT("/:id", UpdateUser)
		//v1.DELETE("/:id", DeleteUser)
	}
	s.ListenAndServe()
}

func FetchSingleUser(c *gin.Context) {
	host := c.Param("ip")
	pinger, err := ping.NewPinger(host)
	if err != nil {
		fmt.Println("错误")
		var result gin.H
		result = gin.H{
			"result": nil,
			"Status": "ip地址有误",
			"Addr":   host,
		}
		c.JSON(http.StatusOK, result)
		return
	}
	pinger.OnRecv = func(pkt *ping.Packet) {
	}
	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
	}

	pinger.Count = 5
	//pinger.Interval = 1 * time.Second
	pinger.Timeout = 5 * time.Second
	pinger.SetPrivileged(true)
	fmt.Printf("开始PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	//ping [-c count] [-i interval] [-t timeout] [--privileged] host
	err = pinger.Run()
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	fmt.Println(stats)
	var status string = "离线"
	if stats.PacketsRecv > 0 {
		status = "在线"
	}
	var result gin.H
	if err != nil {
		// If no results send null
		result = gin.H{
			"result":                nil,
			"Status":                status,
			"Addr":                  stats.Addr,
			"PacketsSent":           stats.PacketsSent,
			"PacketsRecv":           stats.PacketsRecv,
			"PacketsRecvDuplicates": stats.PacketsRecvDuplicates,
			"PacketLoss":            stats.PacketLoss,
			"MinRtt":                stats.MinRtt,
			"MaxRtt":                stats.MaxRtt,
			"AvgRtt":                stats.AvgRtt,
			"StdDevRtt":             stats.StdDevRtt,
		}
	} else {
		result = gin.H{
			"Status":                status,
			"Addr":                  stats.Addr,
			"PacketsSent":           stats.PacketsSent,
			"PacketsRecv":           stats.PacketsRecv,
			"PacketsRecvDuplicates": stats.PacketsRecvDuplicates,
			"PacketLoss":            stats.PacketLoss,
			"MinRtt":                stats.MinRtt,
			"MaxRtt":                stats.MaxRtt,
			"AvgRtt":                stats.AvgRtt,
			"StdDevRtt":             stats.StdDevRtt,
		}
	}
	c.JSON(http.StatusOK, result)
}
