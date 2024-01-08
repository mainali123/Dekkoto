let globalData

let data = {
    videoID: null,
    title: null,
    description: null,
    categoryID: null,
    genreID: null,
    categoryName: null,
    genreName: null
}

function editVideo(videoID, title, description, categoryID, genreID, categoryName, genreName) {
    data.videoID = videoID;
    data.title = title;
    data.description = description;
    data.categoryID = categoryID;
    data.genreID = genreID;
    data.categoryName = categoryName;
    data.genreName = genreName;
    console.log(data); // Check updated data before redirecting

    // send data to adminEditVideo.js
    localStorage.setItem('data', JSON.stringify(data));

    // Redirect after updating data
    window.location.href = '/editVideo'; // Change the URL to your desired location
}

function deleteVideo(videoID) {
    console.log(videoID);
    fetch('/deleteVideo', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json', // Specify content type as JSON
        },
        body: JSON.stringify({videoID: videoID}),
    })
        .then(response => {
            // Check if the response is valid
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            alert('Video deleted successfully');
            window.location.href = '/showVideos';
        })
        .then(data => {
        })
        .catch(error => {
            // Handle errors
            console.error('There was a problem with the fetch operation:', error);
            // If error is TypeError: tableBody is null then don't show error in console
            if (error.name !== 'TypeError') {
                alert('There was a problem with the fetch operation:' + error);
            }
        });
}

fetch('/showVideosPost', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json', // Specify content type as JSON
    },
})
    .then(response => {
        // Check if the response is valid
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        // Parse the JSON data from the response
        return response.json();
    })
    .then(data => {
        // Data received from the backend
        globalData = data;
        // console.log(globalData); // Print data in the console

        // Access the Videos array within the received data
        const videos = globalData.videos.Videos;

        // Reference to the table body where the data will be populated
        const tableBody = document.querySelector('#videoTable tbody');

        let categoryName = "";
        let genreName;

        // Function to populate the table with video data
        async function populateTable() {
            for (const video of videos) {
                const index = videos.indexOf(video);
                const row = tableBody.insertRow();

                // Get the category name
                function getCatName(categoryID) {
                    return fetch('/showCategoriesName', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json', // Specify content type as JSON
                        },
                        body: JSON.stringify({categoryID: categoryID}),
                    })
                        .then(response => {
                            // Check if the response is valid
                            if (!response.ok) {
                                throw new Error('Network response was not ok');
                            }
                            return response.json();
                        })
                        .then(data => {
                            // Data received from the backend
                            categoryName = data.categoryName
                            return categoryName;
                        })
                        .catch(error => {
                            // Handle errors
                            // console.error('There was a problem with the fetch operation:', error);
                            // If error is TypeError: tableBody is null then don't show error in console
                            if (error.name !== 'TypeError') {
                                console.log('There was a problem with the fetch operation:' + error);
                            }
                        });
                }

                categoryName = await getCatName(video.CategoryID);

                // Get the genre name
                function getGenreName(genreID) {
                    return fetch('/showGenresName', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json', // Specify content type as JSON
                        },
                        body: JSON.stringify({genreID: genreID}),
                    })
                        .then(response => {
                            // Check if the response is valid
                            if (!response.ok) {
                                throw new Error('Network response was not ok');
                            }
                            return response.json();
                        })
                        .then(data => {
                            // Data received from the backend
                            genreName = data.genreName
                            return genreName;
                        })
                        .catch(error => {
                            // Handle errors
                            // console.error('There was a problem with the fetch operation:', error);
                            // If error is TypeError: tableBody is null then don't show error in console
                            if (error.name !== 'TypeError') {
                                console.log('There was a problem with the fetch operation:' + error);
                            }
                        });
                }

                genreName = await getGenreName(video.GenreID);


                row.innerHTML = `
                <td>${index + 1}</td>
                <td>${video.Title}</td>
                <td>${video.Description}</td>
                <td>${new Date(video.UploadDate).toDateString()}</td>
                <td>${video.ViewsCount}</td>
                <td>${video.LikesCount}</td>
                <td>${video.Duration}</td>
                <td>${categoryName}</td>
                <td>${genreName}</td>
                <td><button class="btn btn-primary" onclick="editVideo('${video.VideoID}', '${video.Title}', '${video.Description}', '${video.CategoryID}', '${video.GenreID}', '${categoryName}', '${genreName}')">Edit</button></td>
                <td><button class="btn btn-danger" onclick="deleteVideo(${video.VideoID})">Delete</button></td>
                `;
            }
        }

        // Call the function to populate the table
        populateTable();
    })
    .catch(error => {
        // Handle errors
        // console.error('There was a problem with the fetch operation:', error);
        // If error is TypeError: tableBody is null then don't show error in console
        if (error.name !== 'TypeError') {
            alert('There was a problem with the fetch operation:', error);
        }
    });
