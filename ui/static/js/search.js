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
let resultsContainer = document.querySelector('.result-container-wrapper');

// Function to clear the current search results
function clearSearchResults() {
    while (resultsContainer.firstChild) {
        resultsContainer.removeChild(resultsContainer.firstChild);
    }
}

// Get the main element
const mainElement = document.querySelector('main');

// Function to create a new DOM element for each video result
function createVideoElement(video) {
    // Create a new div element for the results container
    const newDiv = document.createElement("div");
    newDiv.classList.add('results-container');

    // Create a new img element for the cover image
    const newImg = document.createElement("img");
    newImg.src = video.ThumbnailURL; // Use the ThumbnailURL property of the video object
    newImg.alt = video.Title + " cover image"; // Use the Title property of the video object
    newImg.classList.add('cover-img');

    // Add an event listener to the image
    newImg.addEventListener('click', function() {
        // save the video details in local storage
        localStorage.setItem('videoDetails', JSON.stringify(video));
        window.location.href = '/watchVideo'; // Redirect to the watch video page
    });

    // Create a new div element for the info
    const infoDiv = document.createElement("div");
    infoDiv.classList.add('info');

    // Create a new h2 element for the title
    const titleH2 = document.createElement("h2");
    titleH2.textContent = video.Title; // Use the Title property of the video object
    titleH2.classList.add('title');

    // Create a new p element for the alternate title
    const altTitleP = document.createElement("p");
    altTitleP.textContent = video.Title; // Use the Title property of the video object
    altTitleP.classList.add('alternate-title');

    // Create a new div element for the minor-infos
    const minorInfosDiv = document.createElement("div");
    minorInfosDiv.classList.add('minor-infos');

    // Create a new span element for the type
    const typeSpan = document.createElement("span");
    typeSpan.textContent = video.Genre; // Assuming the type is always "TV"
    typeSpan.classList.add('type');

    // Create a new span element for the status
    const statusSpan = document.createElement("span");
    // if video.Status is empty, set it to "Not watched"
    if (video.Status === "") {
        video.Status = "Not watched";
    }
    statusSpan.textContent = video.Status; // Assuming the status is always "completed"
    statusSpan.classList.add('status');

    // Append the type and status spans to the minor-infos div
    minorInfosDiv.appendChild(typeSpan);
    minorInfosDiv.appendChild(statusSpan);

    // Append the title, alternate title, and minor-infos div to the info div
    infoDiv.appendChild(titleH2);
    infoDiv.appendChild(altTitleP);
    infoDiv.appendChild(minorInfosDiv);

    // Append the cover image and info div to the results container div
    newDiv.appendChild(newImg);
    newDiv.appendChild(infoDiv);

    resultsContainer.appendChild(newDiv);

    // Append the new results container div to the main element

}

// Attach an event listener to the input field
searchInput.addEventListener('submit', function () {
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