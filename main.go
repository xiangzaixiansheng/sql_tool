package main

import (
	"fmt"
	"sql_tool/conf"
	"sql_tool/utils"
)

func main() {
	conf.Init()
	utils.InitMysql()
	//sql => csv
	//转义用\"
	rows, _ := utils.Query("SELECT * FROM cusCopy WHERE 1")

	// fileName, _ := utils.GetFilePath("test2.csv")
	// err := utils.WriteFile(fileName, rows)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//sql => xlsx
	fileName2, _ := utils.GetFilePath("All.xlsx")
	err := utils.WriteXlsFile(fileName2, rows)
	if err != nil {
		fmt.Println(err)
	}

	//csv=>insert
	//utils.CsvTosql("test", "./test2.csv", "./test2.sql")
}
