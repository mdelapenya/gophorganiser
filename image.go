package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	log "github.com/sirupsen/logrus"
)

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

func processImage(p string) {
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
