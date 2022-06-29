package main

import (
	"sql_tool/conf"
	"sql_tool/utils"
)

func main() {
	conf.Init()
	utils.InitMysql()
	//sql => csv
	//转义用\"
	rows, _ := utils.Query("SELECT * FROM order WHERE id < 200")

	err := utils.WriteFile("./test2.csv", rows)
	if err != nil {
		panic(err)
	}

	//csv=>insert
	//utils.CsvTosql("test", "./test2.csv", "./test2.sql")
}
