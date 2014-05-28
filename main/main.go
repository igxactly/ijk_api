package main

import (
	"fmt"
	"ijk_api/lib"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage:", os.Args[0], "<dbInfoFile>")
		os.Exit(1)
	}

	ijk_api.GetDBAccountString_FromFile(os.Args[1])

	// http.HandleFunc("/", foo)
	http.HandleFunc("/ijk/test/", ijk_api.OtherRequestHandler)

	http.HandleFunc("/ijk/account/", ijk_api.AccountRequestHandler)
	http.HandleFunc("/ijk/friendlist/", ijk_api.FriendListRequestHandler)

	http.HandleFunc("/ijk/profile/", ijk_api.ProfileRequestHandler)
	http.HandleFunc("/ijk/profileimg/", ijk_api.OtherRequestHandler)
	http.HandleFunc("/ijk/profilethumb/", ijk_api.OtherRequestHandler)

	http.HandleFunc("/ijk/settings/", ijk_api.OtherRequestHandler)

	http.ListenAndServe(":3000", nil)
}
