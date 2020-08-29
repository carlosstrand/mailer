package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"
)

func FindAvailableTemplates() []string {
	var paths []string
	var files []string

	root := "templates/mails"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, path := range paths {
		file := strings.Replace(filepath.Base(path), ".html", "", 1)
		files = append(files, file)
	}
	return files
}

func PrintAvailableTemplates(templates []string) {
	if len(templates) > 0 {
		fmt.Println("-----------------------------------------------------")
		fmt.Printf("There are %d templates available:\n", len(templates))
		for _, t := range templates {
			color.Green("- %s", t)
		}
		fmt.Println("-----------------------------------------------------")
	}
}
