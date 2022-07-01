package utils

import (
	"fmt"
	"os"
	conf "sql_tool/conf"
	"time"

	//切记：导入驱动包
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitMysql() {

	driverName := "mysql"
	//数据库连接
	user := conf.DbUser
	pwd := conf.DbPassWord
	host := conf.DbHost
	port := conf.DbPort
	dbname := conf.DbName

	//dbConn := "root:yu271400@tcp(127.0.0.1:3306)/myblog?charset=utf8"
	dbConn := user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8"

	db1, err := sql.Open(driverName, dbConn)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		db = db1
		db.SetConnMaxLifetime(time.Hour * 4)
	}
}

//操作数据库
func ModifyDB(sql string, args ...interface{}) (int64, error) {
	result, err := db.Exec(sql, args...)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return count, nil
}

/**
  database/sql提供了Query和QueryRow方法进行查询数据库
  Query 返回游标，需要迭代Next
  QueryRow 读取单条,
*/
//查询
func QueryRowDB(sql string) *sql.Row {
	return db.QueryRow(sql)
}

//查询多条数据
func Query(sql string) (*sql.Rows, error) {
	return db.Query(sql)
}
