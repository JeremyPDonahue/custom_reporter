package report

import (
	"os"
	"regexp"
)

// Contains is a helper function that is used to check whether a given string is present in a given slice
func Contains(sl []string, str string) bool {
	for _, x := range sl {
		if x == str {
			return true
		}
	}
	return false
}

// SliceContains is a helper function that is used to check whether a given value within a slice is present in a given slice of slices
func SliceContains(slosl [][]string, sl []string) bool {
	for _, x := range slosl {
		if x[0] == sl[0] {
			return true
		}
	}
	return false
}

// Globally defining the users homepath environment variable as well as any regexes needed throughout
var homePath, _ = os.LookupEnv("HOME")
var isItANamespace = regexp.MustCompile("app[0-9]*")
var isItACluster = regexp.MustCompile("nsk-.*-.*")
