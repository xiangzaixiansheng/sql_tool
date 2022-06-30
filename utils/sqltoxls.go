package utils

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/tealeg/xlsx"
)

//exports 写xlsx文件的内容
func WriteXlsFile(csvFileName string, rows *sql.Rows) error {
	return NewXls(rows).WriteXlsFile(csvFileName)
}

type ConverterXls struct {
	Headers      []string // Column headers to use (default is rows.Columns())
	WriteHeaders bool     // Flag to output headers in your CSV (default is true)
	TimeFormat   string   // Format string for any time.Time values (default is time's default)
	FloatFormat  string   // Format string for any float64 and float32 values (default is %v)
	Delimiter    rune     // Delimiter to use in your CSV (default is comma)

	rows *sql.Rows
}

func (c ConverterXls) WriteXlsFile(xlsFileName string) error {
	// 创建文件
	f, err := os.Create(xlsFileName)
	if err != nil {
		f.Close()
		return err
	}

	err = c.Write(xlsFileName)
	if err != nil {
		return err
	}

	return f.Close()
}

func (c ConverterXls) Write(xlsFileName string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var xlsxRow *xlsx.Row
	var cell *xlsx.Cell
	var err error
	var headers []string

	rows := c.rows

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return err
	}

	//xlsx的标题
	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}

	//如果需要写标题的话
	if c.WriteHeaders {
		// use Headers if set, otherwise default to
		// query Columns
		if len(c.Headers) > 0 {
			headers = c.Headers
		} else {
			headers = columnNames
		}
	}

	//生成标题
	firstRow := headers
	xlsxRow = sheet.AddRow()
	for _, colName := range firstRow {
		cell = xlsxRow.AddCell()
		cell.Value = colName
	}

	//生成json数据
	count := len(columnNames)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	//处理数据
	for rows.Next() {
		//处理每一条数据
		row := make([]string, count)

		for i, _ := range columnNames {
			valuePtrs[i] = &values[i]
		}

		if err = rows.Scan(valuePtrs...); err != nil {
			return err
		}

		for i, _ := range columnNames {
			var value interface{}
			rawValue := values[i]

			byteArray, ok := rawValue.([]byte)
			if ok {
				value = string(byteArray)
			} else {
				value = rawValue
			}

			float64Value, ok := value.(float64)
			if ok && c.FloatFormat != "" {
				value = fmt.Sprintf(c.FloatFormat, float64Value)
			} else {
				float32Value, ok := value.(float32)
				if ok && c.FloatFormat != "" {
					value = fmt.Sprintf(c.FloatFormat, float32Value)
				}
			}

			timeValue, ok := value.(time.Time)
			if ok && c.TimeFormat != "" {
				value = timeValue.Format(c.TimeFormat)
			}

			if value == nil {
				row[i] = ""
			} else {
				row[i] = fmt.Sprintf("%v", value)
			}
		}
		//新增一行
		xlsxRow = sheet.AddRow()
		for i := 0; i < len(headers); i++ {
			val := row[i]
			cell = xlsxRow.AddCell()
			cell.Value = val
		}
	}

	fmt.Println("Saved to ", xlsFileName)
	error := file.Save(xlsFileName)
	return error

}

func NewXls(rows *sql.Rows) *ConverterXls {
	return &ConverterXls{
		rows:         rows,
		WriteHeaders: true,
		Delimiter:    ',',
	}
}
