一个自用的sql工具。

想法💡：平时总是需要倒出一些sql的查询结果，但是公司的提供平台还不怎么好使，一些软件还不能用。😅



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

四、查询结果倒出xlsx

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



TODO:

1⃣️查询数据生成json文件

2⃣️~~查询数据生成excel文件~~

3⃣️~~优化生成文件路径~~

4⃣️sql只支持查询，屏蔽其他操作语句。(安全考虑)

5⃣️增加Gui界面操作

6⃣️优化：根绝需要生成的文件名，自动判断生成的逻辑
