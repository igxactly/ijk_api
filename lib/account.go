package ijk_api

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	//"log"
	"errors"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"
)

type Signup struct {
	PhoneNumber string
	Password    string
	Nickname    string
}

func randSeq(n int) string {
	letters := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UTC().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func AccountRequestHandler(w http.ResponseWriter, r *http.Request) {
	s, err := httputil.DumpRequest(r, true)
	if err == nil {
		fmt.Println(string(s))
	}

	if r.Method == "GET" { // 재인증
		err = signInRequestHandler(w, r)

	} else if r.Method == "POST" { // 가입
		err = signupRequestHandler(w, r)

		/*} else if r.Method == "PUT" { // 없음
		err = signInRequestHandler(w, r)*/

	} else if r.Method == "DELETE" { // 탈퇴
		err = deactivateRequestHandler(w, r)

	} else {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		// w.Header().Set(key, value)
		// w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte{0})
	}
}

func signupRequestHandler(w http.ResponseWriter, r *http.Request) (err error) {
	fmt.Println("Signup Request Handler Called")

	// get post request body
	bs, err := ioutil.ReadAll(r.Body)
	checkErr(err)

	// parse xml from request body
	var s Signup
	err = xml.Unmarshal(bs, &s)
	checkErr(err)

	fmt.Println("XML Unmarshaled")
	fmt.Println(s)

	// check data content
	// NOT IMPLEMENTED YET
	if len(s.Nickname) == 0 || len(s.Password) == 0 || len(s.PhoneNumber) == 0 {
		return errors.New("information not sufficient")
	}

	// check duplication
	userPhoneNo := s.PhoneNumber
	userPasswd := s.Password
	userNickname := s.Nickname

	db, err := sql.Open("mysql", dbAccountStr)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	phoneNumExist := false

	qryString := "select phonenum from users" +
		" where phonenum='" + userPhoneNo + "';"

	rows, err := db.Query(qryString)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		phoneNumExist = true
	}

	err = rows.Err()
	checkErr(err)

	accessToken := ""
	// create new user
	if !phoneNumExist {
		qryString := "insert users (phonenum, name, password)" +
			" values (?, ?, ?);"

		stmt, err := db.Prepare(qryString)
		checkErr(err)

		res, err := stmt.Exec(userPhoneNo, userNickname, userPasswd)
		checkErr(err)

		id, err := res.LastInsertId()
		checkErr(err)

		// prepare access token
		accessToken = randSeq(40)

		qryString = "insert apitokens (idusers, apitokenstr)" +
			" values (?, ?);"

		stmt, err = db.Prepare(qryString)
		checkErr(err)

		res, err = stmt.Exec(id, accessToken)
		checkErr(err)

		id, err = res.LastInsertId()
		checkErr(err)

	} else { // if exists
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "user phonenum exists")

		return errors.New("existing phonenum")
	}

	bodyByteStream := []byte(accessToken)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(bodyByteStream)

	return nil
}

func signInRequestHandler(w http.ResponseWriter, r *http.Request) (err error) {
	fmt.Println("SignIn Request Handler Called")

	// get querystring from request url
	userPhoneNo := r.FormValue("phoneno")
	userPasswd := r.FormValue("passwd")

	// check data content
	// NOT IMPLEMENTED YET
	if len(userPasswd) == 0 || len(userPhoneNo) == 0 {
		return errors.New("information not sufficient")
	}

	db, err := sql.Open("mysql", dbAccountStr)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	phoneNumExist := false

	qryString := "select idusers, phonenum, password from users" +
		" where phonenum='" + userPhoneNo + "';"

	rows, err := db.Query(qryString)
	checkErr(err)
	defer rows.Close()

	var (
		id     string
		phone  string
		passwd string
	)

	for rows.Next() {
		err = rows.Scan(&id, &phone, &passwd)
		phoneNumExist = true
	}

	fmt.Println("data:", id, phone, passwd)

	err = rows.Err()
	checkErr(err)

	accessToken := ""
	// create new user
	if phoneNumExist && phone == userPhoneNo && userPasswd == passwd {
		qryString := "update apitokens set apitokenstr=?" +
			" where idusers=?;"

		// prepare access token
		accessToken = randSeq(40)

		stmt, err := db.Prepare(qryString)
		checkErr(err)

		res, err := stmt.Exec(accessToken, id)
		checkErr(err)

		numaff, err := res.RowsAffected()
		checkErr(err)

		fmt.Println("affected:", numaff)

	} else { // if exists
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad account information")

		return errors.New("account seems not exist")
	}

	bodyByteStream := []byte(accessToken)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(bodyByteStream)

	return nil
}

func deactivateRequestHandler(w http.ResponseWriter, r *http.Request) (err error) {
	fmt.Println("Deactivation Request Handler Called")

	// get querystring from request url
	userPhoneNo := r.FormValue("phoneno")
	userPasswd := r.FormValue("passwd")

	// check data content
	// NOT IMPLEMENTED YET
	if len(userPasswd) == 0 || len(userPhoneNo) == 0 {
		return errors.New("information not sufficient")
	}

	db, err := sql.Open("mysql", dbAccountStr)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	phoneNumExist := false

	qryString := "select idusers, phonenum, password from users" +
		" where phonenum='" + userPhoneNo + "';"

	rows, err := db.Query(qryString)
	checkErr(err)
	defer rows.Close()

	var (
		id     string
		phone  string
		passwd string
	)

	for rows.Next() {
		err = rows.Scan(&id, &phone, &passwd)
		phoneNumExist = true
	}

	err = rows.Err()
	checkErr(err)

	// create new user
	fmt.Println("before check account")
	fmt.Println(phoneNumExist, id, phone, '=', userPhoneNo, passwd, '=', userPasswd, err)
	if phoneNumExist && phone == userPhoneNo && passwd == userPasswd {
		qryString := "delete from users" +
			" where idusers=?;"

		stmt, err := db.Prepare(qryString)
		checkErr(err)

		res, err := stmt.Exec(id)
		checkErr(err)

		numaff, err := res.RowsAffected()
		checkErr(err)

		fmt.Println("affected:", numaff)

	} else { // if exists
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad account information")

		return errors.New("account seems not exist")
	}

	w.Write([]byte{0})

	return nil
}
