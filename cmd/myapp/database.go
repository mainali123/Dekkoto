package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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
	queryToInsertVideo := "INSERT INTO videos (Title, Description, URL, ThumbnailURL, UploaderID, UploadDate, ViewsCount, LikesCount, Duration, CategoryID, GenreID) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(queryToInsertVideo, title, desc, videoUrl, thumbnailUrl, uploaderId, uploadDate, 0, 0, duration, categoryId, genre)
	if err != nil {
		return err
	}
	return nil
}

// Function to fetch CategoryID based on CategoryName
func (db *databaseConn) getCategoryID(categoryName string) (int, error) {
	fmt.Println("From database: ", categoryName)
	var categoryID int
	query := "SELECT CategoryID FROM categories WHERE CategoryName = ?"
	fmt.Println("Category query: ", query)
	err := db.DB.QueryRow(query, categoryName).Scan(&categoryID)
	fmt.Println("From database (cat id1): ", err)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, sql.ErrNoRows
		}
		return 0, err
	}
	fmt.Println("From database (cat id): ", categoryID)
	return categoryID, nil
}

// Function to get category name based on category id
func (db *databaseConn) getCategoryName(categoryID int) (string, error) {
	var categoryName string
	query := "SELECT CategoryName FROM categories WHERE CategoryID = ?"
	err := db.DB.QueryRow(query, categoryID).Scan(&categoryName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", sql.ErrNoRows
		}
		return "", err
	}
	return categoryName, nil
}

// Function to fetch GenreID based on GenreName
func (db *databaseConn) getGenreID(genreName string) (int, error) {
	fmt.Println("From database: ", genreName)
	var genreID int
	query := "SELECT GenreID FROM genres WHERE GenreName = ?"
	fmt.Println("Genre query: ", query)
	err := db.DB.QueryRow(query, genreName).Scan(&genreID)
	fmt.Println("From database (genre id1): ", err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, sql.ErrNoRows
		}
		return 0, err
	}
	fmt.Println("From database (genre id): ", genreID)
	return genreID, nil
}

// Function to get genre name based on genre id
func (db *databaseConn) getGenreName(genreID int) (string, error) {
	var genreName string
	query := "SELECT GenreName FROM genres WHERE GenreID = ?"
	err := db.DB.QueryRow(query, genreID).Scan(&genreName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", sql.ErrNoRows
		}
		return "", err
	}
	return genreName, nil
}

type VideoDesc struct {
	VideoID       int
	Title         string
	Description   string
	URL           string
	ThumbnailURL  string
	UploaderID    int
	UploadDate    string
	ViewsCount    int
	LikesCount    int
	DislikesCount int
	Duration      string
	CategoryID    int
	GenreID       int
}

func (db *databaseConn) videoDescForTable(userID int) ([]VideoDesc, error) {
	query := "SELECT * FROM videos WHERE UploaderID = ?"
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []VideoDesc
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func (db *databaseConn) videosBrowser() ([]VideoDesc, error) {
	query := "SELECT * FROM videos"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []VideoDesc
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func (db *databaseConn) videoDescForEdit(videoID int, title string, description string, category int, genre int) error {
	// Update the video details in the database based on the videoID
	query := "UPDATE videos SET Title = ?, Description = ?, CategoryID = ?, GenreID = ? WHERE VideoID = ?"
	_, err := db.DB.Exec(query, title, description, category, genre, videoID)
	if err != nil {
		return err
	}
	return nil
}

func (db *databaseConn) deleteVideo(videoID int) error {
	query := "DELETE FROM videos WHERE VideoID = ?"
	_, err := db.DB.Exec(query, videoID)
	if err != nil {
		return err
	}
	return nil
}

func (db *databaseConn) deleteVideoFromFile(videoID int) (string, string, error) {
	// URL and ThumbnailURL of the video to be deleted
	var videoURL string
	var thumbnailURL string
	query := "SELECT URL, ThumbnailURL FROM videos WHERE VideoID = ?"
	err := db.DB.QueryRow(query, videoID).Scan(&videoURL, &thumbnailURL)
	if err != nil {
		return "", "", err
	}
	return videoURL, thumbnailURL, nil
}

func (db *databaseConn) recentlyAddedVideos() ([]VideoDesc, error) {
	query := "SELECT * FROM videos ORDER BY UploadDate DESC LIMIT 10"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var videos []VideoDesc
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func (db *databaseConn) recommendedVideos() ([]VideoDesc, error) {
	// show random 10 videos
	query := "SELECT * FROM videos ORDER BY RAND() LIMIT 1"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var videos []VideoDesc
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func (db *databaseConn) weeklyTop() ([]VideoDesc, error) {
	query := "SELECT V.VideoID, V.Title, V.Description, V.URL, V.ThumbnailURL, V.UploaderID, V.UploadDate, V.ViewsCount, V.LikesCount, V.DislikesCount, V.Duration, V.CategoryID, V.GenreID FROM Videos V JOIN VideoActions VA ON V.VideoID = VA.VideoID WHERE VA.ActionsDate BETWEEN CURRENT_DATE - INTERVAL DAYOFWEEK(CURRENT_DATE) + 6 DAY AND CURRENT_DATE ORDER BY V.ViewsCount DESC LIMIT 10"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var videos []VideoDesc
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func (db *databaseConn) videoActions(videoID int, userID int) error {
	// check if the user have already action on the video
	query := "SELECT * FROM videoactions WHERE VideoID = ? AND UserID = ?"
	rows, err := db.DB.Query(query, videoID, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		// user have already action on the video
		// update the ActionDate to current date
		currentDate := time.Now().Format("2006-01-02") // format the date as YYYY-MM-DD
		currentTime := time.Now().Format("15:04:05")   // format the time as HH:MM:SS
		// update the ActionDate and ActionTime to current date and time
		updateQuery := "UPDATE videoactions SET ActionsDate = ?, ActionTime = ? WHERE VideoID = ? AND UserID = ?"
		_, err := db.DB.Exec(updateQuery, currentDate, currentTime, videoID, userID)
		if err != nil {
			return err
		}
	} else {
		// user has not actioned on the video, so add the video to the VideoActions table
		currentDate := time.Now().Format("2006-01-02") // format the date as YYYY-MM-DD
		currentTime := time.Now().Format("15:04:05")   // format the time as HH:MM:SS
		insertQuery := "INSERT INTO videoactions (UserID, VideoID, Watching, ActionsDate, ActionTime) VALUES (?, ?, 1, ?, ?)"
		_, err := db.DB.Exec(insertQuery, userID, videoID, currentDate, currentTime)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *databaseConn) continueWatching(userID int) ([]VideoDesc, error) {
	rows, err := db.DB.Query("SELECT VideoID FROM videoactions WHERE UserID = ? AND Watching = 1 ORDER BY ActionsDate DESC, ActionTime DESC LIMIT 10", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videoIDs []int
	for rows.Next() {
		var videoID int
		if err := rows.Scan(&videoID); err != nil {
			return nil, err
		}
		videoIDs = append(videoIDs, videoID)
	}

	if len(videoIDs) == 0 {
		return []VideoDesc{}, nil
	}

	videoIDStrs := make([]string, len(videoIDs))
	for i, videoID := range videoIDs {
		videoIDStrs[i] = strconv.Itoa(videoID)
	}
	videoIDsStr := strings.Join(videoIDStrs, ",")

	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s) ORDER BY FIELD(VideoID, %s)", videoIDsStr, videoIDsStr))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []VideoDesc
	for rows.Next() {
		var video VideoDesc
		if err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *databaseConn) caroselSlide() ([]VideoDesc, error) {
	// show random 10 videos
	query := "SELECT * FROM videos ORDER BY RAND() LIMIT 10"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var videos []VideoDesc
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}
