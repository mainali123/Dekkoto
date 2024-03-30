console.log("search.js loaded")

let value = false;


function onetimeFetch() {
    fetch('/searchData', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({search: ""}),
    })
        .then(response => response.json())
        .then(data => {
            // Clear the current search results
            console.log(data)
            // Add the new search results to the DOM
            clearSearchResults();
            data.videos.forEach(video => {
                createVideoElement(video);
            });
        })
        .catch((error) => {
            console.error('Error:', error);
        });
}

if (!value) {
    onetimeFetch();
    value = true;
}

// Get the search input field
const searchInput = document.querySelector('.search-bar');

// Get the results container
let resultsContainer = document.querySelector('.card-display');

// Function to clear the current search results
function clearSearchResults() {
    while (resultsContainer.firstChild) {
        resultsContainer.removeChild(resultsContainer.firstChild);
    }
}


// Function to create a new DOM element for each video result
function createVideoElement(video) {
    // Encode video details to Base64
    const videoDetailsBase64 = btoa(JSON.stringify(video));

    let cardOfVideo = `
        <div class="card">
            <img class="card-img-top" src="${video.ThumbnailURL}" data-video-details="${videoDetailsBase64}"/>
            <p class="card-text">${video.Title}</p>
            <p class="rating status">${video.Status}</p>
            <p class=" type">${video.Genre}</p>
            <div class="card-body"></div>
        </div>`;
    resultsContainer.innerHTML += cardOfVideo;
}


// Attach an event listener to the input field
searchInput.addEventListener('keyup', function () {
    searchValue = document.querySelector('.search-input').value;

    event.preventDefault();
    console.log("I am here")
    // Get the value of the input field
    console.log(searchValue);

    // Send an AJAX request with the input value
    fetch('/searchData', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({search: searchValue}),
    })
        .then(response => response.json())
        .then(data => {
            // Clear the current search results
            clearSearchResults();
            console.log(data)
            // Add the new search results to the DOM
            data.videos.forEach(video => {
                createVideoElement(video);
            });
        })
        .catch((error) => {
            console.error('Error:', error);
        });
});


// Get the search input field
const searchData = document.querySelector('.search-input');

// Attach an event listener to the input field
searchData.addEventListener('input', function () {
    console.log("I am here")
    // Get the value of the input field
    const searchValue = searchData.value;
    let autoComplete = document.querySelector('#autoComplete_list_1');
    if (searchValue === "") {
        autoComplete.style.display = "none";
    } else {
        autoComplete.style.display = "flex";
    }

});

// Add event listener to the video cards
resultsContainer.addEventListener('click', function (event) {
    let videoDetailsBase64 = event.target.getAttribute('data-video-details');
    if (videoDetailsBase64) {
        // Decode Base64 to JSON
        const videoDetails = JSON.parse(atob(videoDetailsBase64));
        console.log(videoDetails);
        // Store the selected video's details in local storage
        localStorage.setItem('videoDetails', JSON.stringify(videoDetails));
        // Redirect to the watch video page
        window.location.href = '/watchVideo';
    }
});
