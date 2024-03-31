/*
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

let loadOnce = true;

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
            if (loadOnce) {
                getCommentDetails();
                loadOnce = false;
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

/!*function upvote_downvote() {
    console.log('upvote_downvote function loaded')
    const upvote = document.querySelectorAll('.upVotes');
    const downvote = document.querySelectorAll('.downVotes');

    upvote.forEach((upvote) => {
        upvote.addEventListener('click', (event) => {
            console.log("upvote clicked")
            // convert the commentID to a number
            comment_ID = parseInt(event.target.dataset.id);

            /!*fetch('/upvote', {
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
                });*!/

            // upvote a video
            $.ajax({
                url: '/upvote',
                type: 'POST',
                data: JSON.stringify({commentID: comment_ID}),
                success: function(data) {
                    console.log(data)

                    if (data.message === 'Comment is already upvoted') {
                        alert(data.message);
                    }
                    if (!data.success) {
                        alert(data.message);
                    }
                },
                error: function (xhr, status, error) {
                    console.log('An error occurred: ' + error)
                }
            })
        });
    });

    downvote.forEach((downvote) => {
        downvote.addEventListener('click', (event) => {
            console.log("downvote clicked")
            // convert the commentID to a number
            let comment_ID = parseInt(event.target.dataset.id);

            /!*fetch('/downvote', {
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
                    getCommentDetails();
                })
                .catch(error => {
                    // Display the error message to the user
                    alert('An error occurred: ' + error);
                });*!/

            $.ajax({
                url: '/downvote',
                type: 'POST',
                data: JSON.stringify({commentID: comment_ID}),
                success: function(data) {
                    console.log(data)

                    if (data.message === 'Comment is already downvoted') {
                        alert(data.message);
                    }

                    if (!data.success) {
                        alert(data.message);
                    }
                    getCommentDetails();
                },
                error: function (xhr, status, error) {
                    console.log('An error occurred: ' + error)
                }
            })
        });
    });
}*!/

function upvote_downvote() {
    console.log('upvote_downvote function loaded');
    const upvote = document.querySelectorAll('.upVotes');
    const downvote = document.querySelectorAll('.downVotes');

    let upvoteSrc = upvote[0].src.split('/')[upvote[0].src.split('/').length - 1];
    let downVoteSrc = downvote[0].src.split('/')[downvote[0].src.split('/').length - 1];

    console.log(upvoteSrc, downVoteSrc);

    let voteNumber = document.querySelector('.voteNumbers');

    upvote.forEach((upvote) => {
        upvote.addEventListener('click', (event) => {
            const fullUrl = event.target.src;
            const urlParts = fullUrl.split('/');
            const imageName = urlParts[urlParts.length - 1];
            let comment_ID = parseInt(event.target.dataset.id);

            if (imageName === 'upvote_green.svg') {
                console.log('upvoted');

                $.ajax({
                    url: '/reverseUpvote',
                    type: 'POST',
                    data: JSON.stringify({commentID: comment_ID}),
                    success: function (data) {
                        if (!data.success) {
                            alert(data.message);
                        } else {
                            let displayNumber = parseInt(voteNumber.innerText);
                            event.target.src = "../static/images/upDrop.svg";
                            voteNumber.innerText = displayNumber - 1;
                            upvoteSrc = 'upDrop.svg';
                        }
                    },
                    error: function (xhr, status, error) {
                        console.log('An error occurred: ' + error);
                    }
                });

            } else if (imageName === 'upDrop.svg') {
                console.log('not upvoted');

                $.ajax({
                    url: '/upvote',
                    type: 'POST',
                    data: JSON.stringify({commentID: comment_ID}),
                    success: function (data) {
                        if (!data.success) {
                            alert(data.message);
                        } else {
                            let displayNumber = parseInt(voteNumber.innerText);
                            event.target.src = "../static/images/upvote_green.svg";
                            if (downVoteSrc === 'downvote_red.svg') {
                                voteNumber.innerText = displayNumber + 2;
                            } else {
                                voteNumber.innerText = displayNumber + 1;
                            }
                            // Change downvote images to default
                            downvote.forEach(downvote => {
                                downvote.src = "../static/images/downDrop.svg";
                            });
                            upvoteSrc = 'upvote_green.svg';
                        }
                    },
                    error: function (xhr, status, error) {
                        console.log('An error occurred: ' + error);
                    }
                });
            }
        });
    });

    downvote.forEach((downvote) => {
        downvote.addEventListener('click', (event) => {
            const fullUrl = event.target.src;
            const urlParts = fullUrl.split('/');
            const imageName = urlParts[urlParts.length - 1];
            let comment_ID = parseInt(event.target.dataset.id);

            if (imageName === 'downvote_red.svg') {
                console.log('downvoted');

                $.ajax({
                    url: '/reverseDownvote',
                    type: 'POST',
                    data: JSON.stringify({commentID: comment_ID}),
                    success: function (data) {
                        if (!data.success) {
                            alert(data.message);
                        } else {
                            let displayNumber = parseInt(voteNumber.innerText);
                            event.target.src = "../static/images/downDrop.svg";
                            voteNumber.innerText = displayNumber + 1;
                            downVoteSrc = 'downDrop.svg';
                        }
                    },
                    error: function (xhr, status, error) {
                        console.log('An error occurred: ' + error);
                    }
                });

            } else if (imageName === 'downDrop.svg') {
                console.log('not downvoted');

                $.ajax({
                    url: '/downvote',
                    type: 'POST',
                    data: JSON.stringify({commentID: comment_ID}),
                    success: function (data) {
                        if (!data.success) {
                            alert(data.message);
                        } else {
                            let displayNumber = parseInt(voteNumber.innerText);
                            event.target.src = "../static/images/downvote_red.svg";
                            if (upvoteSrc === 'upvote_green.svg') {
                                voteNumber.innerText = displayNumber - 2;
                            } else {
                                voteNumber.innerText = displayNumber - 1;
                            }
                            // Change upvote images to default
                            upvote.forEach(upvote => {
                                upvote.src = "../static/images/upDrop.svg";
                            });
                            downVoteSrc = 'downvote_red.svg';
                        }
                    },
                    error: function (xhr, status, error) {
                        console.log('An error occurred: ' + error);
                    }
                });
            }
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
                    upvoteIcon.src = "../static/images/upDrop.svg";
                    downvoteIcon.src = "../static/images/downDrop.svg";
                }
            });
        })
}


// like_dislike();


// Select the like and dislike buttons
const likeBtn = document.querySelector('.like');
const dislikeBtn = document.querySelector('.dislike');
const fontLike = document.getElementById('font-awe-like');
const fontDislike = document.getElementById('font-awe-dislike');

likeBtn.addEventListener('click', function () {
    // Check if the video is already liked
    if (likeBtn.classList.contains('active')) {
        // If the video is already liked, send a request to unlike the video
        likeBtn.classList.remove('active');
        fontLike.classList.remove('fa-solid')
        fontLike.classList.add('fa-regular')
        reverseLikeVideo();
    } else {
        // If the video is not liked, send a request to like the video
        likeBtn.classList.add('active');
        dislikeBtn.classList.remove('active');
        fontLike.classList.add('fa-solid')
        fontLike.classList.remove('fa-regular')

        fontDislike.classList.remove('fa-solid')
        fontDislike.classList.add('fa-regular')
        likeVideo();
    }
});

dislikeBtn.addEventListener('click', function () {
    // Check if the video is already disliked
    if (dislikeBtn.classList.contains('active')) {
        // If the video is already disliked, send a request to undislike the video
        dislikeBtn.classList.remove('active');
        fontDislike.classList.remove('fa-solid')
        fontDislike.classList.add('fa-regular')
        reverseDislikeVideo();
    } else {
        // If the video is not disliked, send a request to dislike the video
        dislikeBtn.classList.add('active');
        likeBtn.classList.remove('active');
        fontDislike.classList.add('fa-solid')
        fontDislike.classList.remove('fa-regular')

        fontLike.classList.remove('fa-solid')
        fontLike.classList.add('fa-regular')
        dislikeVideo();
    }
});

// Function to like a video

function likeVideo() {
    console.log(JSON.stringify({videoID: videoDetails.VideoID}))
    $.ajax({
        url: '/likeVideo',
        type: 'POST',
        data: JSON.stringify({videoID: videoDetails.VideoID}),
        success: function (data) {
            // If the video is successfully liked, update the UI
            if (data.success) {
                likeBtn.classList.add('active');
                dislikeBtn.classList.remove('active');
                like_dislike_count();
            }
        },
        error: function (xhr, status, error) {
            console.log('An error occurred: ' + error)
        }
    })
}

// Function to unlike a video

function reverseLikeVideo() {
    $.ajax({
        url: '/reverseLikeVideo',
        type: 'POST',
        data: JSON.stringify({videoID: videoDetails.VideoID}),
        success: function (data) {
            // If the video is successfully unliked, update the UI
            if (data.success) {
                likeBtn.classList.remove('active');
                like_dislike_count();
            }
        },
        error: function (xhr, status, error) {
            console.log('An error occurred: ' + error)
        }
    })
}

// Function to dislike a video

function dislikeVideo() {
    $.ajax({
        url: '/dislikeVideo',
        type: 'POST',
        data: JSON.stringify({videoID: videoDetails.VideoID}),
        success: function (data) {
            dislikeBtn.classList.add('active')
            likeBtn.classList.remove('active')
            like_dislike_count();
        },
        error: function (xhr, status, error) {
            console.log("Error: " + error)
        }
    })
}

// Function to undislike a video
function reverseDislikeVideo() {
    $.ajax({
        url: '/reverseDislikeVideo',
        type: 'POST',
        data: JSON.stringify({videoID: videoDetails.VideoID}),
        success: function (data) {
            dislikeBtn.classList.remove('active')
            like_dislike_count();
        },
        error: function (xhr, status, error) {
            console.log("Error: " + error)
        }
    })
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
                fontLike.classList.add('fa-solid');
                likeBtn.classList.add('active');
                console.log('liked');
            } else if (data.isDisliked) {
                fontDislike.classList.remove('fa-regular');
                fontDislike.classList.add('fa-solid');
                dislikeBtn.classList.add('active');
                console.log('disliked');
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
            // isLikedDisliked();
        })
        .catch(error => {
            alert('An error occurred: ' + error);
        });
}

like_dislike_count();*/
