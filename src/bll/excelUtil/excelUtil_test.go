package excelUtil

import (
	"fmt"
	"testing"

	"moqikaka.com/goutil/fileUtil"
)

func TestReadExcel(t *testing.T) {
	excelUtil, err := NewExcelUtil("E:\\Japanese_Translate.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	resultDict := excelUtil.ReadExcelToDict()
	//for key, value := range resultDict {
	//	fmt.Println(fmt.Sprintf("%s--%s", key, value))
	//}

	filePathList, err := fileUtil.GetFileList("E:\\test\\")
	if err != nil {
		fmt.Println(err)
	} else {
		for _, filePath := range filePathList {
			fmt.Println(filePath)
			tempExcelUtil, err := NewExcelUtil(filePath)
			if err != nil {
				fmt.Println(err)
			} else {
				tempExcelUtil.Replace(resultDict, "E:\\pp")
			}
		}
	}
}
