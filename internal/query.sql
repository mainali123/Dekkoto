-- Create database dekkoto
CREATE DATABASE dekkoto;

-- Create Users table
CREATE TABLE Users
(
    UserID   INT PRIMARY KEY AUTO_INCREMENT,
    UserName VARCHAR(50),
    Email    VARCHAR(100),
    Password VARCHAR(100),
    Admin    BOOLEAN DEFAULT FALSE
);

-- Create Categories table
CREATE TABLE Categories
(
    CategoryID   INT PRIMARY KEY AUTO_INCREMENT,
    CategoryName VARCHAR(50)
);

CREATE TABLE Genres
(
    GenreID   INT PRIMARY KEY AUTO_INCREMENT,
    GenreName VARCHAR(50)
);

-- Create Videos table
CREATE TABLE Videos
(
    VideoID       INT PRIMARY KEY AUTO_INCREMENT,
    Title         VARCHAR(100),
    Description   TEXT,
    URL           VARCHAR(255),
    ThumbnailURL  VARCHAR(255),
    UploaderID    INT,
    UploadDate    DATE,
    ViewsCount    INT DEFAULT 0,
    LikesCount    INT DEFAULT 0,
    DislikesCount INT DEFAULT 0,
    Duration      VARCHAR(10),
    CategoryID    INT,
    GenreID       INT,
    FOREIGN KEY (UploaderID) REFERENCES Users (UserID),
    FOREIGN KEY (CategoryID) REFERENCES Categories (CategoryID),
    FOREIGN KEY (GenreID) REFERENCES Genres (GenreID)
);

-- Create VideoActions table
CREATE TABLE VideoActions
(
    VideoActionID INT PRIMARY KEY AUTO_INCREMENT,
    UserID        INT,
    VideoID       INT,
    Recommends    TINYINT DEFAULT 0,
    Watching      TINYINT DEFAULT 0,
    Completed     TINYINT DEFAULT 0,
    On_hold       TINYINT DEFAULT 0,
    Considering   TINYINT DEFAULT 0,
    Dropped       TINYINT DEFAULT 0,
    ActionsDate   DATE,
    ActionTime    TIME,
    FOREIGN KEY (UserID) REFERENCES Users (UserID),
    FOREIGN KEY (VideoID) REFERENCES Videos (VideoID)
);

-- Create Comments table
CREATE TABLE Comments
(
    CommentID   INT PRIMARY KEY AUTO_INCREMENT,
    UserID      INT,
    VideoID     INT,
    CommentText TEXT,
    CommentDate DATE,
    FOREIGN KEY (UserID) REFERENCES Users (UserID),
    FOREIGN KEY (VideoID) REFERENCES Videos (VideoID)
);


-- Insert Action, Comedy, Drama, Horror, Romance, Thriller in categories table
# INSERT INTO Categories (CategoryName) VALUES ('Action'), ('Comedy'), ('Drama'), ('Horror'), ('Romance'), ('Thriller');

INSERT INTO Categories (CategoryName)
VALUES ('action'),
       ('comedy'),
       ('drama'),
       ('horror'),
       ('romance'),
       ('thriller'),
       ('action,comedy'),
       ('action,drama'),
       ('action,horror'),
       ('action,romance'),
       ('action,thriller'),
       ('comedy,drama'),
       ('comedy,horror'),
       ('comedy,romance'),
       ('comedy,thriller'),
       ('drama,horror'),
       ('drama,romance'),
       ('drama,thriller'),
       ('horror,romance'),
       ('horror,thriller'),
       ('romance,thriller'),
       ('action,comedy,drama'),
       ('action,comedy,horror'),
       ('action,comedy,romance'),
       ('action,comedy,thriller'),
       ('action,drama,horror'),
       ('action,drama,romance'),
       ('action,drama,thriller'),
       ('action,horror,romance'),
       ('action,horror,thriller'),
       ('action,romance,thriller'),
       ('comedy,drama,horror'),
       ('comedy,drama,romance'),
       ('comedy,drama,thriller'),
       ('comedy,horror,romance'),
       ('comedy,horror,thriller'),
       ('comedy,romance,thriller'),
       ('drama,horror,romance'),
       ('drama,horror,thriller'),
       ('drama,romance,thriller'),
       ('horror,romance,thriller'),
       ('action,comedy,drama,horror'),
       ('action,comedy,drama,romance'),
       ('action,comedy,drama,thriller'),
       ('action,comedy,horror,romance'),
       ('action,comedy,horror,thriller'),
       ('action,comedy,romance,thriller'),
       ('action,drama,horror,romance'),
       ('action,drama,horror,thriller'),
       ('action,drama,romance,thriller'),
       ('action,horror,romance,thriller'),
       ('comedy,drama,horror,romance'),
       ('comedy,drama,horror,thriller'),
       ('comedy,drama,romance,thriller'),
       ('comedy,horror,romance,thriller'),
       ('drama,horror,romance,thriller'),
       ('action,comedy,drama,horror,romance'),
       ('action,comedy,drama,horror,thriller'),
       ('action,comedy,drama,romance,thriller'),
       ('action,comedy,horror,romance,thriller'),
       ('action,drama,horror,romance,thriller'),
       ('comedy,drama,horror,romance,thriller'),
       ('action,comedy,drama,horror,romance,thriller');


-- Insert Movie, Series, Anime in genres table
INSERT INTO Genres (GenreName)
VALUES ('movie'),
       ('series'),
       ('anime');

CREATE TABLE CommentActions
(
    CommentActionID INT PRIMARY KEY AUTO_INCREMENT,
    Upvotes TINYINT DEFAULT 0 CHECK (Upvotes IN (0, 1)),
    Downvotes TINYINT DEFAULT 0 CHECK (Downvotes IN (0, 1)),
    CommentID INT,
    UserID INT,
    FOREIGN KEY (CommentID) REFERENCES Comments (CommentID),
    FOREIGN KEY (UserID) REFERENCES Users (UserID)
);