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

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime/debug"
	"strings"
)

func initLog(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	fileLog, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}
	traceHandle = io.MultiWriter(traceHandle, fileLog)
	infoHandle = io.MultiWriter(fileLog, infoHandle)
	warningHandle = io.MultiWriter(fileLog, warningHandle)
	errorHandle = io.MultiWriter(fileLog, errorHandle)

	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func initFolder() {
	currentUser, err := user.Current()
	if err != nil {
		errLog(err, debug.Stack())
	}

	//check if exists folder of recycled
	RECYCLED_FOLDER = filepath.Join(currentUser.HomeDir, ".grm")
	if _, err := os.Stat(RECYCLED_FOLDER); os.IsNotExist(err) {
		err = os.Mkdir(RECYCLED_FOLDER, os.ModeDir)
		if err != nil {
			errLog(err, debug.Stack())
		}
		f, err := os.OpenFile(filepath.Join(RECYCLED_FOLDER, "_DO_NOT_REMOVE_"), os.O_CREATE, os.ModePerm)
		if err != nil {
			errLog(err, debug.Stack())
		}
		defer f.Close()
	}

	//check if exists file of list deleted files and init
	RECYCLED_FILEDB = filepath.Join(RECYCLED_FOLDER, ".grm.db")
	saveInfoDeletedFile(nil, true)
}

func emptyRecycle() {
	if !askForConfirmation("", false) {
		return
	}
	err := os.RemoveAll(RECYCLED_FOLDER)
	if err != nil {
		errLog(err, debug.Stack())
	}
	Info.Printf("Empty recycle folder completed.")
	initFolder()
}

//-----------------------------//

func askForConfirmation(fileName string, isFile bool) bool {
	var response string

	if isFile {
		Info.Printf("process file: '" + fileName + "'?")
	} else {
		Info.Printf("are you sure empty recycle folder?")
	}

	_, err := fmt.Scan(&response)
	if err != nil {
		fmt.Print(response)
		errLog(err, debug.Stack())
	}

	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		return true
	} else if strings.ToLower(response) == "n" || strings.ToLower(response) == "no" {
		return false
	} else {
		Info.Printf("Please type yes or no (or y or n) and then press enter:\n")
		return askForConfirmation(fileName, isFile)
	}
}

func errLog(err error, stack []byte) bool {

	if err == nil {
		return true
	}

	////TODO: use 'signal' --> https://golang.org/pkg/os/signal/
	if err.Error() == "EOF" {
		Error.Println("operation has been interrupted")
	} else {
		Error.Println(err.Error(), string(stack))
	}
	os.Exit(1)
	return false
}
