package main

import (
	"embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
)

//go:embed frontend/dist/*
var FS embed.FS

// 关闭主线程，ui会自动退出；退出ui，主线程会自动退出；
// 启动gin 携程会快很多
func main() {
	go func() { // gin协程，是平行的关系；等待用户请求；
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		// func(c *gin.context) {}   c相当于c语言的指针，指向gin的上下文
		router.GET("/", func(c *gin.Context) {
			// 与nodejs不同，go语言使用的是流的方式，使用writer的形式制造response
			c.String(http.StatusOK, "<h1>hello world</h1>")
		})
		router.Run(":8888")
	}()
	//以下为主线程
	// open chrome on mac using go
	// command := "open"
	// // arg1 := "-a"
	// chromePath := "/Applications/Google\\ Chrome.app"
	// cmd := exec.Command(command, "-e", chromePath)
	// cmd.Start()
	// cmd.Process.Kill()
	// select {}
	// 声明ui变量，其类型为lorca.ui
	var ui lorca.UI
	// 初始值为nil
	// fmt.Printf("%v\n", ui)
	//locra.New 返回两个值，忽略第二个值
	ui, _ = lorca.New("http://127.0.0.1:8888/", "", 800, 600, "--disable-sync", "--disable-translate")
	// 创建一个频道，监听操作系统的信号。
	chSingal := make(chan os.Signal, 1)
	// 通知，中断和终止信号。
	signal.Notify(chSingal, syscall.SIGINT, syscall.SIGTERM)
	//channel筛选，有一个成功，就继续运行；会阻塞当前线程，等case出现。
	// <-chSignal
	///select，等待n个频道（需要自己创建），随机等直到channel出现值
	select {
	case <-ui.Done():
	case <-chSingal:
	}
	ui.Close()
}
