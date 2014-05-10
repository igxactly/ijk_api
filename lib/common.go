package ijk_api

import (
	"log"
)

type SignInfo struct {
	PhoneNumber string
	SignKey     string
}

type Machine struct {
	MachineCode  string
	MachineCount int
}

type OwnMachines struct {
	Test []Machine `xml:""`
}

type Profile struct {
	PhoneNumber   string
	NickName      string
	UserStatus    string
	Region        string       `xml:"Region,omitempty"`
	WorkingPeriod int          `xml:"WorkingPeriod,omitempty"`
	OwnMachines   *OwnMachines `xml:"OwnMachines,omitempty"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkErrPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
