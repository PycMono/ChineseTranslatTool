package main

import (
	"fmt"
	"moqikaka.com/ChineseTranslatTool/src/bll"
	"moqikaka.com/goutil/logUtil"
	"sync"
)

var (
	wg            sync.WaitGroup
	con_SEPERATOR = "------------------------------------------------------------------------------"
)

func init() {
	// 设置WaitGroup需要等待的数量，只要有一个服务器出现错误都停止服务器
	wg.Add(1)

	// 设置日志文件的存储目录
	logUtil.SetLogPath("LOG")
}

func main() {
	err := bll.TranslateWord() // 翻译
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("翻译完成")
		logUtil.InfoLog("翻译完成")
	}

	wg.Wait()
}
