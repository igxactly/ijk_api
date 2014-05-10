package ijk_api

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var dbAccountStr string

func GetDBAccountString_FromFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file", err)
		os.Exit(1)
	}
	defer f.Close()

	r := bufio.NewReader(f)

	for {
		str, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("error while reading file", err)
			break
			//return err
		}

		dbAccountStr = strings.Trim(str, "\n")
		fmt.Println(dbAccountStr)

		break
	}
}
