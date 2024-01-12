console.log("hello world")

// videoDetails:"{"VideoID":39,"Title":"Nezuko","Description":"desc","URL":"./userUploadDatas/videos/20240109161922.406048_0b177be7-2f42-4dcf-b4de-2399eb452b29_323500696560_encoded.mp4","ThumbnailURL":"./userUploadDatas/thumbnails/20240109161922.406048_0b177be7-2f42-4dcf-b4de-2399eb452b29_323500696560.png","UploaderID":1,"UploadDate":"2024-01-09T00:00:00Z","ViewsCount":0,"LikesCount":0,"DislikesCount":0,"Duration":"00:00:26","CategoryID":6,"GenreID":1}"
// get the video details from local storage

let videoDetails = localStorage.getItem('videoDetails');
// convert to JSON
videoDetails = JSON.parse(videoDetails);
console.log(videoDetails)

// Select the video player and source elements
let videoPlayer = document.getElementById('videoPlayer');
let videoSource = document.getElementById('videoSource');

// Set the source of the video
videoSource.src = videoDetails.URL;
console.log(videoSource.src)

// Load the video
videoPlayer.load();

// Try to play the video
videoPlayer.play().catch(error => {
    // Auto-play was prevented
    // Show a UI element to let the user manually start the playback
    console.log("Auto-play was prevented by the browser. Please start the video manually.");
});

document.addEventListener('DOMContentLoaded', (event) => {
    const videoPlayer = document.getElementById('videoPlayer');
    const skipButton = document.getElementById('skipButton');

    const forwardSkipButton = document.getElementById('forwardSkipButton');
    const backwardSkipButton = document.getElementById('backwardSkipButton');

    skipButton.addEventListener('click', () => {
        videoPlayer.currentTime += 85;
    });

    forwardSkipButton.addEventListener('click', () => {
        videoPlayer.currentTime += 10;
    });

    backwardSkipButton.addEventListener('click', () => {
        videoPlayer.currentTime -= 10;
    });
});

let isLocked = false; // Variable to store the lock state

document.addEventListener('DOMContentLoaded', (event) => {
    const videoPlayer = document.getElementById('videoPlayer');
    const lockButton = document.getElementById('lockButton');
    const skipButton = document.getElementById('skipButton');
    const forwardSkipButton = document.getElementById('forwardSkipButton');
    const backwardSkipButton = document.getElementById('backwardSkipButton');

    lockButton.addEventListener('click', () => {
        isLocked = !isLocked; // Toggle the lock state
        if (isLocked) {
            lockButton.textContent = 'Unlock'; // Change the button text to 'Unlock'
            videoPlayer.controls = false; // Disable the video controls
            skipButton.style.display = 'none'; // Hide the '+85' button
            forwardSkipButton.style.display = 'none'; // Hide the '+10' button
            backwardSkipButton.style.display = 'none'; // Hide the '-10' button
        } else {
            lockButton.textContent = 'Lock'; // Change the button text to 'Lock'
            videoPlayer.controls = true; // Enable the video controls
            skipButton.style.display = 'block'; // Show the '+85' button
            forwardSkipButton.style.display = 'block'; // Show the '+10' button
            backwardSkipButton.style.display = 'block'; // Show the '-10' button
        }
    });
});

document.addEventListener('DOMContentLoaded', (event) => {
    const videoPlayer = document.getElementById('videoPlayer');
    const speedSelect = document.getElementById('speedSelect');

    speedSelect.addEventListener('change', () => {
        videoPlayer.playbackRate = speedSelect.value;
    });
});