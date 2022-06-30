package utils

//csv 格式数据转insert sql语句

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

var csvSep = []byte{','}

func CsvTosql(tablename, inputpath, outputpath string) {
	var err error
	var inFile *os.File

	inFile, err = os.Open(inputpath)
	if err != nil {
		log.Printf("cannot open file %s: %v", inputpath, err)
		os.Exit(1)
	}
	defer inFile.Close()

	var outFile *os.File
	outFile, err = os.Create(outputpath)
	if err != nil {
		log.Printf("cannot open file %s: %v", outputpath, err)
		os.Exit(1)
	}
	defer outFile.Close()

	br := bufio.NewReader(inFile)
	for i := 0; true; i++ {
		line, err := br.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("cannot read file %s: %v", inputpath, err)
			os.Exit(1)
		}
		line = bytes.Trim(line, "\n\r ")
		//标题
		if i == 0 {
			fields := bytes.Split(line, csvSep)
			err = writeInsertIntoFields(outFile, tablename, fields)
			if err != nil {
				fmt.Printf("cannot write file %s: %v", outputpath, err)
				os.Exit(1)
			}
			continue
		}
		//处理内容的
		if i != 1 {
			_, err := outFile.Write([]byte{',', '\n'})
			if err != nil {
				fmt.Printf("cannot write file %s: %v", outputpath, err)
				os.Exit(1)
			}
		}
		values := bytes.Split(line, csvSep)
		err = writeValues(outFile, values)
		if err != nil {
			fmt.Printf("cannot write file %s: %v", outputpath, err)
			os.Exit(1)
		}
	}
	_, err = outFile.Write([]byte{';', '\n'})
	if err != nil {
		fmt.Printf("cannot write file %s: %v", outputpath, err)
		os.Exit(1)
	}
}

//处理开头的
func writeInsertIntoFields(w io.Writer, tableName string, fields [][]byte) error {
	_, err := io.WriteString(w, fmt.Sprintf("INSERT INTO `%s` (", tableName))
	if err != nil {
		return err
	}
	for i, f := range fields {
		if i != 0 {
			_, err = w.Write([]byte{',', ' '})
			if err != nil {
				return err
			}
		}
		_, err = w.Write([]byte{'`'})
		if err != nil {
			return err
		}

		_, err = w.Write(trimField(f))
		if err != nil {
			return err
		}

		_, err = w.Write([]byte{'`'})
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(w, ") \nVALUES\n")
	if err != nil {
		return err
	}
	return nil
}

func writeValues(w io.Writer, lineValues [][]byte) error {
	_, err := w.Write([]byte{'('})
	if err != nil {
		return err
	}
	for i, v := range lineValues {
		if i != 0 {
			_, err := w.Write([]byte{',', ' '})
			if err != nil {
				return err
			}
		}
		_, err := w.Write([]byte{'\''})
		if err != nil {
			return err
		}

		_, err = w.Write(trimValue(v))
		if err != nil {
			return err
		}

		_, err = w.Write([]byte{'\''})
		if err != nil {
			return err
		}
	}
	_, err = w.Write([]byte{')'})
	if err != nil {
		return err
	}
	return nil
}

func trimField(f []byte) []byte {
	return bytes.Trim(f, " 	\"\uFEFF")
}

func trimValue(f []byte) []byte {
	return bytes.Trim(f, " 	\uFEFF")
}
