let globalData

function editVideo(videoID) {
    console.log(videoID);
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
        console.log(globalData); // Print data in the console

        // Access the Videos array within the received data
        const videos = globalData.videos.Videos;

        // Reference to the table body where the data will be populated
        const tableBody = document.querySelector('#videoTable tbody');



        // Function to populate the table with video data
        function populateTable() {
            videos.forEach((video, index) => {
                const row = tableBody.insertRow();
                row.innerHTML = `
                <td>${index + 1}</td>
                <td>${video.Title}</td>
                <td>${video.Description}</td>
                <td>${new Date(video.UploadDate).toDateString()}</td>
                <td>${video.ViewsCount}</td>
                <td>${video.LikesCount}</td>
                <td>${video.Duration}</td>
                <td>${video.CategoryID}</td>
                <td>${video.GenreID}</td>
                <td><button class="btn btn-primary" onclick="editVideo(${video.VideoID})">Edit</button></td>
                <td><button class="btn btn-danger" onclick="deleteVideo(${video.VideoID})">Delete</button></td>
                `;
            });
        }

        // Call the function to populate the table
        populateTable();
    })
    .catch(error => {
        // Handle errors
        console.error('There was a problem with the fetch operation:', error);
    });