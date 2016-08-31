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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
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

const VERSION string = "0.4"
const SEPARATOR rune = '|'

var (
	InfoCmd         *log.Logger
	InfoFile        *log.Logger
	ErrorCmd        *log.Logger
	ErrorFile       *log.Logger
	RECYCLED_FILEDB string
	RECYCLED_FOLDER string
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		errLog(err, debug.Stack())
	}
	pathLogFile := currentUser.HomeDir
	initLog(pathLogFile)
	initFolder()

	paramListRecycleFiles := flag.Bool("ls", false, "list files into recycle")
	paramRecoverFiles := flag.Bool("u", false, "recover list files from recycle")
	paramDeleteFiles := flag.Bool("d", false, "move files into recycle")
	paramEmptyRecycle := flag.Bool("e", false, "empty recycle folder")
	paramForceDelete := flag.Bool("x", false, "force delete files")

	paramForce := flag.Bool("f", false, "ignore nonexistent files, never prompt")
	paramPrompt := flag.Bool("i", false, "prompt before every files")
	paramRecursive := flag.Bool("r", false, "remove contents recursively")
	paramVerbose := flag.Bool("v", false, "explain what is being done")
	paramVersion := flag.Bool("version", false, "output version information and exit")
	paramPathToRestore := flag.String("t", "", "destination folder to restore")
	flag.Parse()

	if *paramVersion {
		infoLog(fmt.Sprintf("Version: %s", VERSION))
		os.Exit(0)
	}

	if *paramEmptyRecycle {
		if !askForConfirmation("", false) {
			return
		}
		emptyRecycle()
		os.Exit(0)
	}

	if *paramDeleteFiles == false && *paramListRecycleFiles == false && *paramRecoverFiles == false {
		err := errors.New("Please use parameters: '-d', '-ls', '-u'")
		errLog(err, nil)
	}

	pathToProcess := ""
	if (len(flag.Args()) != 1 || flag.Arg(0) == "") && (*paramDeleteFiles || *paramListRecycleFiles || *paramRecoverFiles) {
		err := errors.New("Missing or wrong operand")
		errLog(err, nil)
	} else {
		pathToProcess = flag.Arg(0)
	}

	if *paramVerbose {
		infoLog(fmt.Sprintf("Path to Process: %s", pathToProcess))
		infoLog(fmt.Sprintf("Flag List Files: %t", *paramListRecycleFiles))
		infoLog(fmt.Sprintf("Flag Recovery Files: %t", *paramRecoverFiles))
		infoLog(fmt.Sprintf("Flag Path to Delete: %t", *paramDeleteFiles))
		infoLog(fmt.Sprintf("Flag Empty Recycle: %t", *paramEmptyRecycle))

		infoLog(fmt.Sprintf("Flag Prompt: %t", *paramPrompt))
		infoLog(fmt.Sprintf("Flag Verbose: %t", *paramVerbose))
		infoLog(fmt.Sprintf("Flag Force: %t", *paramForce))
		infoLog(fmt.Sprintf("Flag Recursive: %t", *paramRecursive))
		infoLog(fmt.Sprintf("Flag Delete Files: %t", *paramForceDelete))
		infoLog(fmt.Sprintf("Flag Path destination folder: %t", *paramPathToRestore))

		infoLog(fmt.Sprintf("Flag Version: %t", *paramVersion))
		infoLog(fmt.Sprintf("Recycle Folder: %s", RECYCLED_FOLDER))
		infoLog(fmt.Sprintf("Recycle FileDb: %s", RECYCLED_FILEDB))
	}

	if *paramListRecycleFiles {
		infoLog(fmt.Sprintf("List files: " + pathToProcess))
		getListDeletedFiles(pathToProcess, false)
		os.Exit(0)
	}

	if *paramRecoverFiles {
		infoLog(fmt.Sprintf("Recover files: " + pathToProcess))
		recoverFiles(pathToProcess, *paramPrompt, *paramPathToRestore)
		os.Exit(0)
	}

	if *paramDeleteFiles {
		infoLog(fmt.Sprintf("Delete files: %s", pathToProcess))
		var filesToDelete []string
		filesToDelete = deleteFiles(pathToProcess, *paramRecursive, filesToDelete, *paramPrompt, *paramForceDelete, *paramVerbose)
		os.Exit(0)
	}
}

func deleteFiles(pathToProcess string, recursive bool, listFiles []string, paramPrompt bool, paramForceDelete bool, paramVerbose bool) (filesToDelete []string) {

	pathToProcess = filepath.Clean(pathToProcess)
	pathDir := filepath.Dir(pathToProcess)
	baseName := filepath.Base(pathToProcess)

	if paramVerbose {
		infoLog(fmt.Sprintf("Path to process: %s", pathToProcess))
		infoLog(fmt.Sprintf("Path Dir: %s", pathDir))
		infoLog(fmt.Sprintf("Base Name: %s", baseName))
	}

	listFiles = getFilesFromFolder(pathDir, baseName, recursive, listFiles)

	infoDeletedFile := loadInfoDeletedFile()
	infoDeletedFile = appendInfoDeletedFile(listFiles, infoDeletedFile, paramPrompt)
	moveFilesToRecycle(infoDeletedFile, paramForceDelete)

	if !paramForceDelete {
		saveInfoDeletedFile(infoDeletedFile, false)
	}

	return listFiles
}

func recoverFiles(paramRecoverFiles string, paramPrompt bool, paramPathToRestore string) {
	infoDeletedFile := getListDeletedFiles(paramRecoverFiles, paramPrompt)
	moveFilesFromRecycle(infoDeletedFile, paramPathToRestore)
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

	if s, err := os.Stat(folderPath); os.IsNotExist(err) || !s.IsDir() {
		err = errors.New(fmt.Sprintf("Folder is invalid or not exists (%s).", folderPath))
		errLog(err, nil)
	}

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
					infoLog(fmt.Sprintf("Folder: %s", tmpFilePath))
					listFiles = getFilesFromFolder(tmpFilePath, "*", recursive, listFiles)
				}
			} else {
				infoLog(fmt.Sprintf("File: %s", tmpFilePath))
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

func moveFilesFromRecycle(infoDeletedFile []InfoDeletedFile, paramPathToRestore string) {
	for _, v := range infoDeletedFile {
		if v.toProcess {
			//RESTORE
			destFolder := v.pathFile
			if paramPathToRestore != "" {
				destFolder = paramPathToRestore
			}
			infoLog(fmt.Sprintf("Recovered file: %s", filepath.Join(destFolder, v.name)))
			err := os.MkdirAll(destFolder, os.ModeDir)
			if err != nil {
				errLog(err, debug.Stack())
			}
			err = os.Rename(filepath.Join(RECYCLED_FOLDER, v.uuid.String()+"_"+v.name), filepath.Join(destFolder, v.name))
			if err != nil {
				errLog(err, debug.Stack())
			}
		}
	}
}

func moveFilesToRecycle(infoDeletedFile []InfoDeletedFile, paramForceDelete bool) {
	for _, v := range infoDeletedFile {
		if v.toProcess {
			infoLog(fmt.Sprintf("Deleted file: %s", filepath.Join(v.pathFile, v.name)))
			if paramForceDelete {
				err := os.Remove(filepath.Join(v.pathFile, v.name))
				if err != nil {
					errLog(err, debug.Stack())
				}
			} else {
				err := os.Rename(filepath.Join(v.pathFile, v.name), filepath.Join(RECYCLED_FOLDER, v.uuid.String()+"_"+v.name))
				if err != nil {
					errLog(err, debug.Stack())
				}
			}

		}
	}
}

func getListDeletedFiles(filter string, paramPrompt bool) (resultInfoDeletedFile []InfoDeletedFile) {
	infoDeletedFile := loadInfoDeletedFile()

	for k, v := range infoDeletedFile {
		if b, err := filepath.Match(filter, v.name); b && errLog(err, debug.Stack()) {
			//if strings.Contains(v.name, filter) {
			infoLog(fmt.Sprintf("%d: %s  %s  %s", k, v.uuid.String()[:8], v.pathFile, v.name))
			if !paramPrompt || askForConfirmation(v.name, true) {
				infoDeletedFile[k].toProcess = true
			} else {
				infoDeletedFile[k].toProcess = false
			}
			continue
		}
		if b, err := filepath.Match(filter, v.pathFile); b && errLog(err, debug.Stack()) {
			//if strings.Contains(v.pathFile, filter) {
			infoLog(fmt.Sprintf("%d: %s  %s  %s", k, v.uuid.String()[:8], v.pathFile, v.name))
			if !paramPrompt || askForConfirmation(v.name, true) {
				infoDeletedFile[k].toProcess = true
			} else {
				infoDeletedFile[k].toProcess = false
			}
			continue
		}
		//if b, err := filepath.Match(filter, v.uuid.String()); b && errLog(err, debug.Stack()) {
		if strings.Contains(v.uuid.String(), filter) {
			infoLog(fmt.Sprintf("%d: %s  %s  %s", k, v.uuid.String()[:8], v.pathFile, v.name))
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
