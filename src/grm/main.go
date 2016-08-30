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
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	guuid "github.com/satori/go.uuid"
)

type InfoDeletedFile struct {
	name      string
	date      time.Time
	uuid      guuid.UUID
	pathFile  string
	toProcess bool
}

const VERSION string = "0.1"
const SEPARATOR rune = '|'

var (
	Trace           *log.Logger
	Info            *log.Logger
	Warning         *log.Logger
	Error           *log.Logger
	RECYCLED_FILEDB string
	RECYCLED_FOLDER string
)

func main() {
	initLog(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	initFolder()

	paramListRecycleFiles := flag.Bool("ls", false, "list files into recycle")
	paramRecoverFiles := flag.Bool("u", false, "recover list files from recycle")
	paramDeleteFiles := flag.Bool("d", false, "path to process")
	paramEmptyRecycle := flag.Bool("e", false, "empty recycle folder")
	paramForceDelete := flag.Bool("x", false, "force delete files")

	paramForce := flag.Bool("f", false, "ignore nonexistent files, never prompt")
	paramPrompt := flag.Bool("i", false, "prompt before every removal")
	paramRecursive := flag.Bool("r", false, "remove contents recursively")
	paramVerbose := flag.Bool("v", false, "explain what is being done")
	paramVersion := flag.Bool("version", false, "Version")
	flag.Parse()

	if *paramVersion {
		Info.Printf("Version: %s", VERSION)
		os.Exit(0)
	}

	if *paramEmptyRecycle {
		emptyRecycle()
		os.Exit(0)
	}

	if *paramDeleteFiles == false && *paramListRecycleFiles == false && *paramRecoverFiles == false {
		err := errors.New("Please use parameters: '-d', '-ls', '-u'")
		errLog(err, nil)
	}

	pathToProcess := ""
	if (len(flag.Args()) != 1 || flag.Arg(0) == "") && (*paramDeleteFiles || *paramListRecycleFiles || *paramRecoverFiles) {
		err := errors.New("missing or wrong operand")
		errLog(err, nil)
	} else {
		pathToProcess = flag.Arg(0)
	}

	if *paramVerbose {
		Info.Printf("Path to Process: %s", pathToProcess)
		Info.Printf("Param List Files: %t", *paramListRecycleFiles)
		Info.Printf("Param Recovery Files: %t", *paramRecoverFiles)
		Info.Printf("Param Path to Delete: %t", *paramDeleteFiles)
		Info.Printf("Empty Recycle: %t", *paramEmptyRecycle)

		Info.Printf("Param Prompt: %t", *paramPrompt)
		Info.Printf("Param Verbose: %t", *paramVerbose)
		Info.Printf("Param Force: %t", *paramForce)
		Info.Printf("Param Recursive: %t", *paramRecursive)
		Info.Printf("Force Delete Files: %t", *paramForceDelete)

		Info.Printf("Param Version: %t", *paramVersion)
		Info.Printf("Recycle Folder: %s", RECYCLED_FOLDER)
		Info.Printf("Recycle FileDb: %s", RECYCLED_FILEDB)
	}

	if *paramListRecycleFiles {
		Info.Println("List files: " + pathToProcess)
		getListDeletedFiles(pathToProcess, false)
		os.Exit(0)
	}

	if *paramRecoverFiles {
		Info.Println("Recover files: " + pathToProcess)
		recoverFiles(pathToProcess, *paramPrompt)
		os.Exit(0)
	}

	if *paramDeleteFiles {
		Info.Println("Delete files: " + pathToProcess)
		var filesToDelete []string
		filesToDelete = deleteFiles(pathToProcess, *paramRecursive, filesToDelete, *paramPrompt, *paramForceDelete)
		os.Exit(0)
	}
}

func deleteFiles(pathToProcess string, recursive bool, listFiles []string, paramPrompt bool, paramForceDelete bool) (filesToDelete []string) {

	pathToProcess = filepath.Clean(pathToProcess)
	pathDir := filepath.Dir(pathToProcess)
	baseName := filepath.Base(pathToProcess)

	Info.Println("pathToProcess: " + pathToProcess)
	Info.Println("pathDir: " + pathDir)
	Info.Println("baseName: " + baseName)

	listFiles = getFilesFromFolder(pathDir, baseName, recursive, listFiles)

	infoDeletedFile := loadInfoDeletedFile()
	infoDeletedFile = appendInfoDeletedFile(listFiles, infoDeletedFile, paramPrompt)
	moveFilesToRecycle(infoDeletedFile, paramForceDelete)

	if !paramForceDelete {
		saveInfoDeletedFile(infoDeletedFile, false)
	}

	return listFiles
}

func recoverFiles(paramRecoverFiles string, paramPrompt bool) {
	infoDeletedFile := getListDeletedFiles(paramRecoverFiles, paramPrompt)
	moveFilesFromRecycle(infoDeletedFile)
	var newInfoDeletedFile []InfoDeletedFile
	for _, v := range infoDeletedFile {
		if v.toProcess {
			continue
		} else {
			newInfoDeletedFile = append(newInfoDeletedFile, v)
		}
	}
	saveInfoDeletedFile(newInfoDeletedFile, false)
}

func getFilesFromFolder(folderPath string, matching string, recursive bool, listFiles []string) (resultListFiles []string) {

	tmpListFiles, err := ioutil.ReadDir(folderPath)
	if err != nil {
		errLog(err, debug.Stack())
	}
	for _, k := range tmpListFiles {
		tmpFilePath := filepath.Join(folderPath, k.Name())
		chkMatching, err := filepath.Match(matching, k.Name())
		if err != nil {
			errLog(err, debug.Stack())
		}

		if chkMatching {
			if k.IsDir() {
				if recursive {
					Info.Println("Folder: " + tmpFilePath)
					listFiles = getFilesFromFolder(tmpFilePath, "*", recursive, listFiles)
				}
			} else {
				Info.Println("File: " + tmpFilePath)
				listFiles = append(listFiles, tmpFilePath)
			}
		}
	}

	return listFiles
}

func appendInfoDeletedFile(filesToDelete []string, infoDeletedFile []InfoDeletedFile, paramPrompt bool) (resultInfoDeletedFile []InfoDeletedFile) {

	for _, v := range filesToDelete {
		if !paramPrompt || askForConfirmation(v, true) {
			//DELETE
			var info InfoDeletedFile
			info.date = time.Now()
			info.uuid = guuid.NewV4()
			info.name = filepath.Base(v)
			info.pathFile = filepath.Dir(v)
			info.toProcess = true
			infoDeletedFile = append(infoDeletedFile, info)
		} else {
			//NO_DELETE
			continue
		}
	}
	return infoDeletedFile
}

func moveFilesFromRecycle(infoDeletedFile []InfoDeletedFile) {
	for _, v := range infoDeletedFile {
		if v.toProcess {
			//RESTORE
			Info.Printf("Recovered file: %s", filepath.Join(v.pathFile, v.name))
			//			err := os.Rename(filepath.Join(RECYCLED_FOLDER, v.uuid.String()+"_"+v.name), filepath.Join(v.pathFile, v.name))
			//			if err != nil {
			//				errLog(err, debug.Stack())
			//			}
		}
	}
}

func moveFilesToRecycle(infoDeletedFile []InfoDeletedFile, paramForceDelete bool) {
	for _, v := range infoDeletedFile {
		if v.toProcess {
			Info.Printf("Deleted file: %s", filepath.Join(v.pathFile, v.name))
			if paramForceDelete {
				err := os.Remove(filepath.Join(v.pathFile, v.name))
				if err != nil {
					errLog(err, debug.Stack())
				}
			} else {
				//			err := os.Rename(filepath.Join(v.pathFile, v.name), filepath.Join(RECYCLED_FOLDER, v.uuid.String()+"_"+v.name))
				//			if err != nil {
				//				errLog(err, debug.Stack())
				//			}
			}

		}
	}
}

func getListDeletedFiles(filter string, paramPrompt bool) (resultInfoDeletedFile []InfoDeletedFile) {
	infoDeletedFile := loadInfoDeletedFile()

	for k, v := range infoDeletedFile {
		if b, err := filepath.Match(filter, v.name); b && errLog(err, debug.Stack()) {
			//if strings.Contains(v.name, filter) {
			Info.Printf("%d: %s  %s  %s", k, v.uuid.String()[:8], v.pathFile, v.name)
			if !paramPrompt || askForConfirmation(v.name, true) {
				infoDeletedFile[k].toProcess = true
			} else {
				infoDeletedFile[k].toProcess = false
			}
			continue
		}
		if b, err := filepath.Match(filter, v.pathFile); b && errLog(err, debug.Stack()) {
			//if strings.Contains(v.pathFile, filter) {
			Info.Printf("%d: %s  %s  %s", k, v.uuid.String()[:8], v.pathFile, v.name)
			if !paramPrompt || askForConfirmation(v.name, true) {
				infoDeletedFile[k].toProcess = true
			} else {
				infoDeletedFile[k].toProcess = false
			}
			continue
		}
		//if b, err := filepath.Match(filter, v.uuid.String()); b && errLog(err, debug.Stack()) {
		if strings.Contains(v.uuid.String(), filter) {
			Info.Printf("%d: %s  %s  %s", k, v.uuid.String()[:8], v.pathFile, v.name)
			if !paramPrompt || askForConfirmation(v.name, true) {
				infoDeletedFile[k].toProcess = true
			} else {
				infoDeletedFile[k].toProcess = false
			}
			continue
		}
	}

	return infoDeletedFile
}
