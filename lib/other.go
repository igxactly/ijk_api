package ijk_api

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	// "time"
)

func OtherRequestHandler(w http.ResponseWriter, r *http.Request) {
	s, err := httputil.DumpRequest(r, true)
	if err == nil {
		fmt.Println(string(s))
	}

	b := []byte{}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(b)

	// userPhoneNo := r.FormValue("phoneno")

	// res := "Request: " + time.Now().String() + "Phone:" + userPhoneNo

	// fmt.Println(res)

	// data := []byte(res)
	// w.Header().Set("Content-Type", "text/plain")
	// w.Write(data)
}
