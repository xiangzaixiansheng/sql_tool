package utils

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"io"
	"os"
)

//exports 写json文件的内容
func WriteJsonFile(jsonFileName string, rows *sql.Rows) error {
	return NewJson(rows).WriteJsonFile(jsonFileName)
}

type ConverterJson struct {
	Headers      []string // Column headers to use (default is rows.Columns())
	WriteHeaders bool     // Flag to output headers in your CSV (default is true)
	TimeFormat   string   // Format string for any time.Time values (default is time's default)
	FloatFormat  string   // Format string for any float64 and float32 values (default is %v)
	Delimiter    rune     // Delimiter to use in your CSV (default is comma)

	rows *sql.Rows
}

func (c ConverterJson) WriteJsonFile(jsonFileName string) error {
	// 创建文件
	f, err := os.Create(jsonFileName)
	if err != nil {
		f.Close()
		return err
	}

	err = c.Write(jsonFileName, f)
	if err != nil {
		return err
	}

	return f.Close()
}

func (c ConverterJson) Write(jsonFileName string, writer io.Writer) error {

	var err error
	results := make([]map[string]interface{}, 0)

	columns, _ := c.rows.Columns()
	data := make([][]byte, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range data {
		pointers[i] = &data[i]
	}
	for c.rows.Next() {
		//组装一个json
		row := make(map[string]interface{}, 0)
		err = c.rows.Scan(pointers...)
		if err != nil {
			return err
		}
		for key := range data {
			row[columns[key]] = string(data[key])
		}
		results = append(results, row)
	}
	bytes2, _ := json.Marshal(results)
	stringData2 := string(bytes2)
	//写文件喽
	w := bufio.NewWriter(writer) //创建新的 Writer 对象
	_, err = w.WriteString(stringData2)
	w.Flush()

	return err
}

func NewJson(rows *sql.Rows) *ConverterJson {
	return &ConverterJson{
		rows:         rows,
		WriteHeaders: true,
		Delimiter:    ',',
	}
}
