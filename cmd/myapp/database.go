package main

import (
	"database/sql"
	"errors"
)

type databaseConn struct {
	DB *sql.DB
}

func (db *databaseConn) registerUser(userName string, email string, password string) error {
	queryToCheckIfUserExists := "SELECT * FROM users WHERE Email = ?"
	rows, err := db.DB.Query(queryToCheckIfUserExists, email)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		// User with the provided email already exists in the database
		return errors.New("user already exists")
	}

	queryToInsertUser := "INSERT INTO users (UserName, Email, Password) VALUES (?, ?, ?)"
	_, err = db.DB.Exec(queryToInsertUser, userName, email, password)
	if err != nil {
		return err
	}
	return nil
}

func (db *databaseConn) loginUser(email string, password string) error {
	queryToLoginUser := "SELECT * FROM users WHERE Email = ? AND Password = ?"
	rows, err := db.DB.Query(queryToLoginUser, email, password)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		// User with the provided email and password does not exist in the database
		return errors.New("user does not exist")
	}
	return nil
}
