package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := openSourceFile(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := createDestinationFile(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	fileSize, err := getFileSize(srcFile)
	if err != nil {
		return err
	}

	err = validateOffset(offset, fileSize)
	if err != nil {
		return err
	}

	limit = calculateLimit(offset, limit, fileSize)

	if err = seekToOffset(srcFile, offset); err != nil {
		return err
	}

	reader := io.LimitReader(srcFile, limit)
	buffer := make([]byte, 32*1024)
	copied := int64(0)
	var mu sync.Mutex

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			mu.Lock()
			printProgressBar(copied, limit)
			mu.Unlock()
			if copied >= limit {
				break
			}
		}
	}()

	err = copyData(reader, dstFile, buffer, &copied)
	if err != nil {
		return err
	}

	mu.Lock()
	printProgressBar(limit, limit)
	mu.Unlock()
	fmt.Println()

	return nil
}

func openSourceFile(path string) (*os.File, error) {
	srcFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return srcFile, nil
}

func createDestinationFile(path string) (*os.File, error) {
	dstFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return dstFile, nil
}

func getFileSize(file *os.File) (int64, error) {
	srcInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	if !srcInfo.Mode().IsRegular() {
		return 0, ErrUnsupportedFile
	}

	return srcInfo.Size(), nil
}

func validateOffset(offset, fileSize int64) error {
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

func calculateLimit(offset, limit, fileSize int64) int64 {
	if limit == 0 || offset+limit > fileSize {
		return fileSize - offset
	}

	return limit
}

func seekToOffset(file *os.File, offset int64) error {
	_, err := file.Seek(offset, io.SeekStart)

	return err
}

func copyData(reader io.Reader, dstFile *os.File, buffer []byte, copied *int64) error {
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			_, writeErr := dstFile.Write(buffer[:n])
			if writeErr != nil {
				return writeErr
			}
			*copied += int64(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
