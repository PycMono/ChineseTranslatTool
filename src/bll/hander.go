package bll

import (
	"fmt"
	"github.com/pkg/errors"
	"moqikaka.com/ChineseTranslatTool/src/bll/excelUtil"
	"moqikaka.com/ChineseTranslatTool/src/model"
	"moqikaka.com/goutil/fileUtil"
	"moqikaka.com/goutil/logUtil"
	"path"
)

// 翻译文字
// 参数：无
// 返回值：
// 1.错误对象
func TranslateWord() error {
	baseConfig := config.GetBaseConfig()
	if baseConfig == nil {
		return errors.New("baseConfig为nil")
	}

	// 获取译文
	worldDict, err := getTranslateWordDict(baseConfig)
	if err != nil {
		return err
	}

	// 获取源目录下的所有文件
	filePathList, err := fileUtil.GetFileList(baseConfig.SourceFilePath)
	if err != nil {
		return err
	}

	notInExcelList := make([]string, 0)
	// 遍历文件处理
	for _, filePath := range filePathList {
		tempExcelUtil, err := excelUtil.NewExcelUtil(filePath)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			resultList, err := tempExcelUtil.Replace(worldDict, baseConfig.TarFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}

			notInExcelList = append(notInExcelList, resultList...)
		}
	}

	for _, item := range notInExcelList {
		fmt.Println(item)
		logUtil.ErrorLog(item)
	}

	return nil
}

// 获取译文
// 参数：无
// 返回值：
// 1.译文数据
// 2.错误对象
func getTranslateWordDict(baseConfig *config.BaseConfig) (map[string]string, error) {
	curPath := fileUtil.GetCurrentPath()
	translateFile, err := excelUtil.NewExcelUtil(path.Join(curPath, baseConfig.TranslateFileName))
	if err != nil {
		return nil, err
	}

	return translateFile.ReadExcelToDict(), nil
}
