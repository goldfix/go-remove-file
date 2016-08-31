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

	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

func initLog(logFile string) {

	InfoCmd = log.New(os.Stdout, "", 0)
	ErrorCmd = log.New(os.Stdout, "Error: ", 0)

	if logFile != "" {
		logFile = filepath.Join(logFile, "grm_"+time.Now().Format("20060102")+".log")
		fileLog, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}
		InfoFile = log.New(fileLog, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		ErrorFile = log.New(fileLog, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		InfoFile = log.New(ioutil.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		ErrorFile = log.New(ioutil.Discard, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func initFolder() {
	currentUser, err := user.Current()
	if err != nil {
		errLog(err, debug.Stack())
	}

	//check if exists folder of recycled
	RECYCLED_FOLDER = filepath.Join(currentUser.HomeDir, ".grm")
	if _, err := os.Stat(RECYCLED_FOLDER); os.IsNotExist(err) {
		err = os.Mkdir(RECYCLED_FOLDER, os.ModePerm)
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
	err := os.RemoveAll(RECYCLED_FOLDER)
	if err != nil {
		errLog(err, debug.Stack())
	}
	infoLog(fmt.Sprintf("Empty recycle folder completed."))
	initFolder()
}

//-----------------------------//

func askForConfirmation(fileName string, isFile bool) bool {
	var response string

	if isFile {
		infoLog(fmt.Sprintf("Process file: '" + fileName + "'?"))
	} else {
		infoLog(fmt.Sprintf("Are you sure empty recycle folder?"))
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
		infoLog(fmt.Sprintf("Please type yes or no (or y or n) and then press enter:\n"))
		return askForConfirmation(fileName, isFile)
	}
}

func infoLog(mex string) bool {
	InfoCmd.Printf(mex)
	InfoFile.Printf(mex)
	return true
}

func errLog(err error, stack []byte) bool {

	if err == nil {
		return true
	}

	////TODO: use 'signal' --> https://golang.org/pkg/os/signal/
	if err.Error() == "EOF" {
		ErrorCmd.Println("Pperation has been interrupted")
		ErrorFile.Println("Operation has been interrupted")
	} else {
		ErrorCmd.Println(err.Error(), string(stack))
		ErrorFile.Println(err.Error(), string(stack))
	}
	os.Exit(1)
	return false
}
