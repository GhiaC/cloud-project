package main

import (
	"github.com/labstack/echo/v4"
	"gogs.ghiasi.me/masoud/cloud-project/config"
	handler2 "gogs.ghiasi.me/masoud/cloud-project/handler"
	"net/http"
	"strconv"
)

func main() {
	config.Initialize("")
	e := echo.New()
	handler := handler2.NewHandler()
	handler.Start()
	e.GET("/cpu", func(c echo.Context) error {
		handler.IncreaseLoad(10000)
		return c.String(http.StatusOK, "Cpu:"+handler.GetCpuAvg())
	})
	e.GET("/mem", func(c echo.Context) error {
		handler.IncreaseMemUsage(20)
		return c.String(http.StatusOK, "Mem:"+handler.GetMemAvg())
	})
	e.GET("/cpumem", func(c echo.Context) error {
		handler.IncreaseMemUsage(20)
		handler.IncreaseLoad(10000)
		return c.String(http.StatusOK, "Mem:"+handler.GetMemAvg()+" cpu:"+handler.GetCpuAvg())
	})
	e.GET("/GetCpuUsage", func(c echo.Context) error {
		return c.String(http.StatusOK, handler.GetCpuAvg())
	})
	e.GET("/GetMemUsage", func(c echo.Context) error {
		return c.String(http.StatusOK, handler.GetMemAvg())
	})
	e.GET("/WhatsYourName", func(c echo.Context) error {
		return c.String(http.StatusOK, config.Conf.Core.ServerName)
	})
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Conf.Core.Port)))
}
