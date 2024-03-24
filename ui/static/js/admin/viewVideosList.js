console.log('video list Loaded successfully')

// Helper function to truncate the description
function truncateDescription(description) {
    let words = description.split(' ');
    if (words.length > 3) {
        return words.slice(0, 3).join(' ') + '...';
    } else {
        return description;
    }
}

function convertToMinutesSeconds(timeInSeconds) {
    let minutes = Math.floor(timeInSeconds / 60);
    let seconds = Math.floor(timeInSeconds % 60);
    return `${minutes}:${seconds < 10 ? '0' + seconds : seconds}`;
}

// Fetch the data from the server
fetch('/showVideosPost', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
})
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(videoData => {
        let count=1;
        let videoList = document.getElementById('videoList');
        videoList.innerHTML = '';

        videoData.videos.forEach(video => {
            let videoItem = document.createElement('div');
            videoItem.className = 'videoList';

            // Use the truncateDescription function to truncate the video description
            let truncatedDescription = truncateDescription(video.Description);

            videoItem.innerHTML = `
            <div class="thumbnail-duration">
                <img src="${video.ThumbnailURL}" alt="${video.Title}">
                <p class="duration">${convertToMinutesSeconds(video.Duration)}</p>
            </div>
            <div class="title-desc">
                <h3>${truncateDescription(video.Title)}</h3>
                <p class="description">${truncatedDescription}</p>
            </div>
            <p class="views">${video.ViewsCount}</p>
            <p class="likes">${video.LikesCount}</p>
            <p class="dislikes">${video.DislikesCount}</p>
            <div class="category">
                <p>${video.CategoryID}</p>
            </div>
            <div class="uploaded-date">
                <p>${new Date(video.UploadDate).toDateString()}</p>
            </div>
            
            <div class="btn-container">
                <button class="btn btn-primary" id="editBtn${count}" metaData='${JSON.stringify(video)}'><i class="fa-solid fa-pen-to-square"></i></button>
                <button class="btn btn-primary" id="deleteBtn${count}" metaData='${JSON.stringify(video)}'><i class="fa-solid fa-trash"></i></button>
            </div>
        `;
            count++;
            videoList.appendChild(videoItem);
        });
    })
    .catch(error => {
        console.error('There was a problem with the fetch operation:', error);
    });