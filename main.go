// @title COZE-DISCORD-PROXY
// @version 1.0.0
// @description COZE-DISCORD-PROXY 代理服务
// @BasePath
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"gpltrans/router"
	"gpltrans/utils"
)

func main() {
	parser := utils.NewParallelParser("begin select 1 from dual; end;")
	_, err := parser.Parse()
	fmt.Println(err)

	server := gin.Default()

	router.InitApiRouter(server)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic("failed to start HTTP server: " + err.Error())
		}
	}()

	// 等待中断信号
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// 给 HTTP 服务器一些时间来关闭
	ctxShutDown, cancelShutDown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutDown()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		panic("HTTP server Shutdown failed:" + err.Error())
	}
}
