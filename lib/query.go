package ijk_api

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

// Phoneno + Passwd
func GetUserID_withPasswd(db *sql.DB, phone string, passwd string) (id string, err error) {

	return "", errors.New("")
}

// Phoneno + Signkey
func GetUserID_withSignkey(db *sql.DB, userPhoneNo string, userSignkey string) (id string, err error) {
	signkeyMatch := false

	qryString := "select idusers, apitokenstr from apitokens" +
		" where idusers in" +
		" (select idusers from users" +
		" where phonenum='" +
		userPhoneNo + "');"

	rows, err := db.Query(qryString)
	checkErr(err)
	defer rows.Close()

	var (
		signkey string
	)
	for rows.Next() {
		err := rows.Scan(&id, &signkey)
		checkErr(err)

		if signkey == userSignkey {
			signkeyMatch = true
		}
	}
	err = rows.Err()
	checkErr(err)

	if signkeyMatch {
		return id, nil
	} else {
		return "", errors.New("")
	}
}
