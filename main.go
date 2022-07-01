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
	//select 修改超时时间3分钟 一个超级大的表"SELECT /*+ MAX_EXECUTION_TIME(3 * 60 * 1000) */ * FROM df_property_info_0000 WHERE 1"
	rows, _ := utils.Query("SELECT * FROM cusCopy WHERE 1")
	defer rows.Close()

	// fileName, _ := utils.GetFilePath("devices.csv")
	// err := utils.WriteFile(fileName, rows)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//sql => xlsx
	// fileName2, _ := utils.GetFilePath("All.xlsx")
	// err := utils.WriteXlsFile(fileName2, rows)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//sql => json
	fileName3, _ := utils.GetFilePath("All.json")
	err := utils.WriteJsonFile(fileName3, rows)
	if err != nil {
		fmt.Println(err)
	}

	//csv=>insert
	//utils.CsvTosql("test", "./test2.csv", "./test2.sql")

	if err := rows.Err(); err != nil {
		fmt.Println(err, "Rows error.")
	}

	if err := rows.Close(); err != nil {
		fmt.Println(err, "can't make `rows.Close().")
	}
}
