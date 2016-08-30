/**
*
* The MIT License (MIT)
* Copyright (c) 2016 pietro partescano
*
* Permission is hereby granted, free of charge, to any person obtaining a copy of
* this software and associated documentation files (the "Software"), to deal in
* the Software without restriction, including without limitation the rights to
* use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
* of the Software, and to permit persons to whom the Software is furnished to do
* so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
* WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
* CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*
**/

// data_serialize
package main

import (
	"encoding/csv"
	"os"
	"runtime/debug"
	"strings"
	"time"

	guuid "github.com/satori/go.uuid"
)

func saveInfoDeletedFile(infoDeletedFile []InfoDeletedFile, initFile bool) {

	if initFile {
		if _, err := os.Stat(RECYCLED_FILEDB); os.IsNotExist(err) {
			//continue
		} else {
			return
		}
	}

	f, err := os.OpenFile(RECYCLED_FILEDB, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		errLog(err, debug.Stack())
	}
	_, err = f.WriteString(strings.Join([]string{time.Now().Format("20060102_150405"), "//* Data File 'grm' tool. Please not remove. *//", VERSION, "nd\n"}, string(SEPARATOR)))
	defer f.Close()

	csvFile := csv.NewWriter(f)
	csvFile.Comma = rune(SEPARATOR)
	defer csvFile.Flush()

	var records [][]string
	for _, v := range infoDeletedFile {
		records = append(records, []string{v.date.Format("20060102_150405"), v.name, v.pathFile, v.uuid.String()})
	}
	err = csvFile.WriteAll(records)
	if err != nil {
		errLog(err, debug.Stack())
	}
}

func loadInfoDeletedFile() (resultInfoDeletedFile []InfoDeletedFile) {
	f, err := os.OpenFile(RECYCLED_FILEDB, os.O_RDONLY, os.ModePerm)
	if err != nil {
		errLog(err, debug.Stack())
	}
	defer f.Close()
	csvFile := csv.NewReader(f)
	csvFile.Comma = SEPARATOR
	csvFile.TrimLeadingSpace = true
	var result [][]string
	result, err = csvFile.ReadAll()
	if err != nil {
		errLog(err, debug.Stack())
	}

	for k, v := range result {
		if k == 0 {
			continue
		}
		var info InfoDeletedFile
		info.date, err = time.Parse("20060102_150405", v[0])
		if err != nil {
			errLog(err, debug.Stack())
		}
		info.name = v[1]
		info.pathFile = v[2]
		info.uuid, err = guuid.FromString(v[3])
		if err != nil {
			errLog(err, debug.Stack())
		}
		resultInfoDeletedFile = append(resultInfoDeletedFile, info)
	}

	return resultInfoDeletedFile
}
