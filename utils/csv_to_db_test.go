package utils

import (
	"fmt"
	"testing"
)

func TestFileCache(t *testing.T) {
	fmt.Println("开始导入CSV")
	CSVtoDb("user", "/Users/hanxiang1/work/gogo/sql_tool/output/user.csv")

}
