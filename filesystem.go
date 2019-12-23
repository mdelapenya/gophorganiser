package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	log "github.com/sirupsen/logrus"
)

func copyFile(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists.", dst)
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

func isImage(p string) bool {
	fi, _ := os.Stat(p)
	if fi.IsDir() {
		return false
	}

	ext := strings.ToLower(filepath.Ext(p))
	if ext == ".jpg" || ext == ".jpeg" {
		return true
	}

	return false
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
		processFile(file)
	}
}

func processDirectories() {
	for _, dir := range picturesDirs {
		processDirectory(dir, targetDir)
	}
}

func processFile(p string) {
	log.Infof("Procesing file: %s", p)

	f, err := os.Open(p)
	if err != nil {
		log.Warnf("Could not open %s: %v", p, err)
		return
	}

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		log.Warnf("Could not decode EXIF from %s: %v", p, err)
		return
	}

	dt, _ := x.DateTime()
	formatDate := func(d int) string {
		if d < 10 {
			return fmt.Sprintf("0%d", d)
		}

		return fmt.Sprintf("%d", d)
	}

	newName := fmt.Sprintf(
		"%s%s%s%s%s%s_%s", formatDate(dt.Year()), formatDate(int(dt.Month())),
		formatDate(dt.Day()), formatDate(dt.Hour()), formatDate(dt.Minute()),
		formatDate(dt.Second()), filepath.Base(p))

	err = copyFile(p, path.Join(targetDir, newName), 1000)
	if err != nil {
		log.Warnf("Error copying file: %v", err)
		return
	}

	err = os.Remove(p)
	if err != nil {
		log.Warnf("Error deleting source file: %v", err)
		return
	}

	log.Infof("File %s removed from source", p)
}
