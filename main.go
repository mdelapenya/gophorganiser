package main

import (
	"errors"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
)

var picturesDirs = []string{}

func main() {
	log.Print("Hello Gophorganiser!")

	// add source directories
	for true {
		picturesDir := promptPath()
		picturesDirs = append(picturesDirs, picturesDir)

		if !promptContinue() {
			break
		}
	}

	log.Printf("We are going to process the following directories: %v", picturesDirs)
}

func pathExists(p string) error {
	for _, picturesDir := range picturesDirs {
		if p == picturesDir {
			return errors.New("Path was already added")
		}
	}

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return errors.New("Path does not exist")
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

func promptPath() string {
	prompt := promptui.Prompt{
		Label:    "Location of the pictures: ",
		Default:  "N",
		Validate: pathExists,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return result
}
