package main

import (
	"os"
)

func readFile(fname string) string {
	databyte, err := os.ReadFile(fname)
	checkErr(err)
	return string(databyte)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}

}

func openOrCreateFile(name string) (*os.File, error) {

	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil && err.Error() == "open "+name+": no such file or directory" {
		file, _ = os.Create(name)
		local_log.Printf("File Created: " + name)
		return file, nil
	}

	return file, nil

}
