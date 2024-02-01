console.log("videoList.js loaded")

document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("recommended_backend").addEventListener("click", recommendedVideo);
    document.getElementById("watching_backend").addEventListener("click", watchingVideo);
    document.getElementById("completed_backend").addEventListener("click", completedVideo);
    document.getElementById("on_hold_backend").addEventListener("click", onHoldVideo);
    document.getElementById("considering_backend").addEventListener("click", consideringVideo);
    document.getElementById("dropped_backend").addEventListener("click", droppedVideo);
});


function recommendedVideo() {
    console.log("recommended")
//      Endpoint = "/recommendedVideoList" via POST

    fetch('/videoslist')
        .then(response => response.text())
        .then(html => {
            // Parse the HTML string into a new DOM Document
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');

            // Replace the current document with the new one
            document.replaceChild(document.importNode(doc.documentElement, true), document.documentElement);

            // Now that the '/videoslist' page is loaded, proceed with the fetch operation
            fetch('/recommendedVideoList', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
                .then(response => response.json())
                .then(data => {
                    // Check if the request was successful
                    if (data.success) {
                        // Get the parent element
                        const parentElement = document.querySelector('.card-video-continue-watching');

                        // Check if the parent element exists
                        if (parentElement) {
                            // Iterate over the videos array
                            data.videos.forEach(video => {
                                // Create a new div element
                                const div = document.createElement('div');
                                div.className = 'card';

                                // Create an img element
                                const img = document.createElement('img');
                                img.className = 'card-img-top';
                                img.src = video.ThumbnailURL;
                                img.alt = video.Title;
                                img.dataset.videoDetails = JSON.stringify(video);
                                img.addEventListener('click', function () {
                                    // save the video details in local storage
                                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                                    window.location.href = '/watchVideo'; // Redirect to the watch video page
                                    // console.log(this.dataset.videoDetails); // Log the video details when the image is clicked
                                });

                                // Create a p element for the title
                                const pTitle = document.createElement('p');
                                pTitle.className = 'card-text';
                                pTitle.textContent = video.Title;

                                // Append the img and p elements to the div
                                div.appendChild(img);
                                div.appendChild(pTitle);

                                // Append the div to the parent element
                                parentElement.appendChild(div);
                            });
                        } else {
                            console.error('Parent element not found');
                        }
                    } else {
                        console.error('Failed to fetch videos:', data.message);
                    }
                })
                .catch(error => console.error('Error:', error));
        })
        .catch(error => console.error('Error:', error));
}

function watchingVideo() {
    console.log("watching");
    fetch('/videoslist')
        .then(response => response.text())
        .then(html => {
            // Parse the HTML string into a new DOM Document
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');

            // Replace the current document with the new one
            document.replaceChild(document.importNode(doc.documentElement, true), document.documentElement);

            // Now that the '/videoslist' page is loaded, proceed with the fetch operation
            fetch('/watchingVideoList', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => response.json())
            .then(data => {
                // Check if the request was successful
                if (data.success) {
                    // Get the parent element
                    const parentElement = document.querySelector('.card-video-continue-watching');

                    // Check if the parent element exists
                    if (parentElement) {
                        // Iterate over the videos array
                        data.videos.forEach(video => {
                            // Create a new div element
                            const div = document.createElement('div');
                            div.className = 'card';

                            // Create an img element
                            const img = document.createElement('img');
                            img.className = 'card-img-top';
                            img.src = video.ThumbnailURL;
                            img.alt = video.Title;
                            img.dataset.videoDetails = JSON.stringify(video);
                            img.addEventListener('click', function () {
                                // save the video details in local storage
                                localStorage.setItem('videoDetails', this.dataset.videoDetails);
                                window.location.href = '/watchVideo'; // Redirect to the watch video page
                                // console.log(this.dataset.videoDetails); // Log the video details when the image is clicked
                            });

                            // Create a p element for the title
                            const pTitle = document.createElement('p');
                            pTitle.className = 'card-text';
                            pTitle.textContent = video.Title;

                            // Append the img and p elements to the div
                            div.appendChild(img);
                            div.appendChild(pTitle);

                            // Append the div to the parent element
                            parentElement.appendChild(div);
                        });
                    } else {
                        console.error('Parent element not found');
                    }
                } else {
                    console.error('Failed to fetch videos:', data.message);
                }
            })
            .catch(error => console.error('Error:', error));
        })
        .catch(error => console.error('Error:', error));
}

function completedVideo() {
    // Endpoint = "/completedVideoList" via POST
    console.log("completed")

    fetch('/videoslist')
        .then(response => response.text())
        .then(html => {
            // Parse the HTML string into a new DOM Document
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');

            // Replace the current document with the new one
            document.replaceChild(document.importNode(doc.documentElement, true), document.documentElement);

            // Now that the '/videoslist' page is loaded, proceed with the fetch operation
            fetch('/completedVideoList', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
                .then(response => response.json())
                .then(data => {
                    // Check if the request was successful
                    if (data.success) {
                        // Get the parent element
                        const parentElement = document.querySelector('.card-video-continue-watching');

                        // Check if the parent element exists
                        if (parentElement) {
                            // Iterate over the videos array
                            data.videos.forEach(video => {
                                // Create a new div element
                                const div = document.createElement('div');
                                div.className = 'card';

                                // Create an img element
                                const img = document.createElement('img');
                                img.className = 'card-img-top';
                                img.src = video.ThumbnailURL;
                                img.alt = video.Title;
                                img.dataset.videoDetails = JSON.stringify(video);
                                img.addEventListener('click', function () {
                                    // save the video details in local storage
                                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                                    window.location.href = '/watchVideo'; // Redirect to the watch video page
                                    // console.log(this.dataset.videoDetails); // Log the video details when the image is clicked
                                });

                                // Create a p element for the title
                                const pTitle = document.createElement('p');
                                pTitle.className = 'card-text';
                                pTitle.textContent = video.Title;

                                // Append the img and p elements to the div
                                div.appendChild(img);
                                div.appendChild(pTitle);

                                // Append the div to the parent element
                                parentElement.appendChild(div);
                            });
                        } else {
                            console.error('Parent element not found');
                        }
                    } else {
                        console.error('Failed to fetch videos:', data.message);
                    }
                })
                .catch(error => console.error('Error:', error));
        })
        .catch(error => console.error('Error:', error));
}

function onHoldVideo() {
    // Endpoint = "/onHoldVideoList" via POST
    console.log("on hold")

    fetch('/videoslist')
        .then(response => response.text())
        .then(html => {
            // Parse the HTML string into a new DOM Document
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');

            // Replace the current document with the new one
            document.replaceChild(document.importNode(doc.documentElement, true), document.documentElement);

            // Now that the '/videoslist' page is loaded, proceed with the fetch operation
            fetch('/onHoldVideoList', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
                .then(response => response.json())
                .then(data => {
                    // Check if the request was successful
                    if (data.success) {
                        // Get the parent element
                        const parentElement = document.querySelector('.card-video-continue-watching');

                        // Check if the parent element exists
                        if (parentElement) {
                            // Iterate over the videos array
                            data.videos.forEach(video => {
                                // Create a new div element
                                const div = document.createElement('div');
                                div.className = 'card';

                                // Create an img element
                                const img = document.createElement('img');
                                img.className = 'card-img-top';
                                img.src = video.ThumbnailURL;
                                img.alt = video.Title;
                                img.dataset.videoDetails = JSON.stringify(video);
                                img.addEventListener('click', function () {
                                    // save the video details in local storage
                                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                                    window.location.href = '/watchVideo'; // Redirect to the watch video page
                                    // console.log(this.dataset.videoDetails); // Log the video details when the image is clicked
                                });

                                // Create a p element for the title
                                const pTitle = document.createElement('p');
                                pTitle.className = 'card-text';
                                pTitle.textContent = video.Title;

                                // Append the img and p elements to the div
                                div.appendChild(img);
                                div.appendChild(pTitle);

                                // Append the div to the parent element
                                parentElement.appendChild(div);
                            });
                        } else {
                            console.error('Parent element not found');
                        }
                    } else {
                        console.error('Failed to fetch videos:', data.message);
                    }
                })
                .catch(error => console.error('Error:', error));
        })
        .catch(error => console.error('Error:', error));
}

function consideringVideo() {
    // Endpoint = "/consideringVideoList" via POST
    console.log("considering")

    fetch('/videoslist')
        .then(response => response.text())
        .then(html => {
            // Parse the HTML string into a new DOM Document
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');

            // Replace the current document with the new one
            document.replaceChild(document.importNode(doc.documentElement, true), document.documentElement);

            // Now that the '/videoslist' page is loaded, proceed with the fetch operation
            fetch('/consideringVideoList', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
                .then(response => response.json())
                .then(data => {
                    // Check if the request was successful
                    if (data.success) {
                        // Get the parent element
                        const parentElement = document.querySelector('.card-video-continue-watching');

                        // Check if the parent element exists
                        if (parentElement) {
                            // Iterate over the videos array
                            data.videos.forEach(video => {
                                // Create a new div element
                                const div = document.createElement('div');
                                div.className = 'card';

                                // Create an img element
                                const img = document.createElement('img');
                                img.className = 'card-img-top';
                                img.src = video.ThumbnailURL;
                                img.alt = video.Title;
                                img.dataset.videoDetails = JSON.stringify(video);
                                img.addEventListener('click', function () {
                                    // save the video details in local storage
                                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                                    window.location.href = '/watchVideo'; // Redirect to the watch video page
                                    // console.log(this.dataset.videoDetails); // Log the video details when the image is clicked
                                });

                                // Create a p element for the title
                                const pTitle = document.createElement('p');
                                pTitle.className = 'card-text';
                                pTitle.textContent = video.Title;

                                // Append the img and p elements to the div
                                div.appendChild(img);
                                div.appendChild(pTitle);

                                // Append the div to the parent element
                                parentElement.appendChild(div);
                            });
                        } else {
                            console.error('Parent element not found');
                        }
                    } else {
                        console.error('Failed to fetch videos:', data.message);
                    }
                })
                .catch(error => console.error('Error:', error));
        })
        .catch(error => console.error('Error:', error));
}

function droppedVideo() {
    // Endpoint = "/droppedVideoList" via POST
    console.log("dropped")

    fetch('/videoslist')
        .then(response => response.text())
        .then(html => {
            // Parse the HTML string into a new DOM Document
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');

            // Replace the current document with the new one
            document.replaceChild(document.importNode(doc.documentElement, true), document.documentElement);

            // Now that the '/videoslist' page is loaded, proceed with the fetch operation
            fetch('/droppedVideoList', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
                .then(response => response.json())
                .then(data => {
                    // Check if the request was successful
                    if (data.success) {
                        // Get the parent element
                        const parentElement = document.querySelector('.card-video-continue-watching');

                        // Check if the parent element exists
                        if (parentElement) {
                            // Iterate over the videos array
                            data.videos.forEach(video => {
                                // Create a new div element
                                const div = document.createElement('div');
                                div.className = 'card';

                                // Create an img element
                                const img = document.createElement('img');
                                img.className = 'card-img-top';
                                img.src = video.ThumbnailURL;
                                img.alt = video.Title;
                                img.dataset.videoDetails = JSON.stringify(video);
                                img.addEventListener('click', function () {
                                    // save the video details in local storage
                                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                                    window.location.href = '/watchVideo'; // Redirect to the watch video page
                                    // console.log(this.dataset.videoDetails); // Log the video details when the image is clicked
                                });

                                // Create a p element for the title
                                const pTitle = document.createElement('p');
                                pTitle.className = 'card-text';
                                pTitle.textContent = video.Title;

                                // Append the img and p elements to the div
                                div.appendChild(img);
                                div.appendChild(pTitle);

                                // Append the div to the parent element
                                parentElement.appendChild(div);
                            });
                        } else {
                            console.error('Parent element not found');
                        }
                    } else {
                        console.error('Failed to fetch videos:', data.message);
                    }
                })
                .catch(error => console.error('Error:', error));
        })
        .catch(error => console.error('Error:', error));
}
