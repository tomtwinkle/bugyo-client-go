package main

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const ReadmeTemplateFile = "./README.md.tmpl"
const ReadmeFile = "../README.md"

// export GOOS=linux
// export GOARCH=amd64
// go build -o .github/readme_update .github/readme_update.go
func main() {
	args := os.Args[1:]
	if len(args)%2 == 1 {
		log.Fatal("readme_update [template key1] [val1]  [template key2] [val2] ...")
	}
	tmplData := make(map[string]string)
	var key string
	for i, v := range args {
		if i%2 == 0 {
			key = v
		} else {
			tmplData[key] = v
		}
	}
	log.Println(tmplData)
	if err := run(tmplData); err != nil {
		log.Fatal(err)
	}
}

func run(tmplData map[string]string) error {
	readmeTemplateFilePath, exists := templateFilePath()
	if !exists {
		log.Fatalf("%s not exists", readmeTemplateFilePath)
	}
	readmeFilePath, _ := readmeFilePath()

	log.Println("start generate README")

	t, err := template.ParseFiles(readmeTemplateFilePath)
	if err != nil {
		return err
	}
	f, err := os.Create(readmeFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := t.Execute(f, tmplData); err != nil {
		return err
	}
	return nil
}

func templateFilePath() (string, bool) {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDirPath := filepath.Dir(execPath)
	readmeTemplateFilePath := filepath.Join(execDirPath, ReadmeTemplateFile)
	if _, err := os.Stat(readmeTemplateFilePath); err != nil {
		if os.IsNotExist(err) {
			return readmeTemplateFilePath, false
		}
		log.Fatal(err)
	}
	return readmeTemplateFilePath, true
}

func readmeFilePath() (string, bool) {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDirPath := filepath.Dir(execPath)
	readmeFilePath := filepath.Join(execDirPath, ReadmeFile)
	if _, err := os.Stat(readmeFilePath); err != nil {
		if os.IsNotExist(err) {
			return readmeFilePath, false
		}
		log.Fatal(err)
	}
	return readmeFilePath, true
}
