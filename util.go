package main

import "io/ioutil"

func readFile(fname string) string {
	databyte, err := ioutil.ReadFile(fname)
	checkErr(err)
	return string(databyte)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}

}
