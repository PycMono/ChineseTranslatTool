package excelUtil

import (
	"fmt"
	"path"

	"regexp"

	"strings"

	"unicode"

	"errors"
	"github.com/tealeg/xlsx"
	"moqikaka.com/ChineseTranslatTool/src/model"
	"moqikaka.com/Test/src/BLL/fileUtil"
	"moqikaka.com/goutil/logUtil"
	"os"
)

type ExcelUtil struct {
	// 文件路径
	path string

	// xlsx文件对象
	file *xlsx.File
}

// Excel成字典的形式(每个cell的奇数作为key，偶数作为value)
// 参数：无
// 返回值：
// 字典集合
func (this *ExcelUtil) ReadExcelToDict() map[string]string {
	resultDict := make(map[string]string)
	for _, sheet := range this.file.Sheets {
		for _, row := range sheet.Rows {
			i := 0
			key := ""
			for _, cell := range row.Cells {
				i++
				text := cell.String()
				// 每个cell的奇数作为key，偶数作为value
				if i%2 == 0 {
					resultDict[key] = text
				} else {
					key = text
				}
			}
		}
	}

	return resultDict
}

// 替换字符串
// 参数：
// resultDict：译文
// tarFilePath：保存到目标文件
// 返回值：
// 1.错误对象
func (this *ExcelUtil) Replace(resultDict map[string]string, tarFilePath string) ([]string, error) {
	if !fileUtil.DirExists(tarFilePath) {
		if err := os.MkdirAll(tarFilePath, os.ModePerm|os.ModeTemporary); err != nil {
			return nil, err
		}
	}

	resultList := make([]string, 0)
	baseConfig := config.GetBaseConfig()
	for _, sheet := range this.file.Sheets {
		// 判断当前表是否不翻译
		if _, exists := baseConfig.NotTableDict[sheet.Name]; exists {
			continue
		}

		if sheet.Name == "b_goods_special_enum" {
			if _, exists := baseConfig.NotTableDict[sheet.Name]; !exists {
				for key, value := range baseConfig.NotTableDict {
					logUtil.DebugLog(fmt.Sprintf("字典的key：--%s--，value：--%s--", key, value))
				}

				logUtil.DebugLog(fmt.Sprintf("长度--%d--", len(baseConfig.NotTableDict)))
				logUtil.DebugLog(fmt.Sprintf("居然不存在--%s--", sheet.Name))
			}
		}

		for index, row := range sheet.Rows {
			if index == 0 { // 表头不替换
				continue
			}

			for _, cell := range row.Cells {
				text := cell.String()
				if !this.isChinese(text) {
					continue
				}

				value, exists := resultDict[text]
				if !exists {
					warnInfo := fmt.Sprintf("resultDict不存在key=------%s------的日文翻译sheet.Name=%s,row=%d", text, sheet.Name, index)
					resultList = append(resultList, warnInfo)
					continue
				}

				cell.Value = value
			}
		}
	}

	// 保存数据
	splitList := strings.Split(this.path, "\\")
	err := this.file.Save(path.Join(tarFilePath, splitList[len(splitList)-1]))
	if err != nil {
		return nil, err
	}

	return resultList, nil
}

// 判断是否包含中文字符串
// 参数：
// word：待判断的单词
// 返回值：
// 1.ture：表示中文，false：非中文
func (this *ExcelUtil) isChinese(word string) bool {
	for _, r := range word {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("^[\u4e00-\u9fa5]{3,8}$").MatchString(string(r))) {
			return true
		}
	}

	return false
}

// 创建新的ExcelUtil文件助手
// 参数：
// path：Excel文件路径
// 返回值：
// 1.助手对象
// 2.错误对象
func NewExcelUtil(path string) (*ExcelUtil, error) {
	file, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("NewExcelUtil报错err=%s,path=%s", err, path))
	}

	return &ExcelUtil{
		file: file,
		path: path,
	}, nil
}
