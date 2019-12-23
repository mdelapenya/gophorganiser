package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
)

var picturesDirs = []string{}
var targetDir string

func main() {
	log.Print("Hello Gophorganiser!")

	// add source directories
	for true {
		picturesDir := promptPath("Location of the pictures: ", "N")
		picturesDirs = append(picturesDirs, picturesDir)

		if !promptContinue() {
			break
		}
	}

	log.Printf("We are going to process the following directories: %v", picturesDirs)
	targetDir = promptPath("Now indicate the target directory where to store the processed images: ", "")
	log.Printf("Target directory: %v", targetDir)
	processDirectories()
}

func pathAlreadyAdded(p string) error {
	// always check for absolute file path, to avoid duplicates
	filePath, err := filepath.Abs(p)
	if err != nil {
		log.Fatalf("Cannot get absolute path for %s: %v\n", p, err)
	}

	for _, picturesDir := range picturesDirs {
		if filePath == picturesDir {
			return errors.New("Already added")
		}
	}

	return nil
}

func pathExists(p string) error {
	err := pathAlreadyAdded(p)
	if err != nil {
		return err
	}

	fi, err := os.Stat(p)
	if os.IsNotExist(err) {
		return errors.New("Path does not exist")
	} else if !fi.IsDir() {
		return errors.New("Path is not a directory")
	}

	return nil
}

func promptContinue() bool {
	validate := func(input string) error {
		validInputs := []string{
			"y", "yes", "n", "no", "",
		}
		for _, b := range validInputs {
			if b == strings.ToLower(input) {
				return nil
			}
		}
		return errors.New("Accepted values: y, yes, no, n")
	}

	prompt := promptui.Prompt{
		Label:    "Do you want to add more directories? (y/N)",
		Default:  "N",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	lower := strings.ToLower(result)
	if "y" == lower || "yes" == lower {
		return true
	}

	return false
}

func promptPath(label string, defaultValue string) string {
	prompt := promptui.Prompt{
		Label:    label,
		Default:  defaultValue,
		Validate: pathExists,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	filePath, err := filepath.Abs(result)
	if err != nil {
		log.Fatalf("Cannot get absolute path for %v\n", err)
	}

	return filePath
}
