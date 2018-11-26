package config

import (
	"fmt"
	"moqikaka.com/goutil/configUtil"
	"moqikaka.com/goutil/debugUtil"
	"moqikaka.com/goutil/logUtil"
	"strings"
)

// 基础配置对象
type BaseConfig struct {
	// 源Excel目录
	SourceFilePath string

	// 目标Excel目录
	TarFilePath string

	// 译文
	TranslateFileName string

	// 不翻译的表名字
	NotTableDict map[string]string
}

//----------------------------------------------------------------

var (
	baseConfig *BaseConfig
)

func init() {
	// 读取配置文件内容
	configObj := configUtil.NewXmlConfig()
	err := configObj.LoadFromFile("config.xml")
	if err != nil {
		fmt.Println(fmt.Sprintf("初始化NewXmlConfig报错%s", err))
		logUtil.ErrorLog(fmt.Sprintf("初始化NewXmlConfig报错%s", err))
		return
	}

	initBaseConfig(configObj)
}

// 初始化配置
// 参数：
// config：工具助手
// 返回值：
// 1.错误对象
func initBaseConfig(config *configUtil.XmlConfig) error {
	sourceFilePath, err := config.String("root/BaseConfig/SourceFilePath", "")
	if err != nil {
		return err
	}

	tarFilePath, err := config.String("root/BaseConfig/TarFilePath", "")
	if err != nil {
		return err
	}

	translateFileName, err := config.String("root/BaseConfig/TranslateFileName", "")
	if err != nil {
		return err
	}

	notTable, err := config.String("root/BaseConfig/NotTable", "")
	// 去除空格,制表符，换行符号，防止切割后，出现乱七八糟的问题，导致查询时候数据不在字典里面
	notTable = strings.Replace(notTable, "\r", "", -1)
	notTable = strings.Replace(notTable, "\t", "", -1)
	notTable = strings.Replace(notTable, "\n", "", -1)
	if err != nil {
		return err
	}

	notTableDict := make(map[string]string)
	var s string
	for _, table := range strings.Split(notTable, ";") {
		notTableDict[table] = table
		if s == "" {
			s = table
			continue
		}

		s = fmt.Sprintf("%s,%s", s, table)
	}

	fmt.Println(len(notTableDict))
	logUtil.DebugLog(fmt.Sprintf("不翻译的表名为：%s", s))
	fmt.Sprint(fmt.Sprintf("不翻译的表名为：%s", s))

	baseConfig = &BaseConfig{
		SourceFilePath:    sourceFilePath,
		TarFilePath:       tarFilePath,
		TranslateFileName: translateFileName,
		NotTableDict:      notTableDict,
	}

	debugUtil.Println("BaseConfig:", baseConfig)

	return nil
}

func GetBaseConfig() *BaseConfig {
	return baseConfig
}
