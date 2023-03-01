package utils

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//csv 导入db的脚本

const (
	DELIMITER           = ',' // default delimiter for csv files
	MAX_SQL_CONNECTIONS = 10  // default max_connections of mysql is 150,
	Db                  = "mysql"
	DbHost              = "127.0.0.1"
	DbPort              = "3306"
	DbUser              = "root"
	DbPassWord          = "xiangzai"
	DbName              = "sqlstudy"
)

func CSVtoDb(tableName, filePath string) {

	dbConn := DbUser + ":" + DbPassWord + "@tcp(" + DbHost + ":" + DbPort + ")/" + DbName + "?charset=utf8"
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	db.SetMaxIdleConns(MAX_SQL_CONNECTIONS)
	defer db.Close()
	fmt.Println("=====db=====")
	// --------------------------------------------------------------------------
	// 加载文件并读取
	// --------------------------------------------------------------------------

	//打开文件(只读模式)，创建io.read接口实例
	opencast, err := os.Open(filePath)
	if err != nil {
		log.Printf("%v文件打开失败！", filePath)
	}
	defer opencast.Close()

	start := time.Now()
	query := ""
	callback := make(chan int) // callback channel for insert goroutines
	connections := 0           // number of concurrent connections
	insertions := 0            // counts how many insertions have finished

	available := make(chan bool, MAX_SQL_CONNECTIONS) // buffered channel, holds number of available connections
	for i := 0; i < MAX_SQL_CONNECTIONS; i++ {
		available <- true
	}

	// start status logger
	startLogger(&insertions, &connections)

	// start connection controller
	startConnectionController(&insertions, &connections, callback, available)

	// --------------------------------------------------------------------------
	// read rows and insert into database
	// --------------------------------------------------------------------------

	//创建csv读取接口实例
	ReadCsv := csv.NewReader(opencast)

	var wg sync.WaitGroup
	id := 1
	isFirstRow := true

	for {
		//获取一行内容，一般为第一行内容
		record, err := ReadCsv.Read() //返回切片类型：[chen  hai wei]
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}
		if isFirstRow {
			// 解析第一行的表头
			parseColumns(tableName, record, &query)
			isFirstRow = false
		} else if <-available { // wait for available database connection
			connections += 1
			id += 1
			wg.Add(1)
			go insert(id, query, db, callback, &connections, &wg, strToInterface(record))
		}
		wg.Wait()

		endTime := time.Since(start)
		log.Printf("Status: %d insertions\n", insertions)
		log.Printf("Execution time: %s\n", endTime)

	}

}

func parseColumns(tableName string, columns []string, query *string) {
	*query = "INSERT INTO " + tableName + " ("
	values := "VALUES ("
	for i, c := range columns {
		if i == 0 {
			*query += c
			values += "?"
		} else {
			*query += ", " + c
			values += ", ?"
		}
	}
	values += ")"
	*query += ") " + values
	fmt.Println("query", *query)
}

// inserts data into database
func insert(id int, query string, db *sql.DB, callback chan<- int, conns *int, wg *sync.WaitGroup, args []interface{}) {

	// make a new statement for every insert,
	// this is quite inefficient, but since all inserts are running concurrently,
	// it's still faster than using a single prepared statement and
	// inserting the data sequentielly.
	// we have to close the statement after the routine terminates,
	// so that the connection to the database is released and can be reused
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Printf("ID: %d (%d conns), %s\n", id, *conns, err.Error())
	}

	// finished inserting, send id over channel to signalize termination of routine
	callback <- id
	wg.Done()
}

// controls termination of program and number of connections to database
func startConnectionController(insertions, connections *int, callback <-chan int, available chan<- bool) {

	go func() {
		for {

			<-callback // returns id of terminated routine

			*insertions += 1  // a routine terminated, increment counter
			*connections -= 1 // and unregister its connection

			available <- true // make new connection available
		}
	}()
}

func startLogger(insertions, connections *int) {

	go func() {
		c := time.Tick(time.Second)
		for {
			<-c
			log.Printf("Status: %d insertions, %d database connections\n", *insertions, *connections)
		}
	}()
}

// convert []string to []interface{}
func strToInterface(s []string) []interface{} {
	i := make([]interface{}, len(s))
	for k, v := range s {
		i[k] = v
	}
	return i
}
