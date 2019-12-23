package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func copyFile(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	log.Infof("File copied from %s to %s", src, dst)

	return err
}

func processDirectory(source string, target string) {
	log.Infof("Procesing directory: %s", source)

	var files []string
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Cannot walk through the directory %s", source)
		}

		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("Cannot open file %s", path)
		}
		defer f.Close()

		if isImage(path) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Cannot walk through the directory %s", source)
	}

	for _, file := range files {
		processImage(file)
	}
}

func processDirectories() {
	for _, dir := range picturesDirs {
		processDirectory(dir, targetDir)
	}
}
