/*
console.log("hello world")



// videoDetails:"{"VideoID":39,"Title":"Nezuko","Description":"desc","URL":"./userUploadDatas/videos/20240109161922.406048_0b177be7-2f42-4dcf-b4de-2399eb452b29_323500696560_encoded.mp4","ThumbnailURL":"./userUploadDatas/thumbnails/20240109161922.406048_0b177be7-2f42-4dcf-b4de-2399eb452b29_323500696560.png","UploaderID":1,"UploadDate":"2024-01-09T00:00:00Z","ViewsCount":0,"LikesCount":0,"DislikesCount":0,"Duration":"00:00:26","CategoryID":6,"GenreID":1}"
// get the video details from local storage



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
*/


const container = document.querySelector('.container'),
    mainVideo = container.querySelector('video'),
    videoTimeline = container.querySelector('.video-timeline'),
    progressBar = container.querySelector('.progress-bar'),
    volumeBtn = container.querySelector('.volume i'),
    volumeSlider = container.querySelector('.left input');
(currentVidTime = container.querySelector('.current-time')),
    (videoDuration = container.querySelector('.video-duration')),
    (skipBackward = container.querySelector('.skip-backward i')),
    (skipForward = container.querySelector('.skip-forward i')),
    (playPauseBtn = container.querySelector('.play-pause i')),
    (speedBtn = container.querySelector('.playback-speed span')),
    (speedOptions = container.querySelector('.speed-options')),
    (pipBtn = container.querySelector('.pic-in-pic span')),
    (fullScreenBtn = container.querySelector('.fullscreen i'));



let timer;

const hideControls = () => {
    if (mainVideo.paused) return;
    timer = setTimeout(() => {
        container.classList.remove('show-controls');
    }, 3000);
};
hideControls();

container.addEventListener('mousemove', () => {
    container.classList.add('show-controls');
    clearTimeout(timer);
    hideControls();
});

const formatTime = (time) => {
    let seconds = Math.floor(time % 60),
        minutes = Math.floor(time / 60) % 60,
        hours = Math.floor(time / 3600);

    seconds = seconds < 10 ? `0${seconds}` : seconds;
    minutes = minutes < 10 ? `0${minutes}` : minutes;
    hours = hours < 10 ? `0${hours}` : hours;

    if (hours == 0) {
        return `${minutes}:${seconds}`;
    }
    return `${hours}:${minutes}:${seconds}`;
};

videoTimeline.addEventListener('mousemove', (e) => {
    let timelineWidth = videoTimeline.clientWidth;
    let offsetX = e.offsetX;
    let percent = Math.floor((offsetX / timelineWidth) * mainVideo.duration);
    const progressTime = videoTimeline.querySelector('span');
    offsetX =
        offsetX < 20
            ? 20
            : offsetX > timelineWidth - 20
                ? timelineWidth - 20
                : offsetX;
    progressTime.style.left = `${offsetX}px`;
    progressTime.innerText = formatTime(percent);
});

videoTimeline.addEventListener('click', (e) => {
    let timelineWidth = videoTimeline.clientWidth;
    mainVideo.currentTime = (e.offsetX / timelineWidth) * mainVideo.duration;
});

mainVideo.addEventListener('timeupdate', (e) => {
    let { currentTime, duration } = e.target;
    let percent = (currentTime / duration) * 100;
    progressBar.style.width = `${percent}%`;
    currentVidTime.innerText = formatTime(currentTime);
});

mainVideo.addEventListener('loadeddata', () => {
    videoDuration.innerText = formatTime(mainVideo.duration);
});

const draggableProgressBar = (e) => {
    let timelineWidth = videoTimeline.clientWidth;
    progressBar.style.width = `${e.offsetX}px`;
    mainVideo.currentTime = (e.offsetX / timelineWidth) * mainVideo.duration;
    currentVidTime.innerText = formatTime(mainVideo.currentTime);
};

volumeBtn.addEventListener('click', () => {
    if (!volumeBtn.classList.contains('fa-volume-high')) {
        mainVideo.volume = 0.5;
        volumeBtn.classList.replace('fa-volume-xmark', 'fa-volume-high');
    } else {
        mainVideo.volume = 0.0;
        volumeBtn.classList.replace('fa-volume-high', 'fa-volume-xmark');
    }
    volumeSlider.value = mainVideo.volume;
});

volumeSlider.addEventListener('input', (e) => {
    mainVideo.volume = e.target.value;
    if (e.target.value == 0) {
        return volumeBtn.classList.replace('fa-volume-high', 'fa-volume-xmark');
    }
    volumeBtn.classList.replace('fa-volume-xmark', 'fa-volume-high');
});

speedOptions.querySelectorAll('li').forEach((option) => {
    option.addEventListener('click', () => {
        mainVideo.playbackRate = option.dataset.speed;
        speedOptions.querySelector('.active').classList.remove('active');
        option.classList.add('active');
    });
});

document.addEventListener('click', (e) => {
    if (
        e.target.tagName !== 'SPAN' ||
        e.target.className !== 'material-symbols-rounded'
    ) {
        speedOptions.classList.remove('show');
    }
});

fullScreenBtn.addEventListener('click', () => {
    container.classList.toggle('fullscreen');
    if (document.fullscreenElement) {
        fullScreenBtn.classList.replace('fa-compress', 'fa-expand');
        return document.exitFullscreen();
    }
    fullScreenBtn.classList.replace('fa-expand', 'fa-compress');
    container.requestFullscreen();
});

speedBtn.addEventListener('click', () => speedOptions.classList.toggle('show'));
// pipBtn.addEventListener('click', () => mainVideo.requestPictureInPicture());
skipBackward.addEventListener('click', () => (mainVideo.currentTime -= 5));
skipForward.addEventListener('click', () => (mainVideo.currentTime += 5));
mainVideo.addEventListener('play', () =>
    playPauseBtn.classList.replace('fa-play', 'fa-pause')
);
mainVideo.addEventListener('pause', () =>
    playPauseBtn.classList.replace('fa-pause', 'fa-play')
);
playPauseBtn.addEventListener('click', () =>
    mainVideo.paused ? mainVideo.play() : mainVideo.pause()
);
videoTimeline.addEventListener('mousedown', () =>
    videoTimeline.addEventListener('mousemove', draggableProgressBar)
);
document.addEventListener('mouseup', () =>
    videoTimeline.removeEventListener('mousemove', draggableProgressBar)
);


// My CODE FROM HERE
let videoDetails = localStorage.getItem('videoDetails');
// convert to JSON
videoDetails = JSON.parse(videoDetails);
console.log(videoDetails)

// Video action changed
let videoID = videoDetails.VideoID;

// Get the 'statusSelect' dropdown element
let statusSelect = document.getElementById('statusSelect');

// Send a POST request to the '/videoAction' endpoint
fetch('/videoAction', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({id: videoID}),
})
    .then(response => response.json())
    .then(data => {
        // Use the returned string value to set the selected option in the 'statusSelect' dropdown
        statusSelect.value = data.action;
    })
    .catch((error) => {
        console.error('Error:', error);
    });

// Update the value of changed action
document.addEventListener('DOMContentLoaded', (event) => {
    const statusSelect = document.getElementById('statusSelect');

    statusSelect.addEventListener('change', () => {
        let selectedStatus = statusSelect.value;
        let videoID = videoDetails.VideoID; // Assuming videoDetails is globally available

        fetch('/videoActionChanged', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({VideoID: videoID, Action: selectedStatus}),
        })
            .then(response => response.json())
            .then(data => {
                console.log(data);
            })
            .catch((error) => {
                console.error('Error:', error);
            });
    });
});

// send a post request to the server only once
fetch('/watchVideoPost', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json', // Specify content type as JSON
    },
    body: JSON.stringify(videoDetails),
}).then(response => {
    // Check if the response is valid
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
    console.error('Error:', error);
});


// Select the video player and source elements
let videoPlayer = document.getElementById('video-player');

videoPlayer.src = videoDetails.URL;

console.log(videoPlayer.src)

// Load the video
videoPlayer.load();

// Try to play the video
videoPlayer.play().catch(error => {
    // Auto-play was prevented
    // Show a UI element to let the user manually start the playback
    console.log("Auto-play was prevented by the browser. Please start the video manually.");
});

