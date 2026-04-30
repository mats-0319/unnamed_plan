package mlog

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type HandlerWriter struct { // only one instance
	Writer io.Writer

	File     *os.File
	FileName string
	Size     int64
	MaxSize  int64 // single log file max size, unit: MB

	mu sync.Mutex
}

func (hw *HandlerWriter) New(fileName string, maxSize int64) (err error) {
	if hw == nil {
		err = errors.New("nil HandlerWriter")
		return
	}

	fileIns, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("open log file failed, error:", err)
		return
	}

	fileInfo, err := os.Stat(fileName) // 写在打开文件后面，因为可能还没有创建日志文件
	if err != nil {
		log.Println("get file info failed, error:", err)
		return
	}

	hw.Writer = io.MultiWriter(fileIns, os.Stdout)
	hw.File = fileIns
	hw.FileName = fileName
	hw.Size = fileInfo.Size()
	hw.MaxSize = maxSize

	return
}

func (hw *HandlerWriter) Write(logBytes []byte) (err error) {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	if _, err = hw.Writer.Write(logBytes); err != nil {
		return
	}

	hw.Size += int64(len(logBytes))

	if hw.Size >= hw.MaxSize {
		hw.rotate()
	}

	return
}

func (hw *HandlerWriter) rotate() {
	hw.close()

	timeStr := time.Now().Format("20060102-150405.000")
	newName := fmt.Sprintf("%s.log", timeStr)
	if os.Rename(hw.FileName, newName) != nil {
		return
	}

	go compressFile(newName)

	_ = hw.New(hw.FileName, hw.MaxSize)
}

func (hw *HandlerWriter) close() {
	if hw != nil {
		_ = hw.File.Close()
	}
}

func compressFile(fileName string) {
	gzFile, err := os.Create(fileName + ".gz")
	if err != nil {
		return
	}
	defer gzFile.Close()

	gzWriter := gzip.NewWriter(gzFile)
	defer gzWriter.Close()

	srcFile, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer srcFile.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return
	}

	_ = os.Remove(fileName)
}
