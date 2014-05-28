package ijk_api

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	// "log"
	"net/http"
	"net/http/httputil"
	"strings"
	// "time"
)

type ProfileUpdate struct {
	Profile *Profile `xml:"Profile,omitempty"`
}

func ProfileRequestHandler(w http.ResponseWriter, r *http.Request) {
	s, err := httputil.DumpRequest(r, true)
	if err == nil {
		fmt.Println(string(s))
	}

	if r.Method == "GET" { // 프로필 수신
		err = ProfileRetrieveRequestHandler(w, r)

	} else if r.Method == "PUT" { // 프로필 업로드
		err = ProfileUpdateRequestHandler(w, r)

	} else {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		// w.Header().Set(key, value)
		// w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte{0})
	}
}

func ProfileRetrieveRequestHandler(w http.ResponseWriter, r *http.Request) (err error) {

	userPhoneNo := r.FormValue("phoneno")
	userSignkey := r.FormValue("signkey")

	path := r.URL.Path
	fmt.Println(path)

	rqPrfNumList := strings.SplitAfterN(path, "/", 4)

	rqProfileNum := strings.Trim(rqPrfNumList[len(rqPrfNumList)-1], "/ ")

	db, err := sql.Open("mysql", dbAccountStr)
	checkErr(err)
	defer db.Close()

	// check account signkey
	_, err = GetUserID_withSignkey(db, userPhoneNo, userSignkey)

	//query profiles of friends
	var (
		profile        Profile
		phone          string
		nick           string
		userstatus     string
		region         string
		wperiod        int
		machidliststr  string
		nummachliststr string
		omachines      *OwnMachines
		x              []byte
	)

	if err == nil {
		qryString := "select" +
			" u.phonenum, u.name, u.region, u.workingperiod," +
			" u.currentstatus, m.idmachinelist, m.nummachinelist" +
			" from" +
			" users" +
			" as u," +
			" (select" +
			" idusers," +
			" group_concat(idmachine separator ', ')" +
			" as idmachinelist," +
			" group_concat(nummachine separator ', ')" +
			" as nummachinelist" +
			" from ownmachines group by idusers" +
			" ) as m" +
			" where" +
			" u.idusers" +
			" in (select idusers from users where phonenum='" + rqProfileNum + "')" +
			" and u.idusers=m.idusers;"

		rows, err := db.Query(qryString)
		checkErr(err)
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&phone, &nick, &region, &wperiod,
				&userstatus, &machidliststr, &nummachliststr)
			checkErr(err)

			omachines, _ = parseOwnMachine(machidliststr, nummachliststr)

			profile = Profile{phone, nick, userstatus, region, wperiod, omachines}
		}

		err = rows.Err()
		checkErr(err)

		response := profile
		x, err = xml.MarshalIndent(response, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "profile access cannot be done")
		return errors.New("profile access cannot be done")
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)

	return nil
}

func ProfileUpdateRequestHandler(w http.ResponseWriter, r *http.Request) (err error) {
	// get post request body
	fmt.Println("before read")
	bs, err := ioutil.ReadAll(r.Body)
	checkErr(err)

	path := r.URL.Path
	fmt.Println(path)

	userPhoneNo := r.FormValue("phoneno")
	userSignkey := r.FormValue("signkey")

	db, err := sql.Open("mysql", dbAccountStr)
	checkErr(err)
	defer db.Close()

	// check account signkey
	_, err = GetUserID_withSignkey(db, userPhoneNo, userSignkey)

	if err == nil {
		// unmarshaling xml
		// fmt.Println("Signup Request Handler Called")

		fmt.Println(string(bs))

		// parse xml from request body
		fmt.Println("before parse")

		var (
			profileUpdate ProfileUpdate
			// profile       Profile
			// phone         string
			// nick          string
			// userstatus    string
			// region        string
			// wperiod       int
			// omachines     *OwnMachines
		)

		fmt.Println("going to parse")
		err = xml.Unmarshal(bs, &profileUpdate)
		fmt.Println("just parsed")
		//checkErr(err)

		fmt.Println("XML Unmarshaled")
		fmt.Println(profileUpdate)

		var (
			qryString string
			stmt      *sql.Stmt
			res       sql.Result
		)

		qryString = ""

		if profileUpdate.Profile.UserStatus == "free" {
			qryString = "update users set currentstatus='free'" +
				" where phonenum='" + userPhoneNo + "';"

		} else if profileUpdate.Profile.UserStatus == "busy" {
			qryString = "update users set currentstatus='busy'" +
				" where phonenum='" + userPhoneNo + "';"

		}

		if len(qryString) > 0 {
			stmt, err = db.Prepare(qryString)
			checkErr(err)

			res, err = stmt.Exec()
			checkErr(err)

			_, err = res.LastInsertId()
			checkErr(err)
		}

		fmt.Println(qryString)

	} else { // if exists
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "account infomation error")

		return errors.New("account infomation error")
	}

	bodyByteStream := []byte{0}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(bodyByteStream)

	return nil
}
