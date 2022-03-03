package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

var GB = float32(1 << 30)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/memory", func(c *gin.Context) {
		memory, err := memory.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"memory_total":  float32(memory.Total) / GB,
			"memory_used":   float32(memory.Used) / GB,
			"memory_cached": float32(memory.Cached) / GB,
			"memory_free":   float32(memory.Free) / GB,
		})
	})

	r.GET("/cpu", func(c *gin.Context) {
		before, err := cpu.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		time.Sleep(time.Duration(1) * time.Second)
		after, err := cpu.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		total := float64(after.Total - before.Total)
		c.JSON(http.StatusOK, gin.H{
			"cpu_user":   float64(after.User-before.User) / total * 100,
			"cpu_system": float64(after.System-before.System) / total * 100,
			"cpu_idle":   float64(after.Idle-before.Idle) / total * 100,
		})
	})

	return r
}
func main() {
	r := setupRouter()
	r.Run(":8080")
}
