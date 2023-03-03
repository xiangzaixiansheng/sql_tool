一个自用的sql工具。

想法💡：平时总是需要倒出一些sql的查询结果，但是公司的提供平台还不怎么好使，一些软件还不能用。😅



[TOC]

简单介绍一下：

##### 一、配置文件：

conf/app.ini

这里面配置mysql相关的内容

```
[mysql]
Db = mysql
DbHost = 127.0.0.1
DbPort = 3306
DbUser = root
DbPassWord = xiangzai
DbName = order
```

##### 二、查询结果倒出csv

```
	//转义用\"
	rows, _ := utils.Query("SELECT * FROM order WHERE id < 200")

	err := utils.WriteFile("./test2.csv", rows)
	if err != nil {
		fmt.Println(err)
	}
```

##### 三、csv生成插入insert语句

```
utils.CsvTosql("test", "./test2.csv", "./test2.sql")
```

##### 四、查询结果倒出xlsx

```
	//sql => xlsx
	fileName2, _ := utils.GetFilePath("All.xlsx")
	err := utils.WriteXlsFile(fileName2, rows)
	if err != nil {
		fmt.Println(err)
	}
```

使用包:

```
xlsx: github.com/tealeg/xlsx
mysql: github.com/go-sql-driver/mysql
ini: gopkg.in/ini.v1
```

##### 五、csv导入到数据库脚本

执行脚本：

utils/csv_to_db_test.go

```go
CSVtoDb("user", "/Users/hanxiang1/work/gogo/sql_tool/output/user.csv")
```

配置参数

utils/csv_to_db.go

```go
const (
	DELIMITER           = ',' // default delimiter for csv files
	MAX_SQL_CONNECTIONS = 10  // default max_connections of mysql is 150,
	Db                  = "mysql"
	DbHost              = "127.0.0.1"
	DbPort              = "3306"
	DbUser              = "root"
	DbPassWord          = "xiangzai"
	DbName              = "sqlstudy"
	Step                = 10 //每次插入数据
)
```

本机测试的话 100w数据 66s导入完成



**设置testing的超时时间**
vscode的左下角设置-用户设置-搜索下面的关键字
go test timeout

修改30s->120s



### 错误记录：

1、gorm的 busy buffer问题

```
[mysql] 2022/07/01 12:02:04 packets.go:428: busy buffer
invalid connection
invalid connection Rows error.
```

例如我们查询的数据量比较大的话，可能会耗时比较长。

我们可以先看一下mysql数据库的超时时间：

```
SHOW VARIABLES LIKE '%timeout%'
```



| connect_timeout             | 10       |
| --------------------------- | -------- |
| delayed_insert_timeout      | 300      |
| have_statement_timeout      | YES      |
| innodb_flush_log_at_timeout | 1        |
| innodb_lock_wait_timeout    | 20       |
| innodb_rollback_on_timeout  | OFF      |
| interactive_timeout         | 300      |
| lock_wait_timeout           | 31536000 |
| net_read_timeout            | 30       |
| net_write_timeout           | 60       |
| rpl_stop_slave_timeout      | 31536000 |
| slave_net_timeout           | 60       |
| wait_timeout                | 300      |

**connect_time**

connect_timeout指的是连接过程中握手的超时时间，即MySQL客户端在尝试与MySQL服务器建立连接时，MySQL服务器返回错误握手协议前等待客户端数据包的最大时限。默认10秒。

**interactive_timeout / wait_timeout**

MySQL关闭交互/非交互连接前等待的最大时限。默认28800秒。

**lock_wait_timeout**

sql语句请求元数据锁的最长等待时间，默认为一年。此锁超时对于隐式访问Mysql库中系统表的sql语句无效，但是对于使用select，update语句直接访问MySQL库中标的sql语句有效。

**net_read_timeout / net_write_timeout**

mysql服务器端等待从客户端读取数据 / 向客户端写入数据的最大时限，默认30秒。

**slave_net_timeout**

mysql从复制连结等待读取数据的最大时限，默认3600秒。

解决办法：

MySQL doesn't provide safe and efficient canceling mechanism. When context is cancelled or reached `readTimeout`, `DB.ExecContext` returns without terminating using connection. It cause "invalid connection" next time the connection is used.

https://dev.mysql.com/doc/refman/5.7/en/optimizer-hints.html#optimizer-hints-execution-time

```sql
MAX_EXECUTION_TIME(N)
```

Example with a timeout of 1 second (1000 milliseconds):

```sql
SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM t1 INNER JOIN t2 WHERE ...
```




### TODO:

1⃣️查询数据生成json文件

2⃣️~~查询数据生成excel文件~~

3⃣️~~优化生成文件路径~~

4⃣️sql只支持查询，屏蔽其他操作语句。(安全考虑)

5⃣️增加Gui界面操作

6⃣️优化：根绝需要生成的文件名，自动判断生成的逻辑
