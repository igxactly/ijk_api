package ijk_api

import (
	"log"
	"strings"
)

type SignInfo struct {
	PhoneNumber string
	SignKey     string
}

type Machine struct {
	MachineCode  string
	MachineCount string
}

type OwnMachines struct {
	Machine []Machine `xml:""`
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

func parseOwnMachine(midliststr string, mnumliststr string) (omach *OwnMachines, err error) {
	midlist := strings.Split(midliststr, ",")
	mnumlist := strings.Split(mnumliststr, ",")

	if len(midlist) == len(mnumlist) {
		mlist := []Machine{}

		for i := 0; i < len(midlist); i++ {
			mid, mnum := strings.TrimSpace(midlist[i]), strings.TrimSpace(mnumlist[i])
			mlist = append(mlist, Machine{mid, mnum})
		}

		omach = &OwnMachines{mlist}
	} else {
		omach = nil
	}

	return
}
