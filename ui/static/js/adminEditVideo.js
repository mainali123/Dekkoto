console.log("adminEditVideo.js loaded")

// get data from localStorage
let vid_data = JSON.parse(localStorage.getItem('data'));
console.log(vid_data);

// Select the HTML elements
let titleElement = document.getElementById('title');
let descriptionElement = document.getElementById('description');
let genreElements = document.querySelectorAll('input[name="genre"]');
let typeElements = document.querySelectorAll('input[name="type"]');

// Set the value of the HTML elements
titleElement.value = vid_data.title;
descriptionElement.value = vid_data.description;

// Set the checked property of the genre checkboxes
genreElements.forEach((checkbox) => {
    if (vid_data.genres.includes(checkbox.value)) {
        checkbox.checked = true;
    }
});

// Set the checked property of the type radio buttons
typeElements.forEach((radio) => {
    if (vid_data.types === radio.value) {
        radio.checked = true;
    }
});