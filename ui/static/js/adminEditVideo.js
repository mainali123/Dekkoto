console.log("adminEditVideo.js loaded")

// get data from localStorage
let vid_data = JSON.parse(localStorage.getItem('data'));
console.log(vid_data);

// Select the HTML elements
let titleElement = document.getElementById('title');
let descriptionElement = document.getElementById('description');
let genreElements = document.querySelectorAll('input[name="genre"]');
let typeElements = document.querySelectorAll('input[name="type"]');
let formElement = document.getElementById('editVideoForm');

// Set the value of the HTML elements
titleElement.value = vid_data.title;
descriptionElement.value = vid_data.description;

// Split the genres and types strings into arrays
let genresArray = vid_data.categoryName.split(',');
let typesArray = [vid_data.genreName];

// Set the checked property of the genre checkboxes
genreElements.forEach((checkbox) => {
    if (genresArray.includes(checkbox.value)) {
        checkbox.checked = true;
    }
});

// Set the checked property of the type radio buttons
typeElements.forEach((radio) => {
    if (typesArray.includes(radio.value)) {
        radio.checked = true;
    }
});

// Add event listener for form submission
formElement.addEventListener('submit', function(event) {
    // Prevent default form submission
    event.preventDefault();

    // Get the edited values from the form fields
    let editedTitle = titleElement.value;
    let editedDescription = descriptionElement.value;
    let editedGenres = Array.from(genreElements).filter(checkbox => checkbox.checked).map(checkbox => checkbox.value);
    let editedType = Array.from(typeElements).find(radio => radio.checked).value;

    // Store the edited values in an object
    let editedData = {
        videoID: vid_data.videoID,
        title: editedTitle,
        description: editedDescription,
        // convert the array to comma separated string
        genre: editedGenres.join(','),
        type: editedType
    };

    console.log(editedData)
    // Send a POST request to the backend with the edited data
    fetch('/editVideoPost', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(editedData),
    })
    .then(response => {
        if (!response.ok) {
            console.log('Network response was not ok');
        } else {
            // console.log('Video details updated successfully')
            alert('Video details updated successfully');
            window.location.href = '/showVideos';
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });
});