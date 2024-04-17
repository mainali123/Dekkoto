console.log('profile.js loaded');

function checkAdminAccess() {
    // Fetch user details from the server
    fetch('/checkAdminAccess', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => response.json())
        .then(data => {
            // Check if the request was successful
            if (data.success) {
                console.log(data)
                // Check if the user is an admin
                if (!data.adminAccess) {
                    // hide the dropdown-item class
                    document.querySelector('.admin-dropdownshow').style.display = 'none';
                }
            } else {
                console.error('Failed to fetch user details:', data.message);
            }
        })
        .catch(error => console.error('Error:', error));
}
checkAdminAccess();

function displayQuote() {
    // Fetch quotes from the server
    fetch('/quotes', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            // Check if the request was successful
            // Generate a random index

            // Select a quote
            const quote = data.quote;
            const author = data.author;

            // combine quote and author into one string with (quote - author)
            const quoteAuthor = quote + ' - ' + author;

            // Select the 'user-bio' element
            const userBioElement = document.querySelector('.user-bio');

            // Update the 'user-bio' element with the selected quote
            userBioElement.textContent = quoteAuthor;

        })
        .catch(error => console.error('Error:', error));
}

displayQuote();

// Fetch user details from the server
fetch('/userDetails', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    }
})
    .then(response => response.json())
    .then(data => {
        // Check if the request was successful
        if (data.success) {
            // Select the 'user-name' element
            const userNameElement = document.querySelector('.user-name');

            // Update the 'user-name' element with the received user name
            userNameElement.textContent = data.userName;
        } else {
            console.error('Failed to fetch user details:', data.message);
        }
    })
    .catch(error => console.error('Error:', error));


// Fetch user details from the server
fetch('/videoDatas', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    }
})
    .then(response => response.json())
    .then(data => {
        console.log(data);
        let videoDatas = data.videos;
        document.querySelector('.recommen').textContent = videoDatas.Recommends;
        document.querySelector('.watch').textContent = videoDatas.Watching;
        document.querySelector('.compli').textContent = videoDatas.Completed;
        document.querySelector('.hold').textContent = videoDatas.OnHold;
        document.querySelector('.consider').textContent = videoDatas.Considering;
        document.querySelector('.drop').textContent = videoDatas.Dropped;
    })
    .catch(error => console.error(error));

// Fetch watching videos from the server
fetch('/watchingVideos', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    }
})
    .then(response => response.json())
    .then(data => {
        console.log(data);
        const videos = data.videos;
        const videoContainer = document.querySelector('.watching'); // Select the HTML element where you want to display the videos

        // Check if there are videos
        if (videos.length > 0) {
            videos.forEach(video => {
                // Create HTML elements for each video
                const imageDiv = document.createElement('div');
                imageDiv.classList.add('imageDiv')
                const img = document.createElement('img');
                img.src = '../../' + video.ThumbnailURL; // Adjust the path relative to the HTML file
                img.dataset.videoDetails = JSON.stringify(video); // Store the video details in the img element

                // Add an event listener to the img element
                img.addEventListener('click', function () {
                    // Save the video details in local storage
                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                    window.location.href = '/watchVideo'; // Redirect to the watch video page
                });

                // Append the created HTML elements to the video container
                videoContainer.appendChild(imageDiv);
                imageDiv.appendChild(img);
            });
        }
    })
    .catch(error => console.error(error));

// On hold videos
fetch('/onHoldVideos', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
})
    .then(response => response.json())
    .then(data => {
        console.log(data);
        const videos = data.videos;
        const videoContainer = document.querySelector('.on-hold'); // Select the HTML element where you want to display the videos

        // Check if there are videos
        if (videos.length > 0) {
            videos.forEach(video => {
                // Create HTML elements for each video
                const imageDiv = document.createElement('imageDiv');
                imageDiv.classList.add('imageDiv');
                const img = document.createElement('img');
                img.src = '../../' + video.ThumbnailURL; // Adjust the path relative to the HTML file
                img.dataset.videoDetails = JSON.stringify(video); // Store the video details in the img element

                // Add an event listener to the img element
                img.addEventListener('click', function () {
                    // Save the video details in local storage
                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                    window.location.href = '/watchVideo'; // Redirect to the watch video page
                });

                // Append the created HTML elements to the video container
                videoContainer.appendChild(imageDiv);
                imageDiv.appendChild(img);
            });
        }
    })
    .catch(error => console.error(error));

// Considering videos
fetch('/consideringVideos', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
})
    .then(response => response.json())
    .then(data => {
        console.log(data);
        const videos = data.videos;
        const videoContainer = document.querySelector('.considered-anime'); // Select the HTML element where you want to display the videos

        // Check if there are videos
        if (videos.length > 0) {
            videos.forEach(video => {
                // Create HTML elements for each video
                const img = document.createElement('img');
                img.src = '../../' + video.ThumbnailURL; // Adjust the path relative to the HTML file
                img.dataset.videoDetails = JSON.stringify(video); // Store the video details in the img element

                // Add an event listener to the img element
                img.addEventListener('click', function () {
                    // Save the video details in local storage
                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                    window.location.href = '/watchVideo'; // Redirect to the watch video page
                });

                // Append the created HTML elements to the video container
                videoContainer.appendChild(img);
            });
        }
    })
    .catch(error => console.error(error));

// recentlyCompletedVideos
fetch('/recentlyCompletedVideos', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
})
    .then(response => response.json())
    .then(data => {
        console.log(data);
        const videos = data.videos;
        const videoContainer = document.querySelector('.recently-completed');

        if (videos.length > 0) {
            videos.forEach(video => {
                const finishedDiv = document.createElement('div');
                finishedDiv.className = 'finished';

                const img = document.createElement('img');
                img.src = '../../' + video.ThumbnailURL;
                img.alt = video.Title;
                finishedDiv.appendChild(img);

                const detailsDiv = document.createElement('div');
                detailsDiv.className = 'details';
                finishedDiv.appendChild(detailsDiv);

                const titleP = document.createElement('p');
                titleP.className = 'title-anime';
                titleP.textContent = video.Title;
                detailsDiv.appendChild(titleP);

                const timeSetDiv = document.createElement('div');
                timeSetDiv.className = 'time-set';
                detailsDiv.appendChild(timeSetDiv);

                const finishDateSpan = document.createElement('span');
                finishDateSpan.className = 'finish-date';
                const completedDate = new Date(video.CompletedDate);
                finishDateSpan.textContent = completedDate.toLocaleDateString('en-US', {
                    month: 'short',
                    day: 'numeric',
                    year: 'numeric'
                });
                timeSetDiv.appendChild(finishDateSpan);

                finishedDiv.addEventListener('click', function () {
                    localStorage.setItem('videoDetails', JSON.stringify(video));
                    window.location.href = '/watchVideo';
                });

                videoContainer.appendChild(finishedDiv);
            });
        }
    })
    .catch(error => console.error(error));