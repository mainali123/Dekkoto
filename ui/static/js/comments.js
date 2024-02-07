console.log('comments.js loaded')
let loaded = 0;

// Select the h1 tag where the title will be displayed
let titleElement = document.querySelector('.title h1');
// Get the title from the videoDetails object
let videoTitle = videoDetails.Title;
// Set the innerText property of the h1 tag to the video title
titleElement.innerText = videoTitle;

// Select the p tag where the description will be displayed
let descriptionElement = document.querySelector('.description .desc');
// Get the description from the videoDetails object
let videoDescription = videoDetails.Description;
// Set the innerText property of the p tag to the video description
descriptionElement.innerText = videoDescription;


// COMMENT POST FOR THE BACKEND
// Select the comment input field and the post button
const commentInput = document.querySelector('#User-comments');
const postButton = document.getElementById('comment-on-video');


// Add an event listener to the post button
postButton.addEventListener('click', function (event) {
    // Prevent the default form submission
    event.preventDefault();

    // Get the value of the comment input field
    const comment = commentInput.value;
    console.log(comment)

    if (comment.length == 0) {
        alert('Please enter a comment');
        return;
    }


    // Send a POST request to the '/comment' endpoint
    fetch('/comment', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({comment: comment, videoID: videoDetails.VideoID})
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // Display the error message to the user
                alert(data.error);
            } else {
                // Clear the comment input field and display the success message
                commentInput.value = '';
                alert('Comment posted successfully');
            }
        })
        .catch(error => {
            // Display the error message to the user
            alert('An error occurred: ' + error);
        });
});

// Fetch comments from the server


function comments() {
    fetch('/displayComments', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
        .then(response => response.json())
        .then(data => {
            const comments = data.comments;
            console.log(comments);

            // Select the comments section in the HTML
            const commentsSection = document.querySelector('.comments-section');

            // Clear the comments section
            // commentsSection.innerHTML = '';

            // For each comment, create a new div and append it to the comments section
            comments.forEach(comment => {
                const commentDiv = document.createElement('div');
                commentDiv.classList.add('other-comments');
                let commentVotes = comment.Upvotes - comment.Downvotes;
                commentDiv.innerHTML = `
                <div class="comment">
                    <div class="other-user-name">
                        <div class="userrrrr">
                            <p class="name">${comment.UserName}</p>
                            <span class="time">${new Date(comment.CommentDate).toLocaleDateString()}</span>
                        </div>
                        <img src="../static/images/boruto.jpg" alt="default profile" class="other-user-profile" />
                    </div>
                    <p class="other-user-comment">
                        ${comment.CommentText}
                    </p>
                    <div class="upvotes">
                        <img src="../static/images/upDrop.svg" data-id="${comment.CommentID}" alt="" class="upVotes">
                        <span class="voteNumbers" style="color: ${commentVotes < 0 ? '#FF0000' : (commentVotes > 0 ? '#00FF0A' : '')}">${commentVotes}</span>
                        <img src="../static/images/downDrop.svg" data-id="${comment.CommentID}" alt="" class="downVotes">
                    </div>
                </div>
                `;
                commentsSection.appendChild(commentDiv);
            });

            // Add event listeners to the upvote and downvote buttons
            upvote_downvote();
            getCommentDetails();
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function upvote_downvote() {
    console.log('upvote_downvote function loaded')
    const upvote = document.querySelectorAll('.upVotes');
    const downvote = document.querySelectorAll('.downVotes');

    upvote.forEach((upvote) => {
        upvote.addEventListener('click', (event) => {
            console.log("upvote clicked")
            // convert the commentID to a number
            comment_ID = parseInt(event.target.dataset.id);

            fetch('/upvote', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({commentID: comment_ID})
            })
                .then(response => response.json())
                .then(data => {
                    console.log(data)

                    // If the message is 'Comment is already upvoted', alert the user
                    if (data.message === 'Comment is already upvoted') {
                        alert(data.message);
                    }

                    if (!data.success) {
                        // Display the error message to the user
                        alert(data.message);
                    }
                })
                .catch(error => {
                    // Display the error message to the user
                    alert('An error occurred: ' + error);
                });
        });
    });

    downvote.forEach((downvote) => {
        downvote.addEventListener('click', (event) => {
            console.log("downvote clicked")
            // convert the commentID to a number
            let comment_ID = parseInt(event.target.dataset.id);

            fetch('/downvote', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({commentID: comment_ID})
            })
                .then(response => response.json())
                .then(data => {
                    console.log(data)

                    // If the message is 'Comment is already downvoted', alert the user
                    if (data.message === 'Comment is already downvoted') {
                        alert(data.message);
                    }

                    if (!data.success) {
                        // Display the error message to the user
                        alert(data.message);
                    }
                })
                .catch(error => {
                    // Display the error message to the user
                    alert('An error occurred: ' + error);
                });
        });
    });
}

comments();

// for turning upvote colour to green and downvote colour to red
function getCommentDetails() {
    let json = JSON.stringify({videoID: videoDetails.VideoID});
    fetch('/commentDetails', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: json
    })
    .then(response => response.json())
    .then(data => {
        console.log(data);
        data.comments.forEach(commentDetail => {
            let upvoteIcon = document.querySelector(`.upVotes[data-id="${commentDetail.CommentID}"]`);
            let downvoteIcon = document.querySelector(`.downVotes[data-id="${commentDetail.CommentID}"]`);
            if (upvoteIcon && commentDetail.Upvote > commentDetail.Downvote) {
                upvoteIcon.src = "../static/images/upvote_green.svg";
            } else if (downvoteIcon && commentDetail.Downvote > commentDetail.Upvote) {
                downvoteIcon.src = "../static/images/downvote_red.svg";
            } else if (upvoteIcon && downvoteIcon && commentDetail.Upvote === commentDetail.Downvote) {
                upvoteIcon.src = "../static/images/upvote.svg";
                downvoteIcon.src = "../static/images/downvote.svg";
            }
        });
    })
}


// Like and Dislike functionality

/*function likeVideo() {

    fetch('/likeVideo', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            if (data.message === 'Video is already liked') {
                alert(data.message);
            }
        })
        .catch(error => {
            alert('An error occurred: ' + error);
        });
}

function reverseLikeVideo() {

    fetch('/reverseLikeVideo', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            if (data.message === 'Video is already liked') {
                alert(data.message);
            }
        })
        .catch(error => {
            alert('An error occurred: ' + error);
        });
}

function dislikeVideo() {

    fetch('/dislikeVideo', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            if (data.message === 'Video is already disliked') {
                alert(data.message);
            }
        })
        .catch(error => {
            alert('An error occurred: ' + error);
        });
}

function reverseDislikeVideo() {

        fetch('/reverseDislikeVideo', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({videoID: videoDetails.VideoID})
        })
            .then(response => response.json())
            .then(data => {
                console.log(data)
                if (data.message === 'Video is already disliked') {
                    alert(data.message);
                }
            })
            .catch(error => {
                alert('An error occurred: ' + error);
            });
}

function isLikedDisliked() {
    fetch('/isLikedDisliked', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            if (data.isLiked) {
                fontLike.classList.remove('fa-regular');
                fontLike.classList.add('fa-solid', 'fa-thumbs-up');
                likeBtn.classList.add('active');
            } else if (data.isDisliked) {
                fontDislike.classList.remove('fa-regular');
                fontDislike.classList.add('fa-solid', 'fa-thumbs-down');
                disLikeBtn.classList.add('active');
            }
        })
        .catch(error => {
            alert('An error occurred: ' + error);
        });
}

function like_dislike() {

    likeBtn.addEventListener('click', () => {
        // Check if the dislike button is clicked
        if (disLikeBtn.classList.contains('active')) {
            // Undo the dislike button click
            disLikeBtn.classList.remove('active');
            fontDislike.classList.remove('fa-solid', 'fa-thumbs-down');
            fontDislike.classList.add('fa-regular', 'fa-thumbs-down');
        }

        // Toggle the like button click
        if (fontLike.classList.contains('fa-regular')) {
            fontLike.classList.remove('fa-regular');
            fontLike.classList.add('fa-solid', 'fa-thumbs-up');
            likeBtn.classList.add('active');
        } else {
            fontLike.classList.remove('fa-solid', 'fa-thumbs-up');
            fontLike.classList.add('fa-regular', 'fa-thumbs-up');
            likeBtn.classList.remove('active');
        }
    });

    disLikeBtn.addEventListener('click', () => {
        // Check if the like button is clicked
        if (likeBtn.classList.contains('active')) {
            // Undo the like button click
            likeBtn.classList.remove('active');
            fontLike.classList.remove('fa-solid', 'fa-thumbs-up');
            fontLike.classList.add('fa-regular', 'fa-thumbs-up');
        }

        // Toggle the dislike button click
        if (fontDislike.classList.contains('fa-regular')) {
            fontDislike.classList.remove('fa-regular');
            fontDislike.classList.add('fa-solid', 'fa-thumbs-down');
            disLikeBtn.classList.add('active');
        } else {
            fontDislike.classList.remove('fa-solid', 'fa-thumbs-down');
            fontDislike.classList.add('fa-regular', 'fa-thumbs-down');
            disLikeBtn.classList.remove('active');
        }
    });
}*/

// like_dislike();


// Select the like and dislike buttons
const likeBtn = document.querySelector('.like');
const dislikeBtn = document.querySelector('.dislike');
const fontLike = document.getElementById('font-awe-like');
const fontDislike = document.getElementById('font-awe-dislike');

likeBtn.addEventListener('click', function() {
    // Check if the video is already liked
    if (likeBtn.classList.contains('active')) {
        // If the video is already liked, send a request to unlike the video
        likeBtn.classList.remove('active');
        reverseLikeVideo();
    } else {
        // If the video is not liked, send a request to like the video
        likeBtn.classList.add('active');
        dislikeBtn.classList.remove('active');
        likeVideo();
    }
});

dislikeBtn.addEventListener('click', function() {
    // Check if the video is already disliked
    if (dislikeBtn.classList.contains('active')) {
        // If the video is already disliked, send a request to undislike the video
        dislikeBtn.classList.remove('active');
        reverseDislikeVideo();
    } else {
        // If the video is not disliked, send a request to dislike the video
        dislikeBtn.classList.add('active');
        likeBtn.classList.remove('active');
        dislikeVideo();
    }
});

// Function to like a video
function likeVideo() {
    fetch('/likeVideo', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
    .then(response => response.json())
    .then(data => {
        // If the video is successfully liked, update the UI
        if (data.success) {
            likeBtn.classList.add('active');
            dislikeBtn.classList.remove('active');
        }
    });
}

// Function to unlike a video
function reverseLikeVideo() {
    fetch('/reverseLikeVideo', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
    .then(response => response.json())
    .then(data => {
        // If the video is successfully unliked, update the UI
        if (data.success) {
            likeBtn.classList.remove('active');
        }
    });
}

// Function to dislike a video
function dislikeVideo() {
    fetch('/dislikeVideo', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
    .then(response => response.json())
    .then(data => {
        // If the video is successfully disliked, update the UI
        if (data.success) {
            dislikeBtn.classList.add('active');
            likeBtn.classList.remove('active');
        }
    });
}

// Function to undislike a video
function reverseDislikeVideo() {
    fetch('/reverseDislikeVideo', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
    .then(response => response.json())
    .then(data => {
        // If the video is successfully undisliked, update the UI
        if (data.success) {
            dislikeBtn.classList.remove('active');
        }
    });
}

function isLikedDisliked() {
    fetch('/isLikedDisliked', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            if (data.isLiked) {
                fontLike.classList.remove('fa-regular');
                fontLike.classList.add('fa-solid', 'fa-thumbs-up');
                likeBtn.classList.add('active');
            } else if (data.isDisliked) {
                fontDislike.classList.remove('fa-regular');
                fontDislike.classList.add('fa-solid', 'fa-thumbs-down');
                dislikeBtn.classList.add('active');
            }
        })
        .catch(error => {
            alert('An error occurred: ' + error);
        });
}

isLikedDisliked();

function like_dislike_count() {
    fetch('/likeDislikeCount', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({videoID: videoDetails.VideoID})
    })
    .then(response => response.json())
    .then(data => {
        console.log(data)
        document.querySelector('.liked-video-count').innerText = data.likes;
        document.querySelector('.dislike-video-count').innerText = data.dislikes;
    })
    .catch(error => {
        alert('An error occurred: ' + error);
    });
}

like_dislike_count();