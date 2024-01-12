console.log("Homepage.js loaded")

fetch('/showVideosPost', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
})
.then(response => {
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
    return response.json();
})
.then(data => {
    console.log(data)
    const videos = data.videos.Videos;
    const videoContainer = document.querySelector('.card-video'); // Select the HTML element where you want to display the videos

    // Check if there are videos
    if (videos.length > 0) {
        videos.forEach(video => {
            // Create HTML elements for each video
            const card = document.createElement('div');
            card.className = 'card';
            card.style.width = '10rem';

            const img = document.createElement('img');
            img.className = 'card-img-top';
            img.src = '../../' + video.ThumbnailURL; // Adjust the path relative to the HTML file
            img.dataset.videoDetails = JSON.stringify(video); // Store the video details in the img element
            img.addEventListener('click', function() {
                // save the video details in local storage
                localStorage.setItem('videoDetails', this.dataset.videoDetails);
                window.location.href = '/watchVideo'; // Redirect to the watch video page
                // console.log(this.dataset.videoDetails); // Log the video details when the image is clicked
            });

            const cardText = document.createElement('p');
            cardText.className = 'card-text';
            cardText.textContent = video.Title;

            // Append the created HTML elements to the card
            card.appendChild(img);
            card.appendChild(cardText); // Append the title to the card

            const cardBody = document.createElement('div');
            cardBody.className = 'card-body';

            // Create rating and colored background elements
            const rating = document.createElement('p');
            rating.className = 'rating'; // Add the 'rating' class to the rating element
            rating.textContent = 0; // Assuming the video object has a Rating property

            const coloredBackground = document.createElement('div');
            coloredBackground.className = 'colored-background';

            // Append the created HTML elements to the card
            card.appendChild(coloredBackground);
            card.appendChild(rating);
            card.appendChild(cardBody);

            // Append the card to the video container
            videoContainer.appendChild(card);
        });
    }
})
.catch((error) => {
    console.error('Error:', error);
});