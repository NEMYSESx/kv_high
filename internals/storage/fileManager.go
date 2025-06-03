package storage

import (
	"fmt"
	"os"
	"sync"
)

type FileManager struct {
	mu sync.Mutex
	logDir	string
	baseFileName string
	fileNumber int
	currentFile *os.File
	currentSize int64
}

func NewFileManager(logDir, baseFileName string) (*FileManager,error){
	if err := os.MkdirAll(logDir, 0755); err != nil{
		return nil, fmt.Errorf("Failed to create the log directory: %v",err)
	} 
	
	fm :=&FileManager{
		logDir: logDir,
		baseFileName: baseFileName,
		fileNumber: 0,
	}

	if err := fm.CreateNewFile(); err != nil{
		return nil ,err
	}

	return fm, nil
}


func (fm *FileManager) CreateNewFile() error{
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.currentFile != nil {
		fm.currentFile.Close()
	}

	return nil
}