console.log("loaded gadha");
document.querySelector('#submitBtn').addEventListener('click', function (event) {
    event.preventDefault();

    let title = document.getElementById('title').value;
    let description = document.getElementById('description').value;
    let category = document.querySelector('input[name="type"]:checked').value;
    let genreElements = document.querySelectorAll('input[name="genre"]:checked');
    let genre = Array.from(genreElements).map(function (element) {
        return element.value;
    }).join(',');

    let videoForm = document.getElementById('uploadVideoForm');
    let thumbnailForm = document.getElementById('uploadThumbnailForm');
    let bannerForm = document.getElementById('uploadBannerForm');

    videoForm.submit().preventDefault();
    thumbnailForm.submit().preventDefault();
    bannerForm.submit().preventDefault();

    let data = {
        title: title,
        description: description,
        genres: genre,
        types: category
    };

    fetch('/videoDetails', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data),
    })
        .then(response => response.json())
        .then(jsonResponse => {
            console.log(jsonResponse);
            alert(jsonResponse.message);
            // window.location.href = jsonResponse.redirect;
        })
        .catch(error => {
            console.error('There was a problem with the form submission:', error);
        });


    fetch('/uploadVideo', {
    method: "POST",
    headers: {'Content-Type': 'application/json'}
})
.then(response => {
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
    return response.json();
})
.then(data => {
    // Handle your JSON data here
    console.log(data);
    // For example, you can display a success message
    alert('Video uploaded successfully');
    // Or redirect the user to another page
    // window.location.href = '/some-page';
})
.catch(error => {
    console.error('There was a problem with the fetch operation:', error);
});
});
