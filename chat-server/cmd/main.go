package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat-server/internal/app"
	"chat-server/internal/http"
	"chat-server/internal/socket"
)

func main() {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 初始化
	app.Init()

	// 启动服务
	go func() {
		err := http.NewHttp().Start()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := socket.NewSocket().Start()
		if err != nil {
			panic(err)
		}
	}()

EXIT:
	for {
		sig := <-sc

		// 信号处理
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}
	time.Sleep(time.Second)
	os.Exit(state)
}
