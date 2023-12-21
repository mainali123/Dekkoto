-- Create Users table
CREATE TABLE Users (
                       UserID INT PRIMARY KEY AUTO_INCREMENT,
                       UserName VARCHAR(50),
                       Email VARCHAR(100),
                       Password VARCHAR(100),
                       Admin BOOLEAN DEFAULT FALSE
);

-- Create Categories table
CREATE TABLE Categories (
                            CategoryID INT PRIMARY KEY,
                            CategoryName VARCHAR(50)
);

-- Create Videos table
CREATE TABLE Videos (
                        VideoID INT PRIMARY KEY,
                        Title VARCHAR(100),
                        Description TEXT,
                        URL VARCHAR(255),
                        UploaderID INT,
                        UploadDate DATE,
                        ViewsCount INT DEFAULT 0,
                        LikesCount INT DEFAULT 0,
                        DislikesCount INT DEFAULT 0,
                        Duration TIME,
                        CategoryID INT,
                        FOREIGN KEY (UploaderID) REFERENCES Users(UserID),
                        FOREIGN KEY (CategoryID) REFERENCES Categories(CategoryID)
);

-- Create VideoActions table
CREATE TABLE VideoActions (
                              ActionID INT PRIMARY KEY,
                              UserID INT,
                              VideoID INT,
                              ActionType VARCHAR(20),
                              FOREIGN KEY (UserID) REFERENCES Users(UserID),
                              FOREIGN KEY (VideoID) REFERENCES Videos(VideoID)
);

-- Create Comments table
CREATE TABLE Comments (
                          CommentID INT PRIMARY KEY,
                          UserID INT,
                          VideoID INT,
                          CommentText TEXT,
                          Upvotes INT DEFAULT 0,
                          Downvotes INT DEFAULT 0,
                          FOREIGN KEY (UserID) REFERENCES Users(UserID),
                          FOREIGN KEY (VideoID) REFERENCES Videos(VideoID)
);
