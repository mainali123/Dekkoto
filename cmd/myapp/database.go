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

func (db *databaseConn) userId(email string) (int, error) {
	queryToGetUserId := "SELECT UserID FROM users WHERE Email = ?"
	rows, err := db.DB.Query(queryToGetUserId, email)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var userId int
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			return 0, err
		}
	}
	return userId, nil
}

func (db *databaseConn) uploadVideo(title string, desc string, videoUrl string, thumbnailUrl string, uploaderId string, uploadDate string, duration string, categoryId int, genre int) error {
	queryToInsertVideo := "INSERT INTO videos (Title, Description, URL, ThumbnailURL, UploaderID, UploadDate, Duration, CategoryID, GenreID) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(queryToInsertVideo, title, desc, videoUrl, thumbnailUrl, uploaderId, uploadDate, duration, categoryId, genre)
	if err != nil {
		return err
	}
	return nil
}

// Function to fetch CategoryID based on CategoryName
func (db *databaseConn) getCategoryID(categoryName string) (int, error) {
	var categoryID int
	query := "SELECT CategoryID FROM categories WHERE CategoryName = ?"
	err := db.DB.QueryRow(query, categoryName).Scan(&categoryID)
	if err != nil {
		return 0, err
	}
	return categoryID, nil
}

// Function to fetch GenreID based on GenreName
func (db *databaseConn) getGenreID(genreName string) (int, error) {
	var genreID int
	query := "SELECT GenreID FROM genres WHERE GenreName = ?"
	err := db.DB.QueryRow(query, genreName).Scan(&genreID)
	if err != nil {
		return 0, err
	}
	return genreID, nil
}
