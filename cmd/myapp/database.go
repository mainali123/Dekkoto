package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

type databaseConn struct {
	DB *sql.DB
}

// registerUser is a function that registers a new user in the database.
// It takes in three parameters: userName, email, and password.
// It first checks if a user with the provided email already exists in the database.
// If the user exists, it returns an error.
// If the user does not exist, it inserts a new user record into the database with the provided userName, email, and password.
// It returns an error if the insertion fails.
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

	queryToGetUserId := "SELECT UserID FROM users WHERE Email = ?"
	rows, err = db.DB.Query(queryToGetUserId, email)
	if err != nil {
		return err
	}
	defer rows.Close()

	var userId int
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			return err
		}
	}

	// Insert into useraccesslevels
	queryToInsertUserAccessLevels := "INSERT INTO useraccesslevels (UserID) VALUES (?)"
	_, err = db.DB.Exec(queryToInsertUserAccessLevels, userId)
	if err != nil {
		return err
	}

	return nil
}

// loginUser is a function that logs in a user.
// It takes in two parameters: email and password.
// It checks if a user with the provided email and password exists in the database.
// If the user exists, it returns nil.
// If the user does not exist, it returns an error.
func (db *databaseConn) loginUser(email string, password string) error {
	// Query the database for the user with the provided email
	row := db.DB.QueryRow("SELECT password FROM users WHERE email = ?", email)

	// We create a variable to hold the hashed password from the database
	var hashedPassword string

	// Scan the result into our hashedPassword variable
	if err := row.Scan(&hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no rows were returned from the query, it means the provided email does not exist in the database
			return errors.New("the provided email does not exist in our records")
		}
		// If another error occurred, return it
		return err
	}

	// At this point, we have the hashed password, and the user-provided password. We will use bcrypt to compare the passwords
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		// If the passwords do not match, return an error
		return errors.New("the provided password is incorrect")
	}

	// If we've reached this point, it means the user-provided password matches the hashed password in the database. The user is authenticated!
	return nil
}

// userId is a function that retrieves the user ID of a user.
// It takes in one parameter: email.
// It queries the database for the user ID of the user with the provided email.
// It returns the user ID and nil if the user exists.
// It returns 0 and an error if the user does not exist or if the query fails.
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

// uploadVideo is a function that uploads a video.
// It takes in several parameters including title, desc, videoUrl, thumbnailUrl, uploaderId, uploadDate, duration, categoryId, and genre.
// It inserts a new video record into the database with the provided parameters.
// It returns an error if the insertion fails.
func (db *databaseConn) uploadVideo(title string, desc string, videoUrl string, thumbnailUrl string, uploaderId string, uploadDate string, duration string, categoryId int, genre int) error {
	queryToInsertVideo := "INSERT INTO videos (Title, Description, URL, ThumbnailURL, UploaderID, UploadDate, ViewsCount, LikesCount, Duration, CategoryID, GenreID) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(queryToInsertVideo, title, desc, videoUrl, thumbnailUrl, uploaderId, uploadDate, 0, 0, duration, categoryId, genre)
	if err != nil {
		return err
	}
	return nil
}

// getCategoryID is a function that retrieves the category ID of a category.
// It takes in one parameter: categoryName.
// It queries the database for the category ID of the category with the provided categoryName.
// It returns the category ID and nil if the category exists.
// It returns 0 and an error if the category does not exist or if the query fails.
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

// getCategoryName is a function that retrieves the category name of a category.
// It takes in one parameter: categoryID.
// It queries the database for the category name of the category with the provided categoryID.
// It returns the category name and nil if the category exists.
// It returns an empty string and an error if the category does not exist or if the query fails.
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

// getGenreID is a function that retrieves the genre ID of a genre.
// It takes in one parameter: genreName.
// It queries the database for the genre ID of the genre with the provided genreName.
// It returns the genre ID and nil if the genre exists.
// It returns 0 and an error if the genre does not exist or if the query fails.
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

// getGenreName is a function that retrieves the genre name of a genre.
// It takes in one parameter: genreID.
// It queries the database for the genre name of the genre with the provided genreID.
// It returns the genre name and nil if the genre exists.
// It returns an empty string and an error if the genre does not exist or if the query fails.
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

// VideoDesc is a struct that represents the description of a video.
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

// videoDescForTable is a function that retrieves the video descriptions for a user.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos uploaded by the user with the provided userID.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
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

// videosBrowser is a function that retrieves all video descriptions.
// It does not take any parameters.
// It queries the database for the video descriptions of all videos.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
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

// videoDescForEdit is a function that edits the details of a video.
// It takes in several parameters including videoID, title, description, category, and genre.
// It updates the video record in the database with the provided parameters.
// It returns an error if the update fails.
func (db *databaseConn) videoDescForEdit(videoID int, title string, description string, category int, genre int) error {
	// Update the video details in the database based on the videoID
	query := "UPDATE videos SET Title = ?, Description = ?, CategoryID = ?, GenreID = ? WHERE VideoID = ?"
	_, err := db.DB.Exec(query, title, description, category, genre, videoID)
	if err != nil {
		return err
	}
	return nil
}

// deleteVideo is a function that deletes a video.
// It takes in one parameter: videoID.
// It deletes the video record in the database with the provided videoID.
// It returns an error if the deletion fails.
func (db *databaseConn) deleteVideo(videoID int) error {
	query := "DELETE FROM videos WHERE VideoID = ?"
	_, err := db.DB.Exec(query, videoID)
	if err != nil {
		return err
	}
	return nil
}

// deleteVideoFromFile is a function that deletes a video and its related data.
// It takes in one parameter: videoID.
// It deletes the video record and its related data in the VideoActions and Comments tables in the database with the provided videoID.
// It returns the video URL, thumbnail URL, and nil if the deletion is successful.
// It returns empty strings and an error if the deletion fails.
func (db *databaseConn) deleteVideoFromFile(videoID int) (string, string, error) {
	// URL and ThumbnailURL of the video to be deleted
	var videoURL string
	var thumbnailURL string
	// Delete related data from the VideoActions table
	query := "DELETE FROM videoactions WHERE VideoID = ?"
	_, err := db.DB.Exec(query, videoID)
	if err != nil {
		return "", "", err
	}

	// Delete related data from the Comments table
	/*query = "DELETE FROM comments WHERE VideoID = ?"
	_, err = db.DB.Exec(query, videoID)
	if err != nil {
		return "", "", err
	}*/

	query = "SELECT URL, ThumbnailURL FROM videos WHERE VideoID = ?"
	err = db.DB.QueryRow(query, videoID).Scan(&videoURL, &thumbnailURL)
	if err != nil {
		return "", "", err
	}

	// Delete the video from the Videos table
	query = "DELETE FROM videos WHERE VideoID = ?"
	_, err = db.DB.Exec(query, videoID)
	if err != nil {
		return "", "", err
	}

	return videoURL, thumbnailURL, nil
}

// recentlyAddedVideos is a function that retrieves the recently added videos.
// It does not take any parameters.
// It queries the database for the video descriptions of the 10 most recently added videos.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) recentlyAddedVideos() ([]VideoDesc, error) {
	query := "SELECT * FROM videos ORDER BY UploadDate DESC LIMIT 15"
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

// recommendedVideos is a function that retrieves recommended videos.
// It does not take any parameters.
// It queries the database for the video description of a random video.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
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

// weeklyTop is a function that retrieves the top videos of the week.
// It does not take any parameters.
// It queries the database for the video descriptions of the 10 videos with the most views in the past week.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) weeklyTop() ([]VideoDesc, error) {
	query := "SELECT DISTINCT V.VideoID, V.Title, V.Description, V.URL, V.ThumbnailURL, V.UploaderID, V.UploadDate, V.ViewsCount, V.LikesCount, V.DislikesCount, V.Duration, V.CategoryID, V.GenreID FROM Videos V JOIN VideoActions VA ON V.VideoID = VA.VideoID WHERE VA.ActionsDate BETWEEN CURRENT_DATE - INTERVAL DAYOFWEEK(CURRENT_DATE) + 6 DAY AND CURRENT_DATE ORDER BY V.ViewsCount DESC LIMIT 10"
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

// videoActions is a function that records a user's action on a video.
// It takes in two parameters: videoID and userID.
// It checks if the user has already actioned on the video.
// If the user has already actioned on the video, it updates the ActionDate and ActionTime to the current date and time.
// If the user has not actioned on the video, it adds the video to the VideoActions table.
// It returns an error if the update or insertion fails.
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

func (db *databaseConn) updateViews(videoID int) error {
	// Prepare the SQL query to increment the ViewsCount field
	query := "UPDATE videos SET ViewsCount = ViewsCount + 1 WHERE VideoID = ?"

	// Execute the query
	_, err := db.DB.Exec(query, videoID)

	// If there is an error, return it
	if err != nil {
		return err
	}

	// If the query is executed successfully, return nil
	return nil
}

// continueWatching is a function that retrieves the videos that a user is currently watching.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the 10 most recently watched videos by the user with the provided userID.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
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

// caroselSlide is a function that retrieves random videos for the carousel slide.
// It does not take any parameters.
// It queries the database for the video descriptions of 10 random videos.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
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

// VideosearchDesc is a struct that represents the description of a video for search results.
type VideosearchDesc struct {
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
	Genre         string
	Status        string
}

// searchVideos is a function that searches for videos based on a search query.
// It takes in two parameters: searchQuery and userID.
// It queries the database for the video descriptions of the videos that match the search query.
// It returns a slice of VideosearchDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) searchVideos(searchQuery string, userID int) ([]VideosearchDesc, error) {
	// Define the status fields
	statusFields := []string{"Watching", "Completed", "On_hold", "Considering", "Dropped"}

	// If the searchQuery is empty, return all the videos in ascending order of their title
	if searchQuery == "" {
		query1 := "SELECT V.*, G.GenreName, VA.Watching, VA.Completed, VA.On_hold, VA.Considering, VA.Dropped FROM videos V JOIN genres G ON V.GenreID = G.GenreID LEFT JOIN videoactions VA ON V.VideoID = VA.VideoID AND VA.UserID = ? ORDER BY V.Title ASC LIMIT 18"
		rows, err := db.DB.Query(query1, userID)
		if err != nil {
			return nil, err
		}

		var videos []VideosearchDesc
		for rows.Next() {
			var video VideosearchDesc
			var statusValues [5]*int
			err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID, &video.Genre, &statusValues[0], &statusValues[1], &statusValues[2], &statusValues[3], &statusValues[4])
			if err != nil {
				return nil, err
			}

			// Set the status of the video
			for i, status := range statusValues {
				if status != nil && *status == 1 {
					video.Status = statusFields[i]
					break
				}
			}

			videos = append(videos, video)
		}
		return videos, nil
	}

	query := "SELECT V.*, G.GenreName, VA.Watching, VA.Completed, VA.On_hold, VA.Considering, VA.Dropped FROM videos V JOIN genres G ON V.GenreID = G.GenreID LEFT JOIN videoactions VA ON V.VideoID = VA.VideoID AND VA.UserID = ? WHERE V.Title LIKE ?"
	rows, err := db.DB.Query(query, userID, "%"+searchQuery+"%")
	if err != nil {
		return nil, err
	}

	var videos []VideosearchDesc
	for rows.Next() {
		var video VideosearchDesc
		var statusValues [5]*int
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID, &video.Genre, &statusValues[0], &statusValues[1], &statusValues[2], &statusValues[3], &statusValues[4])
		if err != nil {
			return nil, err
		}

		// Set the status of the video
		for i, status := range statusValues {
			if status != nil && *status == 1 {
				video.Status = statusFields[i]
				break
			}
		}

		videos = append(videos, video)
	}

	return videos, nil
}

// videoAction is a function that retrieves the status of a user's action on a video.
// It takes in two parameters: videoID and userID.
// It queries the database for the status of the user's action on the video with the provided videoID.
// It returns the status and nil if the query is successful.
// It returns an empty string and an error if the query fails.
func (db *databaseConn) videoAction(videoID int, userID int) (string, error) {
	query := "SELECT Watching, Completed, On_hold, Considering, Dropped FROM videoactions WHERE VideoID = ? AND UserID = ?"
	rows, err := db.DB.Query(query, videoID, userID)
	if err != nil {
		return "", err
	}

	var status string
	for rows.Next() {
		var watching, completed, onHold, considering, dropped int
		err := rows.Scan(&watching, &completed, &onHold, &considering, &dropped)
		if err != nil {
			return "", err
		}
		if watching == 1 {
			status = "Watching"
		} else if completed == 1 {
			status = "Completed"
		} else if onHold == 1 {
			status = "On-hold"
		} else if considering == 1 {
			status = "Considering"
		} else if dropped == 1 {
			status = "Dropped"
		}
	}
	// convert the status to lowercase
	status = strings.ToLower(status)
	return status, nil
}

// videoActionChanged is a function that changes the status of a user's action on a video.
// It takes in three parameters: videoID, userID, and status.
// It updates the status of the user's action on the video with the provided videoID to the provided status.
// It returns an error if the update fails.
func (db *databaseConn) videoActionChanged(videoID int, userID int, status string) error {
	// Convert the status to lowercase
	status = strings.ToLower(status)

	// Prepare the query
	query := "UPDATE videoactions SET Watching = ?, Completed = ?, On_hold = ?, Considering = ?, Dropped = ? WHERE VideoID = ? AND UserID = ?"

	// Initialize all status values to 0
	watching, completed, onHold, considering, dropped := 0, 0, 0, 0, 0

	// Depending on the status, set the corresponding value to 1
	switch status {
	case "watching":
		watching = 1
	case "completed":
		completed = 1
	case "on-hold":
		onHold = 1
	case "considering":
		considering = 1
	case "dropped":
		dropped = 1
	default:
		return errors.New("invalid status")
	}

	// Execute the query
	_, err := db.DB.Exec(query, watching, completed, onHold, considering, dropped, videoID, userID)
	if err != nil {
		return err
	}

	return nil
}

type videoActionInfo struct {
	Recommends  int
	Watching    int
	Completed   int
	OnHold      int
	Considering int
	Dropped     int
}

// userProfileVideosData is a function that retrieves the video action data for a user.
// It takes in one parameter: userID.
// It queries the database for the count of different video actions (like recommends, watching, completed, on-hold, considering, dropped) performed by the user with the provided userID.
// It returns a struct of type videoActionInfo containing the counts of different video actions and nil if the query is successful.
// It returns an empty videoActionInfo struct and an error if the query fails.
func (db *databaseConn) userProfileVideosData(userID int) (videoActionInfo, error) {
	// set value fo videoActionInfo to 0 by default
	videoActionInfo := videoActionInfo{
		Recommends:  0,
		Watching:    0,
		Completed:   0,
		OnHold:      0,
		Considering: 0,
		Dropped:     0,
	}

	query := "SELECT VideoID, Recommends, Watching, Completed, On_hold, Considering, Dropped FROM videoactions WHERE UserID = ?"
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return videoActionInfo, err
	}
	defer rows.Close()

	var recommendsCount, watchingCount, completedCount, onHoldCount, consideringCount, droppedCount int
	for rows.Next() {
		var videoID, recommends, watching, completed, onHold, considering, dropped int
		err := rows.Scan(&videoID, &recommends, &watching, &completed, &onHold, &considering, &dropped)
		if err != nil {
			return videoActionInfo, err
		}
		if recommends == 1 {
			recommendsCount++
		}
		if watching == 1 {
			watchingCount++
		} else if completed == 1 {
			completedCount++
		} else if onHold == 1 {
			onHoldCount++
		} else if considering == 1 {
			consideringCount++
		} else if dropped == 1 {
			droppedCount++
		}
	}

	videoActionInfo.Recommends = recommendsCount
	videoActionInfo.Watching = watchingCount
	videoActionInfo.Completed = completedCount
	videoActionInfo.OnHold = onHoldCount
	videoActionInfo.Considering = consideringCount
	videoActionInfo.Dropped = droppedCount

	return videoActionInfo, nil
}

// watchingVideos is a function that retrieves the videos that a user is currently watching.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos that the user with the provided userID is currently watching.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) watchingVideos(userID int) ([]VideoDesc, error) {
	// Query the videoactions table for videos that the user is currently watching
	rows, err := db.DB.Query("SELECT VideoID FROM videoactions WHERE UserID = ? AND Watching = 1 ORDER BY ActionsDate DESC, ActionTime DESC LIMIT 8", userID)
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

	// Query the videos table for the details of the videos
	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s)", videoIDsStr))
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

// onHoldVideos is a function that retrieves the videos that a user has put on hold.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos that the user with the provided userID has put on hold.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) onHoldVideos(userID int) ([]VideoDesc, error) {
	rows, err := db.DB.Query("SELECT VideoID FROM videoactions WHERE UserID = ? AND On_hold = 1 ORDER BY ActionsDate DESC, ActionTime DESC LIMIT 8", userID)
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

	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s)", videoIDsStr))
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

// consideringVideos is a function that retrieves the videos that a user is considering to watch.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos that the user with the provided userID is considering to watch.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) consideringVideos(userID int) ([]VideoDesc, error) {
	// Query the videoactions table for videos that the user is considering to watch
	rows, err := db.DB.Query("SELECT VideoID FROM videoactions WHERE UserID = ? AND Considering = 1 ORDER BY ActionsDate DESC, ActionTime DESC LIMIT 8", userID)
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

	// Query the videos table for the details of the videos
	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s)", videoIDsStr))
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

type recentlyCompletedVideoStruct struct {
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
	CompletedDate string
}

// recentlyCompletedVideos is a function that retrieves the videos that a user has recently completed watching.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos that the user with the provided userID has recently marked as completed.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) recentlyCompletedVideos(userID int) ([]recentlyCompletedVideoStruct, error) {
	// Query the videoactions table for videos that the user has recently completed
	rows, err := db.DB.Query("SELECT VideoID, ActionsDate FROM videoactions WHERE UserID = ? AND Completed = 1 ORDER BY ActionsDate DESC, ActionTime DESC LIMIT 10", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videoIDs []int
	var completedDates []string
	for rows.Next() {
		var videoID int
		var completedDate string
		if err := rows.Scan(&videoID, &completedDate); err != nil {
			return nil, err
		}
		videoIDs = append(videoIDs, videoID)
		completedDates = append(completedDates, completedDate)
	}

	if len(videoIDs) == 0 {
		return []recentlyCompletedVideoStruct{}, nil
	}

	videoIDStrs := make([]string, len(videoIDs))
	for i, videoID := range videoIDs {
		videoIDStrs[i] = strconv.Itoa(videoID)
	}
	videoIDsStr := strings.Join(videoIDStrs, ",")

	// Query the videos table for the details of the videos
	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s)", videoIDsStr))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []recentlyCompletedVideoStruct
	for rows.Next() {
		var video recentlyCompletedVideoStruct
		if err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID); err != nil {
			return nil, err
		}
		for i, videoID := range videoIDs {
			if video.VideoID == videoID {
				video.CompletedDate = completedDates[i]
				break
			}
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// userDetails is a function that retrieves the details of a user.
// It takes in one parameter: userID.
// It queries the database for the username, email, and admin status of the user with the provided userID.
// It returns the username, email, admin status, and nil if the user exists.
// It returns empty strings and an error if the user does not exist or if the query fails.
func (db *databaseConn) userDetails(userID int) (string, string, string, error) {
	var userName, email, isAdmin string
	query := "SELECT UserName, Email, Admin FROM users WHERE UserID = ?"
	err := db.DB.QueryRow(query, userID).Scan(&userName, &email, &isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", "", sql.ErrNoRows
		}
		return "", "", "", err
	}
	return userName, email, isAdmin, nil
}

// recommendedVideoListDatabase is a function that retrieves the videos that a user has recommended.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos that the user with the provided userID has recommended.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) recommendedVideoListDatabase(userID int) ([]VideoDesc, error) {
	// Query the videoactions table for videos that the user has recommended
	rows, err := db.DB.Query("SELECT VideoID FROM videoactions WHERE UserID = ? AND Recommends = 1", userID)
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

	// Query the videos table for the details of the videos
	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s)", videoIDsStr))
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

// watchingVideoListDatabase is a function that retrieves the videos that a user is currently watching.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos that the user with the provided userID is currently watching.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) watchingVideoListDatabase(userID int) ([]VideoDesc, error) {
	// Query the videoactions table for videos that the user is currently watching
	rows, err := db.DB.Query("SELECT VideoID FROM videoactions WHERE UserID = ? AND Watching = 1", userID)
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

	// Query the videos table for the details of the videos
	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s)", videoIDsStr))
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

// completedVideoListDatabase is a function that retrieves the videos that a user has completed watching.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos that the user with the provided userID has marked as completed.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) completedVideoListDatabase(userID int) ([]VideoDesc, error) {
	// Query the videoactions table for videos that the user has completed
	rows, err := db.DB.Query("SELECT VideoID FROM videoactions WHERE UserID = ? AND Completed = 1", userID)
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

	// Query the videos table for the details of the videos
	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s)", videoIDsStr))
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

// onHoldVideoListDatabase is a function that retrieves the videos that a user has put on hold.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos that the user with the provided userID has put on hold.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) onHoldVideoListDatabase(userID int) ([]VideoDesc, error) {
	// Query the videoactions table for videos that the user has put on hold
	rows, err := db.DB.Query("SELECT VideoID FROM videoactions WHERE UserID = ? AND On_hold = 1", userID)
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

	// Query the videos table for the details of the videos
	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s)", videoIDsStr))
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

// consideringVideoListDatabase is a function that retrieves the videos that a user is considering to watch.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos that the user with the provided userID is considering to watch.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) consideringVideoListDatabase(userID int) ([]VideoDesc, error) {
	// Query the videoactions table for videos that the user is considering
	rows, err := db.DB.Query("SELECT VideoID FROM videoactions WHERE UserID = ? AND Considering = 1", userID)
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

	// Query the videos table for the details of the videos
	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s)", videoIDsStr))
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

// droppedVideoListDatabase is a function that retrieves the videos that a user has dropped.
// It takes in one parameter: userID.
// It queries the database for the video descriptions of the videos that the user with the provided userID has dropped.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// It returns nil and an error if the query fails.
func (db *databaseConn) droppedVideoListDatabase(userID int) ([]VideoDesc, error) {
	// Query the videoactions table for videos that the user has dropped
	rows, err := db.DB.Query("SELECT VideoID FROM videoactions WHERE UserID = ? AND Dropped = 1", userID)
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

	// Query the videos table for the details of the videos
	rows, err = db.DB.Query(fmt.Sprintf("SELECT * FROM videos WHERE VideoID IN (%s)", videoIDsStr))
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

// commentOnVideo is a function that allows a user to comment on a video.
// It takes in three parameters: userID, videoID, and comment.
// It inserts a new record into the comments table in the database with the provided userID, videoID, and comment.
// It returns an error if the insertion fails.
func (db *databaseConn) commentOnVideo(userID int, videoID int, comment string) error {

	// Check if the user has already commented on the video
	query := "SELECT * FROM comments WHERE UserID = ? AND VideoID = ?"
	rows, err := db.DB.Query(query, userID, videoID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// If the user has already commented on the video, update the comment
	if rows.Next() {
		// Prepare the query to update the comment
		query = "UPDATE comments SET CommentText = ?, CommentDate = ? WHERE UserID = ? AND VideoID = ?"
		_, err := db.DB.Exec(query, comment, time.Now().Format("2006-01-02"), userID, videoID)
		if err != nil {
			return err
		}
		return nil
	}

	// Prepare the query to insert the comment into the comments table
	query = "INSERT INTO comments (UserID, VideoID, CommentText, CommentDate) VALUES (?, ?, ?, ?)"

	// Execute the query
	_, err = db.DB.Exec(query, userID, videoID, comment, time.Now().Format("2006-01-02"))

	// If there is an error, return it
	if err != nil {
		return err
	}

	// If the query is executed successfully, return nil
	return nil
}

type Comment struct {
	CommentID   int
	UserID      int
	VideoID     int
	UserName    string
	CommentText string
	CommentDate string
	Upvotes     int
	Downvotes   int
	ImageURL    string
}

// getComments is a function that retrieves the comments for a video.
// It takes in one parameter: videoID.
// It queries the database for the comments on the video with the provided videoID.
// It returns a slice of Comment structs and nil if the query is successful.
// Each Comment struct contains the comment ID, user ID, video ID, username, comment text, comment date, upvotes, downvotes, and image URL.
// It returns nil and an error if the query fails.
func (db *databaseConn) getComments(videoID int) ([]Comment, error) {
	query := `SELECT C.CommentID, C.UserID, C.VideoID, U.UserName, C.CommentText, C.CommentDate, 
				COALESCE(SUM(CA.Upvotes), 0) AS Upvotes, COALESCE(SUM(CA.Downvotes), 0) AS Downvotes, 
				MAX(I.ImageURL) AS ImageURL
				FROM comments C 
				JOIN users U ON C.UserID = U.UserID 
				LEFT JOIN CommentActions CA ON C.CommentID = CA.CommentID 
				LEFT JOIN userprofileimages I ON C.UserID = I.UserID 
				WHERE C.VideoID = ? 
				GROUP BY C.CommentID 
				ORDER BY C.CommentDate DESC;`
	rows, err := db.DB.Query(query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.CommentID, &comment.UserID, &comment.VideoID, &comment.UserName, &comment.CommentText, &comment.CommentDate, &comment.Upvotes, &comment.Downvotes, &comment.ImageURL); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// upvoteComment is a function that allows a user to upvote a comment.
// It takes in two parameters: commentID and userID.
// It updates the upvotes count for the comment with the provided commentID in the database.
// It returns an error if the update fails.
func (db *databaseConn) upvoteComment(commentID int, userID int) error {
	// Check if the comment exists
	var exists int
	/*err := db.DB.QueryRow("SELECT COUNT(*) FROM comments WHERE CommentID = ?", commentID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("comment does not exist")
	}*/

	// Check if the user exists
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE UserID = ?", userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("user does not exist")
	}

	// Query the CommentActions table to check if the user has already upvoted or downvoted the comment
	query := "SELECT Upvotes, Downvotes FROM CommentActions WHERE CommentID = ? AND UserID = ?"
	row := db.DB.QueryRow(query, commentID, userID)

	var upvotes, downvotes int
	err = row.Scan(&upvotes, &downvotes)

	// If the user has already upvoted the comment, return an error message
	/*if err == nil && upvotes == 1 {
		return errors.New("Comment is already upvoted")
	}*/

	// If the user has not performed any action before, insert a new upvote
	if errors.Is(err, sql.ErrNoRows) {
		query = "INSERT INTO CommentActions (CommentID, UserID, Upvotes, Downvotes) VALUES (?, ?, 1, 0)"
		_, err = db.DB.Exec(query, commentID, userID)
	} else { // If the user has previously downvoted the comment, update the record
		query = "UPDATE CommentActions SET Upvotes = 1, Downvotes = 0 WHERE CommentID = ? AND UserID = ?"
		_, err = db.DB.Exec(query, commentID, userID)
	}
	return err
}

// reverseUpvoteComment is a function that allows a user to reverse their upvote on a comment.
// It takes in two parameters: commentID and userID.
// It updates the upvotes count for the comment with the provided commentID in the database, effectively reversing the upvote.
// It returns an error if the update fails.
func (db *databaseConn) reverseUpvoteComment(commentID int, userID int) error {
	// Check if the comment exists
	//var exists int

	// Check if the user exists
	/*err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE UserID = ?", userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("User does not exist")
	}*/

	// Query the CommentActions table to check if the user has already upvoted or downvoted the comment
	query := "SELECT Upvotes, Downvotes FROM CommentActions WHERE CommentID = ? AND UserID = ?"
	row := db.DB.QueryRow(query, commentID, userID)

	var upvotes, downvotes int
	err := row.Scan(&upvotes, &downvotes)

	query = "UPDATE CommentActions SET Upvotes = 0, Downvotes = 0 WHERE CommentID = ? AND UserID = ?"
	_, err = db.DB.Exec(query, commentID, userID)
	return err
}

// downvoteComment is a function that allows a user to downvote a comment.
// It takes in two parameters: commentID and userID.
// It updates the downvotes count for the comment with the provided commentID in the database.
// It returns an error if the update fails.
func (db *databaseConn) downvoteComment(commentID int, userID int) error {
	// Check if the comment exists
	var exists int
	/*err := db.DB.QueryRow("SELECT COUNT(*) FROM comments WHERE CommentID = ?", commentID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("Comment does not exist")
	}*/

	// Check if the user exists
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE UserID = ?", userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("User does not exist")
	}

	// Query the CommentActions table to check if the user has already upvoted or downvoted the comment
	query := "SELECT Upvotes, Downvotes FROM CommentActions WHERE CommentID = ? AND UserID = ?"
	row := db.DB.QueryRow(query, commentID, userID)

	var upvotes, downvotes int
	err = row.Scan(&upvotes, &downvotes)

	// If the user has already downvoted the comment, return an error message
	/*if err == nil && downvotes == 1 {
		return errors.New("Comment is already downvoted")
	}*/

	// If the user has not performed any action before, insert a new downvote
	if err == sql.ErrNoRows {
		query = "INSERT INTO CommentActions (CommentID, UserID, Upvotes, Downvotes) VALUES (?, ?, 0, 1)"
		_, err = db.DB.Exec(query, commentID, userID)
	} else { // If the user has previously upvoted the comment, update the record
		query = "UPDATE CommentActions SET Upvotes = 0, Downvotes = 1 WHERE CommentID = ? AND UserID = ?"
		_, err = db.DB.Exec(query, commentID, userID)
	}

	return err
}

// reverseDownvoteComment is a function that allows a user to reverse their downvote on a comment.
// It takes in two parameters: commentID and userID.
// It updates the downvotes count for the comment with the provided commentID in the database, effectively reversing the downvote.
// It returns an error if the update fails.
func (db *databaseConn) reverseDownvoteComment(commentID int, userID int) error {
	// Check if the comment exists
	var exists int
	/*err := db.DB.QueryRow("SELECT COUNT(*) FROM comments WHERE CommentID = ?", commentID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("Comment does not exist")
	}*/

	// Check if the user exists
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE UserID = ?", userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("User does not exist")
	}

	// Query the CommentActions table to check if the user has already upvoted or downvoted the comment
	query := "SELECT Upvotes, Downvotes FROM CommentActions WHERE CommentID = ? AND UserID = ?"
	row := db.DB.QueryRow(query, commentID, userID)

	var upvotes, downvotes int
	err = row.Scan(&upvotes, &downvotes)

	query = "UPDATE CommentActions SET Upvotes = 0, Downvotes = 0 WHERE CommentID = ? AND UserID = ?"
	_, err = db.DB.Exec(query, commentID, userID)

	return err
}

type commentDetails struct {
	Upvote    int
	Downvote  int
	CommentID int
	VideoID   int
}

// commentDetails is a function that retrieves the details of a comment.
// It takes in two parameters: videoID and userID.
// It queries the database for the upvotes, downvotes, commentID, and videoID of the comment on the video with the provided videoID made by the user with the provided userID.
// It returns a slice of commentDetails structs and nil if the query is successful.
// Each commentDetails struct contains the upvotes, downvotes, commentID, and videoID.
// It returns nil and an error if the query fails.
func (db *databaseConn) commentDetails(videoID int, userID int) ([]commentDetails, error) {
	// Prepare the SQL query
	query := `SELECT c.CommentID, c.VideoID, ca.Upvotes, ca.Downvotes
				FROM comments c
				JOIN commentactions ca ON c.CommentID = ca.CommentID
				WHERE c.VideoID = ?
				AND ca.UserID = ?`

	// Execute the query
	rows, err := db.DB.Query(query, videoID, userID)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold the comment details
	var details []commentDetails

	// Iterate over the rows in the result set
	for rows.Next() {
		var detail commentDetails

		// Scan the values into a commentDetails struct
		if err := rows.Scan(&detail.CommentID, &detail.VideoID, &detail.Upvote, &detail.Downvote); err != nil {
			return nil, err
		}
		// Append the struct to the slice
		details = append(details, detail)
	}

	// Return the slice and nil (no error)
	return details, nil
}

// likeVideo is a function that allows a user to like a video.
// It takes in two parameters: videoID and userID.
// It updates the likes count for the video with the provided videoID in the database.
// It returns an error if the update fails.
func (db *databaseConn) likeVideo(videoID int, userID int) error {

	queryToCheck := "SELECT Recommends FROM videoactions WHERE VideoID = ? AND UserID = ?"

	row := db.DB.QueryRow(queryToCheck, videoID, userID)

	var recommends int
	err := row.Scan(&recommends)
	if err != nil {
		return err
	}
	if recommends == -1 {
		query := "UPDATE videos SET DislikesCount = DislikesCount - 1 WHERE VideoID = ?"
		_, err := db.DB.Exec(query, videoID)
		if err != nil {
			return err
		}
	}

	// Prepare the SQL query to increment the LikesCount field
	query := "UPDATE videos SET LikesCount = LikesCount + 1 WHERE VideoID = ?"

	// Execute the query
	_, err = db.DB.Exec(query, videoID)

	// If there is an error, return it
	if err != nil {
		return err
	}

	// Prepare the SQL query to set the Recommends field to 1
	query = "UPDATE videoactions SET Recommends = 1 WHERE VideoID = ? AND UserID = ?"

	// Execute the query
	_, err = db.DB.Exec(query, videoID, userID)

	// If there is an error, return it
	if err != nil {
		return err
	}

	// If the query is executed successfully, return nil
	return nil
}

// reverseLikeVideo is a function that allows a user to reverse their like on a video.
// It takes in two parameters: videoID and userID.
// It updates the likes count for the video with the provided videoID in the database, effectively reversing the like.
// It returns an error if the update fails.
func (db *databaseConn) reverseLikeVideo(videoID int, userID int) error {
	// Prepare the SQL query to decrement the LikesCount field
	query := "UPDATE videos SET LikesCount = LikesCount - 1 WHERE VideoID = ?"

	// Execute the query
	_, err := db.DB.Exec(query, videoID)

	// If there is an error, return it
	if err != nil {
		return err
	}

	// Prepare the SQL query to set the Recommends field to 0
	query = "UPDATE videoactions SET Recommends = 0 WHERE VideoID = ? AND UserID = ?"

	// Execute the query
	_, err = db.DB.Exec(query, videoID, userID)

	// If there is an error, return it
	if err != nil {
		return err
	}

	// If the query is executed successfully, return nil
	return nil
}

// dislikeVideo is a function that allows a user to dislike a video.
// It takes in two parameters: videoID and userID.
// It updates the dislikes count for the video with the provided videoID in the database.
// It returns an error if the update fails.
func (db *databaseConn) dislikeVideo(videoID int, userID int) error {

	queryToCheck := "SELECT Recommends FROM videoactions WHERE VideoID = ? AND UserID = ?"

	row := db.DB.QueryRow(queryToCheck, videoID, userID)

	var recommends int
	err := row.Scan(&recommends)
	if err != nil {
		return err
	}

	if recommends == 1 {
		query := "UPDATE videos SET LikesCount = LikesCount - 1 WHERE VideoID = ?"
		_, err := db.DB.Exec(query, videoID)
		if err != nil {
			return err
		}
	}

	// Prepare the SQL query to increment the DislikesCount field
	query := "UPDATE videos SET DislikesCount = DislikesCount + 1 WHERE VideoID = ?"

	// Execute the query
	_, err = db.DB.Exec(query, videoID)

	// If there is an error, return it
	if err != nil {
		return err
	}

	// Prepare the SQL query to set the Recommends field to 0
	query = "UPDATE videoactions SET Recommends = -1 WHERE VideoID = ? AND UserID = ?"

	// Execute the query
	_, err = db.DB.Exec(query, videoID, userID)

	// If there is an error, return it
	if err != nil {
		return err
	}

	// If the query is executed successfully, return nil
	return nil
}

// reverseDislikeVideo is a function that allows a user to reverse their dislike on a video.
// It takes in two parameters: videoID and userID.
// It updates the dislikes count for the video with the provided videoID in the database, effectively reversing the dislike.
// It returns an error if the update fails.
func (db *databaseConn) reverseDislikeVideo(videoID int, userID int) error {
	// Prepare the SQL query to decrement the DislikesCount field
	query := "UPDATE videos SET DislikesCount = DislikesCount - 1 WHERE VideoID = ?"

	// Execute the query
	_, err := db.DB.Exec(query, videoID)

	// If there is an error, return it
	if err != nil {
		return err
	}

	query = "UPDATE videoactions SET Recommends = 0 WHERE VideoID = ? AND UserID = ?"

	// Execute the query
	_, err = db.DB.Exec(query, videoID, userID)

	// If there is an error, return it
	if err != nil {
		return err
	}

	// If the query is executed successfully, return nil
	return nil
}

// isLikedDisliked is a function that checks if a user has liked or disliked a video.
// It takes in two parameters: videoID and userID.
// It queries the database for the 'Recommends' status of the video with the provided videoID for the user with the provided userID.
// It returns the 'Recommends' status, which is 1 if the user has liked the video, -1 if the user has disliked the video, and 0 if the user has not actioned on the video.
// It also returns an error if the query fails.
func (db *databaseConn) isLikedDisliked(videoID int, userID int) (int, int, error) {
	// Prepare the SQL query to check if the user has liked or disliked the video
	query := "SELECT Recommends FROM videoactions WHERE VideoID = ? AND UserID = ?"

	// Execute the query
	row := db.DB.QueryRow(query, videoID, userID)

	// Initialize a variable to hold the result
	var recommends int

	// Scan the result into the variable
	err := row.Scan(&recommends)

	// If there is an error, return it
	if err != nil {
		return 0, 0, err
	}

	// If the value is 1, the user has liked the video
	if recommends == 1 {
		return 1, 0, nil
	} else if recommends == -1 {
		return 0, 1, nil
	} else {
		return 0, 0, nil
	}
}

// likeDislikeCount is a function that retrieves the count of likes and dislikes for a video.
// It takes in one parameter: videoID.
// It queries the database for the count of 'Recommends' status of the video with the provided videoID.
// It returns the count of likes and dislikes, and nil if the query is successful.
// It returns 0 for both counts and an error if the query fails.
func (db *databaseConn) likeDislikeCount(videoID int) (int, int, error) {
	// Prepare the SQL query to get the LikesCount and DislikesCount fields
	query := "SELECT LikesCount, DislikesCount FROM videos WHERE VideoID = ?"

	// Execute the query
	row := db.DB.QueryRow(query, videoID)

	// Initialize variables to hold the results
	var likesCount, dislikesCount int

	// Scan the results into the variables
	err := row.Scan(&likesCount, &dislikesCount)

	// If there is an error, return it
	if err != nil {
		return 0, 0, err
	}

	// Return the results
	return likesCount, dislikesCount, nil
}

// autoComplete is a function that retrieves a list of video descriptions for auto-completion purposes.
// It does not take any parameters.
// It queries the database for the video descriptions of all videos.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// Each VideoDesc struct contains the details of a video, including the video ID, title, description, uploader ID, upload date, duration, category ID, genre, likes count, dislikes count, views count, and video URL.
// It returns nil and an error if the query fails.
func (db *databaseConn) autoComplete() ([]VideoDesc, error) {
	// Prepare the SQL query
	query := "SELECT * FROM videos"

	// Execute the query
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the results
	var videos []VideoDesc

	// Iterate over the rows in the result set
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	// Return the results
	return videos, nil
}

// editUserProfile is a function that allows a user to edit their profile.
// It takes in three parameters: userName, email, and prevEmail.
// It first checks if a user with the provided prevEmail exists in the database.
// If the user exists, it updates the user's userName and email in the database to the provided userName and email.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) editUserProfile(userName string, email string, prevEmail string) error {
	query := "UPDATE users SET UserName = ?, Email = ? WHERE Email = ?"
	_, err := db.DB.Exec(query, userName, email, prevEmail)
	if err != nil {
		return err
	}
	return nil
}

// ChangePassword is a function that allows a user to change their password.
// It takes in three parameters: oldPassword, newPassword, and email.
// It first checks if a user with the provided email exists in the database.
// If the user exists, it verifies the oldPassword with the password in the database.
// If the oldPassword is correct, it updates the user's password in the database to the provided newPassword.
// It returns an error if the user does not exist, if the oldPassword is incorrect, or if the update fails.
func (db *databaseConn) changePassword(oldPassword string, newPassword string, email string) error {
	query := "SELECT Password FROM users WHERE Email = ?"
	row := db.DB.QueryRow(query, email)
	var password string
	err := row.Scan(&password)
	if err != nil {
		return err
	}
	res := bcrypt.CompareHashAndPassword([]byte(password), []byte(oldPassword))
	if res != nil {
		return errors.New("old password is incorrect")

	}
	//if password != oldPassword {
	//	return errors.New("old password is incorrect")
	//}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	query = "UPDATE users SET Password = ? WHERE Email = ?"
	_, err = db.DB.Exec(query, encryptedPassword, email)
	if err != nil {
		return err
	}
	return nil
}

// resetPassword is a function that allows a user to reset their password.
// It takes in two parameters: email and password.
// It first checks if a user with the provided email exists in the database.
// If the user exists, it updates the user's password in the database to the provided password.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) resetPassword(email string, password string) error {
	fmt.Println("Reset password is triggered with email: ", email, "and password: ", password)
	query := "SELECT Email FROM users WHERE Email = ?"
	row := db.DB.QueryRow(query, email)
	var userEmail string
	err := row.Scan(&userEmail)
	if err != nil {
		return err
	}
	if userEmail != email {
		return errors.New("email does not exist")
	}

	query = "UPDATE users SET Password = ? WHERE Email = ?"
	_, err = db.DB.Exec(query, password, email)
	if err != nil {
		return err
	}

	return nil
}

// mostViewedVideos is a function that retrieves the most viewed videos.
// It does not take any parameters.
// It queries the database for the video descriptions of the videos sorted by views count in descending order.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// Each VideoDesc struct contains the details of a video, including the video ID, title, description, uploader ID, upload date, duration, category ID, genre, likes count, dislikes count, views count, and video URL.
// It returns nil and an error if the query fails.
func (db *databaseConn) mostViewedVideos() ([]VideoDesc, error) {
	// Prepare the SQL query
	query := "SELECT * FROM videos ORDER BY ViewsCount DESC LIMIT 10"

	// Execute the query
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the results
	var videos []VideoDesc

	// Iterate over the rows in the result set
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	// Return the results
	return videos, nil
}

// likeVsDislike is a function that retrieves the videos with the most likes and dislikes.
// It does not take any parameters.
// It queries the database for the video descriptions of the videos sorted by likes count and dislikes count in descending order.
// It returns two slices of VideoDesc structs and nil if the query is successful.
// The first slice contains the videos with the most likes, and the second slice contains the videos with the most dislikes.
// Each VideoDesc struct contains the details of a video, including the video ID, title, description, uploader ID, upload date, duration, category ID, genre, likes count, dislikes count, views count, and video URL.
// It returns nil for both slices and an error if the query fails.
func (db *databaseConn) likeVsDislike() ([]VideoDesc, []VideoDesc, error) {
	// Initialize a map to keep track of already selected video IDs
	selectedVideoIDs := make(map[int]bool)

	// Query for top 5 most liked videos
	query := "SELECT * FROM videos ORDER BY LikesCount DESC LIMIT 5"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var likedVideos []VideoDesc
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, nil, err
		}

		// Check if the video ID is already in the map
		if _, exists := selectedVideoIDs[video.VideoID]; !exists {
			// If not, add the video to the likedVideos slice and add its ID to the map
			likedVideos = append(likedVideos, video)
			selectedVideoIDs[video.VideoID] = true
		}
	}

	// If the count of liked videos is less than 5, handle this case
	if len(likedVideos) < 5 {
		// Handle this case as per your requirements
	}

	// Query for top 5 most disliked videos
	query = "SELECT * FROM videos ORDER BY DislikesCount DESC LIMIT 5"
	rows, err = db.DB.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var dislikedVideos []VideoDesc
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, nil, err
		}

		// Check if the video ID is already in the map
		if _, exists := selectedVideoIDs[video.VideoID]; !exists {
			// If not, add the video to the dislikedVideos slice and add its ID to the map
			dislikedVideos = append(dislikedVideos, video)
			selectedVideoIDs[video.VideoID] = true
		}
	}

	// If the count of disliked videos is less than 5, handle this case
	if len(dislikedVideos) < 5 {
		// Handle this case as per your requirements
	}

	return likedVideos, dislikedVideos, nil
}

// duration is a function that retrieves the duration of videos.
// It does not take any parameters.
// It queries the database for the duration of all videos.
// It returns a slice of VideoDesc structs and nil if the query is successful.
// Each VideoDesc struct contains the details of a video, including the video ID, title, description, uploader ID, upload date, duration, category ID, genre, likes count, dislikes count, views count, and video URL.
// It returns nil and an error if the query fails.
func (db *databaseConn) duration() ([]VideoDesc, error) {
	// Prepare the SQL query
	query := "SELECT * FROM videos ORDER BY ViewsCount"

	// Execute the query
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the results
	var videos []VideoDesc

	// Iterate over the rows in the result set
	for rows.Next() {
		var video VideoDesc
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	// Return the results
	return videos, nil
}

// deviceInfo is a function that retrieves the information about a device.
// It does not take any parameters.
// It queries the system or database for the details of the device, including the device ID, model, manufacturer, OS version, and other relevant details.
// It returns a struct containing the device details and nil if the query is successful.
// Each device detail struct contains the details of a device, including the device ID, model, manufacturer, OS version, and other relevant details.
// It returns nil and an error if the query fails.
func (db *databaseConn) deviceInfo(IP string, deviceType string, deviceOS string, browser string, loginTime string, countryCode string, countryName string, regionName string, cityName string, latitude float64, longitude float64, zipCode string, timeZone string, asn string, as_ string, isProxy bool) error {
	// Prepare the SQL query to select a record
	query := `SELECT LastLogin FROM ServerLogs WHERE IP = ? AND device_type = ? AND device_os = ? AND Browser = ? AND country_code = ? AND country_name = ? AND region_name = ? AND city_name = ?`
	row := db.DB.QueryRow(query, IP, deviceType, deviceOS, browser, countryCode, countryName, regionName, cityName)

	// Check if a record exists
	var lastLoginData string
	err := row.Scan(&lastLoginData)
	if err != nil && err != sql.ErrNoRows {
		// If an error occurred other than no rows found, return it
		return err
	}

	var lastLogins []string

	if lastLoginData != "" {
		// Parse the existing JSON data into a slice of strings
		err := json.Unmarshal([]byte(lastLoginData), &lastLogins)
		if err != nil {
			// If unable to unmarshal as []string, try to unmarshal as map[string]interface{}
			var lastLoginMap map[string]interface{}
			err := json.Unmarshal([]byte(lastLoginData), &lastLoginMap)
			if err != nil {
				return err
			}

			// Convert map values to []string
			for _, v := range lastLoginMap {
				if s, ok := v.(string); ok {
					lastLogins = append(lastLogins, s)
				}
			}
		}
	}

	// Append the new login time to the slice
	lastLogins = append(lastLogins, loginTime)

	// Convert the updated slice back into JSON
	updatedLastLoginData, err := json.Marshal(lastLogins)
	if err != nil {
		return err
	}

	if lastLoginData == "" {
		// If no record exists, insert a new one
		query = `INSERT INTO ServerLogs (IP, device_type, device_os, Browser, LastLogin, country_code, country_name, region_name, city_name, latitude, longitude, zip_code, time_zone, asn, as_, is_proxy) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = db.DB.Exec(query, IP, deviceType, deviceOS, browser, string(updatedLastLoginData), countryCode, countryName, regionName, cityName, latitude, longitude, zipCode, timeZone, asn, as_, isProxy)
	} else {
		// Update the LastLogin field in the database with the new JSON data
		query = `UPDATE ServerLogs SET LastLogin = ? WHERE IP = ? AND device_type = ? AND device_os = ? AND Browser = ? AND country_code = ? AND country_name = ? AND region_name = ? AND city_name = ?`
		_, err = db.DB.Exec(query, string(updatedLastLoginData), IP, deviceType, deviceOS, browser, countryCode, countryName, regionName, cityName)
	}

	if err != nil {
		return err
	}

	return nil
}

type ServerLog struct {
	IP          string
	DeviceType  string
	DeviceOS    string
	Browser     string
	LoginTime   string
	CountryCode string
	CountryName string
	RegionName  string
	CityName    string
	Latitude    float64
	Longitude   float64
	ZipCode     string
	TimeZone    string
	ASN         string
	AS          string
	IsProxy     bool
}

// serverLog is a function that retrieves the server logs.
// It does not take any parameters.
// It queries the database for all server logs.
// Each server log contains details such as IP, device type, device OS, browser, login time, country code, country name, region name, city name, latitude, longitude, zip code, time zone, ASN, AS, and whether it is a proxy.
// It returns a slice of ServerLog structs and nil if the query is successful.
// Each ServerLog struct contains the details of a server log.
// It returns nil and an error if the query fails.
func (db *databaseConn) serverLog() ([]ServerLog, error) {
	// Prepare the SQL query
	//query := "SELECT * FROM serverlogs"

	// select everything except ID
	query := "SELECT IP, device_type, device_os, Browser, LastLogin, country_code, country_name, region_name, city_name, latitude, longitude, zip_code, time_zone, asn, as_, is_proxy FROM ServerLogs"
	// Execute the query
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the results
	var logs []ServerLog

	// Iterate over the rows in the result set
	for rows.Next() {
		var log ServerLog
		err := rows.Scan(&log.IP, &log.DeviceType, &log.DeviceOS, &log.Browser, &log.LoginTime, &log.CountryCode, &log.CountryName, &log.RegionName, &log.CityName, &log.Latitude, &log.Longitude, &log.ZipCode, &log.TimeZone, &log.ASN, &log.AS, &log.IsProxy)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	// Return the results
	return logs, nil
}

// locationAnalysis is a function that performs an analysis of user locations.
// It does not take any parameters.
// It queries the database for all server logs and groups them by country name and country code.
// It counts the number of occurrences for each country name and country code.
// It also counts the total number of server logs.
// It returns four values:
// 1. A map where the keys are country names and the values are the counts of each country name.
// 2. A map where the keys are country codes and the values are the counts of each country code.
// 3. The total count of server logs.
// 4. An error if the query fails.
func (db *databaseConn) locationAnalysis() (map[string]int, map[string]int, int, error) {
	// Prepare the SQL query
	query := "SELECT country_name, country_code FROM ServerLogs"

	// Execute the query
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()

	// Initialize maps to hold the results
	countryNameCount := make(map[string]int)
	countryCodeCount := make(map[string]int)

	// Initialize a variable to hold the count
	count := 0

	// Iterate over the rows in the result set
	for rows.Next() {
		var countryName, countryCode string
		if err := rows.Scan(&countryName, &countryCode); err != nil {
			return nil, nil, 0, err
		}
		countryNameCount[countryName]++
		countryCodeCount[countryCode]++
		count++
	}

	// Return the results
	return countryNameCount, countryCodeCount, count, nil
}

type VideoDescEdit struct {
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
	CategoryName  string
	GenreName     string
}

// allVideoList is a function that retrieves a list of all videos.
// It does not take any parameters.
// It queries the database for the details of all videos, including the video ID, title, description, URL, thumbnail URL, uploader ID, upload date, views count, likes count, dislikes count, duration, category ID, genre ID, category name, and genre name.
// It returns a slice of VideoDescEdit structs and nil if the query is successful.
// Each VideoDescEdit struct contains the details of a video, including the video ID, title, description, URL, thumbnail URL, uploader ID, upload date, views count, likes count, dislikes count, duration, category ID, genre ID, category name, and genre name.
// It returns nil and an error if the query fails.
func (db *databaseConn) allVideoList() ([]VideoDescEdit, error) {
	// Prepare the SQL query
	query := "SELECT * FROM videos"

	// Execute the query
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the results
	var videos []VideoDescEdit

	// Iterate over the rows in the result set
	for rows.Next() {
		var video VideoDescEdit
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &video.UploaderID, &video.UploadDate, &video.ViewsCount, &video.LikesCount, &video.DislikesCount, &video.Duration, &video.CategoryID, &video.GenreID)
		video.GenreName, _ = db.getGenreName(video.GenreID)
		video.CategoryName, _ = db.getCategoryName(video.CategoryID)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	// Return the results
	return videos, nil
}

// saveUserProfile is a function that saves the profile image of a user.
// It takes in two parameters: imageFileName and userID.
// The imageFileName is the name of the image file to be saved.
// The userID is the ID of the user whose profile image is to be saved.
// It first checks if a user with the provided userID exists in the database.
// If the user exists, it updates the user's profile image in the database to the provided imageFileName.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) saveUserProfile(imageFileName string, userID int) error {
	// check if the user exists
	query := "SELECT COUNT(*) FROM users WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int

	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user does not exist")
	}

	// Check if the user has already uploaded a profile picture
	query = "SELECT ImageURL FROM userprofileimages WHERE UserID = ?"
	row = db.DB.QueryRow(query, userID)
	var imageURL string
	err = row.Scan(&imageURL)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// If the user has already uploaded a profile picture, update the record

	if imageFileName == "" {
		// If the user has not uploaded a profile picture, insert a new record for the first time
		query = "INSERT INTO userprofileimages (UserID) VALUES (?)"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
		return nil
	}

	if imageURL != "" {
		query = "UPDATE userprofileimages SET ImageURL = ? WHERE UserID = ?"
		_, err = db.DB.Exec(query, imageFileName, userID)
		if err != nil {
			return err
		}
		return nil
	} else {
		// If the user has not uploaded a profile picture, insert a new record
		query = "INSERT INTO userprofileimages (UserID, ImageURL) VALUES (?, ?)"
		_, err = db.DB.Exec(query, userID, imageFileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// displayUserProfileImage is a function that retrieves the profile image of a user.
// It takes in one parameter: userID.
// The userID is the ID of the user whose profile image is to be retrieved.
// It queries the database for the profile image URL of the user with the provided userID.
// It returns the profile image URL and nil if the user exists and has a profile image.
// It returns an empty string and an error if the user does not exist, does not have a profile image, or if the query fails.
func (db *databaseConn) displayUserProfileImage(userID int) (string, error) {
	query := "SELECT ImageURL FROM userprofileimages WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var imageURL string

	err := row.Scan(&imageURL)
	if err != nil {
		return "", err
	}
	return imageURL, nil
}

type access struct {
	UserId           int    `json:"user_id"`
	UserName         string `json:"user_name"`
	Email            string `json:"email"`
	DashboardAccess  int    `json:"dashboard_access"`
	UploadAccess     int    `json:"upload_access"`
	EditDeleteAccess int    `json:"edit_delete_access"`
	AnalyticsAccess  int    `json:"analytics_access"`
	ServerLogAccess  int    `json:"server_log_access"`
	UserAccess       int    `json:"user_access"`
}

// allUserAdminAccess is a function that retrieves the admin access details for all users.
// It takes in one parameter: userID.
// The userID is the ID of the user whose admin access details are to be retrieved.
// It queries the database for the admin access details of the user with the provided userID.
// The admin access details include user ID, user name, email, dashboard access, upload access, edit/delete access, analytics access, server log access, and user access.
// It returns a slice of 'access' structs and nil if the query is successful.
// Each 'access' struct contains the admin access details of a user.
// It returns nil and an error if the query fails.
func (db *databaseConn) allUserAdminAccess(userID int) ([]access, error) {
	// Initialize a slice to hold the access objects
	var accesses []access

	// Prepare the SQL query
	query := "SELECT U.UserID AS user_id, U.UserName AS user_name, U.Email AS email, UA.Dashboard AS dashboard_access, UA.Upload AS upload_access, UA.Edit_Delete AS edit_delete_access, UA.Analytics AS analytics_access, UA.ServerLogs AS server_log_access, UA.UserAccess AS user_access FROM Users U LEFT JOIN useraccesslevels UA ON U.UserID = UA.UserID ORDER BY U.UserName ASC"

	// Execute the query
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows in the result set
	for rows.Next() {
		var a access
		err := rows.Scan(&a.UserId, &a.UserName, &a.Email, &a.DashboardAccess, &a.UploadAccess, &a.EditDeleteAccess, &a.AnalyticsAccess, &a.ServerLogAccess, &a.UserAccess)
		if err != nil {
			return nil, err
		}

		// Append the access object to the slice
		accesses = append(accesses, a)
	}

	// Return the slice of access objects
	return accesses, nil
}

// giveDashboardAccess is a function that grants dashboard access to a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to be granted dashboard access.
// It updates the 'Dashboard' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively granting them dashboard access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) giveDashboardAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		query = "INSERT INTO useraccesslevels (UserID, Dashboard) VALUES (?, 1)"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	} else {
		query = "UPDATE useraccesslevels SET Dashboard = 1 WHERE UserID = ?"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// removeDashboardAccess is a function that revokes dashboard access from a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to have their dashboard access revoked.
// It updates the 'Dashboard' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively revoking their dashboard access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) removeDashboardAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user does not have access")
	}
	query = "UPDATE useraccesslevels SET Dashboard = 0 WHERE UserID = ?"
	_, err = db.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

// giveUploadAccess is a function that grants upload access to a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to be granted upload access.
// It updates the 'Upload' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively granting them upload access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) giveUploadAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		query = "INSERT INTO useraccesslevels (UserID, Upload) VALUES (?, 1)"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	} else {
		query = "UPDATE useraccesslevels SET Upload = 1 WHERE UserID = ?"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// removeUploadAccess is a function that revokes upload access from a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to have their upload access revoked.
// It updates the 'Upload' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively revoking their upload access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) removeUploadAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user does not have access")
	}
	query = "UPDATE useraccesslevels SET Upload = 0 WHERE UserID = ?"
	_, err = db.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

// giveEditDeleteAccess is a function that grants edit and delete access to a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to be granted edit and delete access.
// It updates the 'Edit_Delete' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively granting them edit and delete access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) giveEditDeleteAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		query = "INSERT INTO useraccesslevels (UserID, Edit_Delete) VALUES (?, 1)"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	} else {
		query = "UPDATE useraccesslevels SET Edit_Delete = 1 WHERE UserID = ?"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// removeEditDeleteAccess is a function that revokes edit and delete access from a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to have their edit and delete access revoked.
// It updates the 'Edit_Delete' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively revoking their edit and delete access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) removeEditDeleteAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user does not have access")
	}
	query = "UPDATE useraccesslevels SET Edit_Delete = 0 WHERE UserID = ?"
	_, err = db.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

// giveAnalyticsAccess is a function that grants analytics access to a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to be granted analytics access.
// It updates the 'Analytics' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively granting them analytics access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) giveAnalyticsAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		query = "INSERT INTO useraccesslevels (UserID, Analytics) VALUES (?, 1)"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	} else {
		query = "UPDATE useraccesslevels SET Analytics = 1 WHERE UserID = ?"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// removeAnalyticsAccess is a function that revokes analytics access from a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to have their analytics access revoked.
// It updates the 'Analytics' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively revoking their analytics access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) removeAnalyticsAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user does not have access")
	}
	query = "UPDATE useraccesslevels SET Analytics = 0 WHERE UserID = ?"
	_, err = db.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

// giveServerLogAccess is a function that grants server log access to a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to be granted server log access.
// It updates the 'Server_Log' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively granting them server log access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) giveServerLogAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		query = "INSERT INTO useraccesslevels (UserID, ServerLogs) VALUES (?, 1)"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	} else {
		query = "UPDATE useraccesslevels SET ServerLogs = 1 WHERE UserID = ?"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// removeServerLogAccess is a function that revokes server log access from a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to have their server log access revoked.
// It updates the 'Server_Log' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively revoking their server log access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) removeServerLogAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user does not have access")
	}
	query = "UPDATE useraccesslevels SET ServerLogs = 0 WHERE UserID = ?"
	_, err = db.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

// giveUserAccess is a function that grants general user access to a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to be granted user access.
// It updates the 'User_Access' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively granting them user access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) giveUserAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		query = "INSERT INTO useraccesslevels (UserID, UserAccess) VALUES (?, 1)"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	} else {
		query = "UPDATE useraccesslevels SET UserAccess = 1 WHERE UserID = ?"
		_, err = db.DB.Exec(query, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// removeUserAccess is a function that revokes general user access from a user.
// It takes in one parameter: userID.
// The userID is the ID of the user who is to have their user access revoked.
// It updates the 'User_Access' field in the 'useraccesslevels' table in the database for the user with the provided userID, effectively revoking their user access.
// It returns an error if the user does not exist or if the update fails.
func (db *databaseConn) removeUserAccess(userID int) error {
	query := "SELECT COUNT(*) FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user does not have access")
	}
	query = "UPDATE useraccesslevels SET UserAccess = 0 WHERE UserID = ?"
	_, err = db.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

// userAccess is a function that retrieves the access level of a user.
// It takes in one parameter: userID.
// The userID is the ID of the user whose access level is to be retrieved.
// It queries the 'useraccesslevels' table in the database for the access level of the user with the provided userID.
// It returns the access level and nil if the user exists and has an access level.
// It returns an empty string and an error if the user does not exist, does not have an access level, or if the query fails.
func (db *databaseConn) userAccess(userID int) (int, int, int, int, int, int, error) {
	query := "SELECT Dashboard, Upload, Edit_Delete, Analytics, ServerLogs, UserAccess FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var dashboard, upload, editDelete, analytics, serverLogs, userAccess int
	err := row.Scan(&dashboard, &upload, &editDelete, &analytics, &serverLogs, &userAccess)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	return dashboard, upload, editDelete, analytics, serverLogs, userAccess, nil
}

// updateVideo is a function that updates the details of a video.
// It takes in several parameters including videoID, title, description, category, and genre.
// The videoID is the ID of the video to be updated.
// The title, description, category, and genre are the new details to be updated in the video record.
// It updates the video record in the database with the provided parameters.
// It returns an error if the video does not exist or if the update fails.
func (db *databaseConn) updateVideo(title string, description string, category string, genre string, videoID int) error {
	categoryID, err := db.getCategoryID(genre)
	if err != nil {
		return err
	}
	genreID, err := db.getGenreID(category)
	if err != nil {
		return err
	}

	fmt.Println("categoryID: ", categoryID)
	fmt.Println("genreID: ", genreID)

	query := "UPDATE videos SET Title = ?, Description = ?, CategoryID = ?, GenreID = ? WHERE VideoID = ?"
	_, err = db.DB.Exec(query, title, description, categoryID, genreID, videoID)
	if err != nil {
		fmt.Println("error: ", err)
		return err
	}
	return nil
}

// checkAdminAccess is a function that checks if a user has administrative access.
// It takes in one parameter: userID.
// The userID is the ID of the user whose admin access status is to be checked.
// It queries the 'useraccesslevels' table in the database for the admin access status of the user with the provided userID.
// It returns the admin access status and nil if the user exists and has an admin access status.
// It returns false and an error if the user does not exist, does not have an admin access status, or if the query fails.
func (db *databaseConn) checkAdminAccess(userID int) (int, error) {
	//SELECT
	//CASE
	//WHEN Dashboard = 1 OR Upload = 1 OR Edit_Delete = 1 OR Analytics = 1 OR ServerLogs = 1 OR UserAccess = 1 THEN 1
	//ELSE 0
	//END AS AnyValueIs1
	//FROM
	//useraccesslevels;

	query := "SELECT CASE WHEN Dashboard = 1 OR Upload = 1 OR Edit_Delete = 1 OR Analytics = 1 OR ServerLogs = 1 OR UserAccess = 1 THEN 1 ELSE 0 END AS AnyValueIs1 FROM useraccesslevels WHERE UserID = ?"
	row := db.DB.QueryRow(query, userID)
	var anyValueIs1 int
	err := row.Scan(&anyValueIs1)
	if err != nil {
		return 0, err
	}
	return anyValueIs1, nil
}
