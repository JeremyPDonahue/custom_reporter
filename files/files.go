package files

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

// GetFileList function takes in three string arguments: directory, toGrepFor and toExclude.
// It uses the exec package to execute a shell command that lists all the files in the directory and
// filters them based on the passed values of toGrepFor and toExclude.
// It returns a slice of strings representing the filtered list of files in the directory.
func GetFileList(directory, toGrepFor, toExclude string) []string {

	var cmd *exec.Cmd
	if toExclude == "" {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("ls %s/ | grep %s", directory, toGrepFor))
	}
	output, _ := cmd.Output()
	files := strings.Fields(string(output))
	return files
}

func WriteToFile(fileName, content string) string {
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

	return fileName
}

func structsToCSV(structSlice interface{}, fileName string) {

	structType := reflect.TypeOf(structSlice).Elem().Elem()
	var headers string
	for i := 0; i < structType.NumField(); i++ {
		headers += structType.Field(i).Name + ","
	}
	headers = headers[:len(headers)-1]
	fileName = WriteToFile(fileName, headers string)

	data := [][]string{}

	structValue := reflect.ValueOf(structSlice).Elem()

	for i := 0; i < structValue.Len(); i++ {
		row := []string{}

		for j := 0; j < structType.NumField(); j++ {
			row = append(row, fmt.Sprint(structValue.Index(i).Field(j).Interface()))
		}
	
		data = append(data, row)
	}

	file, _ := os.Open(fileName)
	defer file.Close()
	writer := csv.NewWriter(file)


	writer.Write([]string{headers})
	writer.WriteAll(data)
	writer.Flush()



	


}
