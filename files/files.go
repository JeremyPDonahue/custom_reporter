package files

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// GetFileList function takes in three string arguments: directory, toGrepFor and toExclude.
// It uses the exec package to execute a shell command that lists all the files in the directory and
// filters them based on the passed values of toGrepFor and toExclude.
// It returns a slice of strings representing the filtered list of files in the directory.
func GetFileList(directory, toGrepFor string) []string {

	var cmd *exec.Cmd

	cmd = exec.Command("sh", "-c", fmt.Sprintf("ls %s/ | grep %s", directory, toGrepFor))

	output, _ := cmd.Output()
	files := strings.Fields(string(output))
	return files
}

func WriteToFile(fileName, content string) {
	now := time.Now()
	year, month, day := now.Date()
	fileName = fmt.Sprintf("%s_%d_%d_%d.csv", fileName, year, month, day)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Sprintf("Unable to open %s", fileName)
		fmt.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(content + "\n"); err != nil {
		fmt.Sprintf("Unable to write to %s", fileName)
		fmt.Println(err)
	}
}

func HeadersFromStructFields(struc interface{}) string {
	var headers string
	for i := 0; i < reflect.TypeOf(struc).NumField(); i++ {
		headers += reflect.TypeOf(struc).Field(i).Name + ","
	}
	headers = strings.TrimSuffix(headers, ",")
	return headers
}

func StructsToAPICSV(structs []interface{}, fileName string) error {
	now := time.Now()
	year, month, day := now.Date()
	fileName = fmt.Sprintf("%s_%d_%d_%d.csv", fileName, year, month, day)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Trouble creating file")
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// writer.Write(headers)
	var i int
	for _, struc := range structs {
		var fields []string
		val := reflect.ValueOf(struc)
		fmt.Sprintf("Let's count the structs: %d", i)
		i++
		for i := 0; i < val.NumField(); i++ {
			fmt.Sprintf("Let's count the values: %d", i)

			valField := val.Field(i)
			fields = append(fields, valField.Interface().(string))
		}

		writer.Write(fields)
	}

	return nil
}

// This function takes in a slice of strings representing file paths and returns a new slice of strings containing the base of each file path.
// The base of a file path is the last element of the path, which is usually the file name.
// For example, if the input is a slice of paths ["/Users/yx4h/.kube/nsk-beet-prod", "/Users/yx4h/.kube/nsk-curry-nonprod"],
// the output will be a slice containing ["nsk-beet-prod", "nsk-beet-nonprod"]
func ExtractLastPaths(paths []string) []string {
	var lastPaths []string
	for _, path := range paths {
		lastPaths = append(lastPaths, filepath.Base(path))
	}
	return lastPaths
}
