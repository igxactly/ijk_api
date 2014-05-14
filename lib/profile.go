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
	// "time"
)

func ProfileRequestHandler(w http.ResponseWriter, r *http.Request) {
	s, err := httputil.DumpRequest(r, true)
	if err == nil {
		fmt.Println(string(s))
	}

	if r.Method == "GET" { // 친구목록 수신
		err = ProfileListRetrieveRequestHandler(w, r)

	} else if r.Method == "PUT" { // 친구목록 업로드
		err = ProfileListUpdateRequestHandler(w, r)

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

	db, err := sql.Open("mysql", dbAccountStr)
	checkErr(err)
	defer db.Close()

	// check account signkey
	id, err := GetUserID_withSignkey(db, userPhoneNo, userSignkey)

	//query profiles of friends
	var (
		profile    Profile
		phone      string
		nick       string
		userstatus string
		region     string
		wperiod    int
		omachines  *OwnMachines
		x          []byte
	)

	if err == nil {
		qryString := "select phonenum, name, region, workingperiod, currentstatus from users" +
			" where phonenum =" + phone + ");"

		rows, err := db.Query(qryString)
		checkErr(err)
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&phone, &nick, &region, &wperiod, &userstatus)
			checkErr(err)

			omachines = nil //&OwnMachines{[]Machine{}}
			flist = append(flist, Profile{phone, nick, userstatus, region, wperiod, omachines})
		}

		err = rows.Err()
		checkErr(err)

		flrrResponce := FriendListRetrieveResult{flist}
		x, err = xml.MarshalIndent(flrrResponce, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "user phonenum exists")
		return errors.New("existing phonenum")
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

		var flUpd FriendListUpdate
		fmt.Println("going to parse")
		err = xml.Unmarshal(bs, &flUpd)
		fmt.Println("just parsed")
		//checkErr(err)

		fmt.Println("XML Unmarshaled")
		//fmt.Println(flUpd)

		list := ""

		fmt.Println("before loop", flUpd)
		for idx, e := range flUpd.NewFriends {
			fmt.Println("in loop" + e)
			if idx > 0 {
				list += ","
			}
			list += "'" + e + "'"
		}

		fmt.Println("list: " + list)

		qryString := "insert into friendrelation (fk_idusers, fk_idusers_friend)" +
			" (select a.idusers, b.idusers" +
			" from users a, users b" +
			" where a.phonenum='" + userPhoneNo + "' and b.phonenum in (" + list + "));"

		stmt, err := db.Prepare(qryString)
		checkErr(err)

		res, err := stmt.Exec()
		checkErr(err)

		_, err = res.LastInsertId()
		checkErr(err)

		list = ""

		for idx, e := range flUpd.DelFriends {
			if idx > 0 {
				list += ","
			}
			list += "'" + e + "'"
		}

		fmt.Println("list2:" + list)

		qryString = "SET SQL_SAFE_UPDATES = 0;"

		stmt, err = db.Prepare(qryString)
		checkErr(err)

		res, err = stmt.Exec()
		checkErr(err)

		_, err = res.LastInsertId()
		checkErr(err)

		qryString = " delete from friendrelation" +
			" where idfriendrelation in (select temp.idfriendrelation from (" +
			" select idfriendrelation" +
			" from friendrelation" +
			" where" +
			" fk_idusers in (" +
			" select idusers from users" +
			" where phonenum='" + userPhoneNo + "'" +
			" )" +
			" and" +
			" fk_idusers_friend in (" +
			" select idusers from users" +
			" where phonenum in (" + list +
			" ) )" +
			" ) as temp" +
			" );"

		fmt.Println(qryString)

		stmt, err = db.Prepare(qryString)
		checkErr(err)

		res, err = stmt.Exec()
		checkErr(err)

		_, err = res.LastInsertId()
		checkErr(err)

		qryString = " SET SQL_SAFE_UPDATES = 1;"

		stmt, err = db.Prepare(qryString)
		checkErr(err)

		res, err = stmt.Exec()
		checkErr(err)

		_, err = res.LastInsertId()
		checkErr(err)

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
