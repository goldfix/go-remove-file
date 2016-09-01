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

// grm project main.go
package main

import (
	"os"
	"runtime/debug"
	"strconv"
	"testing"
	"time"

	guuid "github.com/satori/go.uuid"
)

func Test_saveInfoDeletedFile(t *testing.T) {
	initLog("")
	initFolder()

	var infoDeletedFile []InfoDeletedFile

	for i := 0; i < 11; i++ {
		var info InfoDeletedFile
		info.date = time.Now()
		info.name = "Alfa_Test_" + strconv.Itoa(i)
		info.pathFile = "Path_" + strconv.Itoa(i)
		info.toProcess = false
		info.uuid = guuid.NewV4()
		infoDeletedFile = append(infoDeletedFile, info)
	}

	saveInfoDeletedFile(infoDeletedFile, false)

	f, err := os.OpenFile(RECYCLED_FILEDB, os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Errorf("Test_saveInfoDeletedFile: %s - %s", err.Error(), string(debug.Stack()))
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		t.Errorf("Test_saveInfoDeletedFile: %s - %s", err.Error(), string(debug.Stack()))
	}

	if stat.Size() != 866 {
		t.Errorf("Test_saveInfoDeletedFile: stat.Size() != 866")
	}
}

func Test_loadInfoDeletedFile(t *testing.T) {
	initFolder()

	f, err := os.OpenFile(RECYCLED_FILEDB, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Errorf("TestloadInfoDeletedFile: %s - %s", err.Error(), string(debug.Stack()))
	}
	defer f.Close()
	var s string
	s = `20160830_162703|//* Data File 'grm' tool. Please not remove. *//|0.1|nd\n
20160830_162703|file_1 - Copy (2).txt|c:\tmp\folder_test|9980aa96-ac59-4342-8b36-e815c714c4a3
20160830_162703|file_1 - Copy (3).txt|c:\tmp\folder_test|1ed23664-39aa-44bb-bda7-85fec56242bd
20160830_162703|file_1 - Copy (4).txt|c:\tmp\folder_test|95f2c6ea-4cba-4bc0-bc8c-ad7017f888ea
20160830_162703|file_1 - Copy (5).txt|c:\tmp\folder_test|10c73f51-2088-470e-8851-4bc2a141ae2c
20160830_162703|file_1 - Copy (6).txt|c:\tmp\folder_test|227d903b-de47-405d-81b0-32477b4b057d
20160830_162703|file_1 - Copy (7).txt|c:\tmp\folder_test|4a4c3323-ae2d-4c09-8d24-a6ec497db92d
20160830_162703|file_1 - Copy (8).txt|c:\tmp\folder_test|4b7724ed-1cbf-4df9-befb-9c327803e302
20160830_162703|file_1 - Copy (9).txt|c:\tmp\folder_test|0902a2eb-ac0a-46d4-b922-15673b9e3857
20160830_162703|file_1 - Copy.txt|c:\tmp\folder_test|f1908a6a-a0a7-4898-b9c3-f594e9a3bff6
20160830_162703|file_1.txt|c:\tmp\folder_test|eeea63e4-21e4-4f5a-a13b-b0038a4376bb
`
	_, err = f.WriteString(s)
	if err != nil {
		t.Errorf("TestloadInfoDeletedFile: %s - %s", err.Error(), string(debug.Stack()))
	}

	resultInfoDeletedFile := loadInfoDeletedFile()
	if len(resultInfoDeletedFile) != 10 {
		t.Errorf("TestloadInfoDeletedFile: len(resultInfoDeletedFile) != 10")
	}
}
