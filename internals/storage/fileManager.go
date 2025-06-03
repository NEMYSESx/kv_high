package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
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

	filename := fmt.Sprintf("%s.%d.log",fm.baseFileName,fm.fileNumber)
	filePath := filepath.Join(fm.logDir,filename)

	file, err := os.OpenFile(filePath,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0644)
	if err !=nil {
		fmt.Errorf("Failed tot create the file: %v",err)
	}

	stat, err := file.Stat()
	if err != nil{
		file.Close()
		return fmt.Errorf("failed to get file stats: %v",err)
	}

	fm.currentFile = file
	fm.currentSize = stat.Size()
	fm.fileNumber++

	return nil
}

func (fm *FileManager) CheckSpaceLeft(dataSize int) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	requiredSize := int64(recordHeaderSize + dataSize)
	return fm.currentSize+requiredSize <= recordSize
}

func (fm *FileManager) RotateFile () error{
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.currentFile != nil {
		return fmt.Errorf("no current file to rotate")
	}

	oldFileName := fm.currentFile.Name()
	if err := fm.currentFile.Close(); err != nil {
		return fmt.Errorf("failed to close current file: %v", err)
	}

	archiveDir := filepath.Join(fm.logDir,"archive")
	if err := os.MkdirAll(archiveDir,0755); err != nil{
		return fmt.Errorf("failed to create archive dictonary: %v",err)
	}

	timestamp := time.Now().Format("20060102_150405")
	archiveFileName := fmt.Sprintf("%s_%s.log", 
		filepath.Base(oldFileName[:len(oldFileName)-4]), timestamp)
	archivePath := filepath.Join(archiveDir, archiveFileName)

	if err := os.Rename(oldFileName, archivePath); err != nil {
		return fmt.Errorf("failed to move file to archive: %v", err)
	}

	return fm.CreateNewFile()
}

func (fm *FileManager) GetCurrentFile()(*os.File,error){
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.currentFile != nil{
		return nil, fmt.Errorf("no current file available")
	}

	return fm.currentFile, nil
}

func (fm *FileManager) UpdateSize(deltaSize int64) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	fm.currentSize += deltaSize
}

func (fm *FileManager) Close() error {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	if fm.currentFile != nil {
		return fm.currentFile.Close()
	}
	return nil
}