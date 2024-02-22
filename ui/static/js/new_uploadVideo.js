console.log("loaded new upload");

// Get form element
const form = document.querySelector('form');
const title = document.getElementById('title');
const description = document.getElementById('description');
const genreSelect = document.querySelectorAll('input[name="genre"]');
const typeSelect = document.querySelector('input[name="type"]:checked');

// Add event listener to form submit
form.addEventListener('submit', function(event) {

    console.log("What is happening here")
    event.preventDefault();

    // Get form data
    const formData = new FormData(form);
    let selectedGenres = Array.from(genreSelect).filter(checkbox => checkbox.checked).map(checkbox => checkbox.value);
    let selectedType = typeSelect.value;

    // Prepare video description data
    const vidDesc = {
        title: title.value,
        description: description.value,
        genre: selectedGenres,
        type: selectedType
    };

    // Send video file
    fetch('/uploadVideo', {
        method: 'POST',
        body: formData.get('video')
    })
        .then()

    // Send thumbnail file to new endpoint
    fetch('/uploadThumbnail', {
        method: 'POST',
        body: formData.get('thumbnail')
    });

    // Send banner file to new endpoint
    fetch('/uploadBanner', {
        method: 'POST',
        body: formData.get('banner')
    });

    // Send video description data
    fetch('/videoDetails', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(vidDesc),
    })
    .then(function (response) {
        return response.json();
    })
    .then(function (jsonResponse) {
        // Display response back to user
        console.log(jsonResponse);
        alert(jsonResponse.message);
    })
    .catch(function (error) {
        console.log(error);
    });
});