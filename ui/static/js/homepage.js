console.log("Homepage.js loaded")


// carousel
// Fetch carousel data from the server
fetch('/caroselSlide', {
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
        const carouselData = data.videos;
        const carouselContainer = document.querySelector('.carousel-inner'); // Select the HTML element where you want to display the carousel items

        // Check if there are carousel items
        if (carouselData.length > 0) {
            carouselData.forEach((item, index) => {
                // Create HTML elements for each carousel item
                const carouselItem = document.createElement('div');
                carouselItem.className = index === 0 ? 'carousel-item active' : 'carousel-item';

                const img = document.createElement('img');
                img.className = 'd-block w-100';
                img.src = '../../' + item.ThumbnailURL.replace('thumbnails', 'banners'); // Adjust the path relative to the HTML file
                img.alt = item.Title;
                img.dataset.videoDetails = JSON.stringify(item); // Store the video details in the img element
                console.log(img.dataset.videoDetails)
                // on click to the active carousel item, redirect to the watch video page
                carouselItem.addEventListener('click', function () {
                    console.log('Carousel item clicked'); // Add this line
                    // save the video details in local storage
                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                    window.location.href = '/watchVideo'; // Redirect to the watch video page
                });

                // Append the created HTML elements to the carousel item
                carouselItem.appendChild(img);

                // Append the carousel item to the carousel container
                carouselContainer.appendChild(carouselItem);
            });
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });

// Recently Added Videos
fetch('/recentlyAdded', {
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
        const videoContainer = document.querySelector('.card-video-recently-added'); // Select the HTML element where you want to display the videos

        // Check if there are videos
        if (videos.length > 0) {
            videos.forEach(video => {
                // Create HTML elements for each video
                const card = document.createElement('div');
                card.className = 'card';

                const img = document.createElement('img');
                img.className = 'card-img-top';
                img.src = '../../' + video.ThumbnailURL; // Adjust the path relative to the HTML file
                img.dataset.videoDetails = JSON.stringify(video); // Store the video details in the img element
                img.addEventListener('click', function () {
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


// Recommended Videos
// Recommended Anime
fetch('/recommendedVideos', {
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
        const videos = data.videos;
        const videoContainer = document.querySelector('.recommended-div'); // Select the HTML element where you want to display the videos
        const recommendedAnime = document.querySelector('.recommended_anime'); // Select the HTML element where you want to display the videos

        // Check if there are videos
        if (videos.length > 0) {
            videos.forEach(video => {
                // Create HTML elements for each video

                recommendedAnime.style.backgroundImage = "linear-gradient(0deg, rgba(0,0,0,0) 0%, rgba(0,0,0,0.9136029411764706) 83%),url('../../" + video.ThumbnailURL.replace('thumbnails', 'banners') + "')";
                // Adjust the path relative to the HTML file

                // Add video details to the video container
                recommendedAnime.dataset.videoDetails = JSON.stringify(video);
                // On click, redirect to the watch video page
                recommendedAnime.addEventListener('click', function () {
                    // save the video details in local storage
                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                    window.location.href = '/watchVideo'; // Redirect to the watch video page
                });

                // const img = document.createElement('img');
                // img.className = 'image reco-img';
                // img.src = '../../' + video.ThumbnailURL.replace('thumbnails', 'banners'); // Adjust the path relative to the HTML file
                // img.alt = video.Title;

                const cardText = document.createElement('p');
                cardText.className = ' reco-text';
                cardText.textContent = video.Title;

                const rating = document.createElement('p');
                rating.className = 'rating reco-rating';
                rating.textContent = 0; // Assuming the video object has a Rating property

                // Append the created HTML elements to the video container
                // videoContainer.appendChild(img);
                videoContainer.appendChild(cardText);
                videoContainer.appendChild(rating);
            });
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });

// Weekly Top Videos
fetch('/weeklyTop', {
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
        const videos = data.videos;
        const videoContainer = document.querySelector('.card-video-weekly-top'); // Select the HTML element where you want to display the videos

        // Check if there are videos
        if (videos.length > 0) {
            videos.forEach(video => {
                // Create HTML elements for each video
                const card = document.createElement('div');
                card.className = 'card';

                const img = document.createElement('img');
                img.className = 'card-img-top';
                img.src = '../../' + video.ThumbnailURL; // Adjust the path relative to the HTML file
                img.dataset.videoDetails = JSON.stringify(video); // Store the video details in the img element
                img.addEventListener('click', function () {
                    // save the video details in local storage
                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                    window.location.href = '/watchVideo'; // Redirect to the watch video page
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

// Continue Watching Videos
fetch('/continueWatching', {
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
        const videos = data.videos;
        const videoContainer = document.querySelector('.card-video-continue-watching'); // Select the HTML element where you want to display the videos

        // Check if there are videos
        if (videos.length > 0) {
            videos.forEach(video => {
                // Create HTML elements for each video
                const card = document.createElement('div');
                card.className = 'card';

                const img = document.createElement('img');
                img.className = 'card-img-top';
                img.src = '../../' + video.ThumbnailURL; // Adjust the path relative to the HTML file
                img.dataset.videoDetails = JSON.stringify(video); // Store the video details in the img element
                img.addEventListener('click', function () {
                    // save the video details in local storage
                    localStorage.setItem('videoDetails', this.dataset.videoDetails);
                    window.location.href = '/watchVideo'; // Redirect to the watch video page
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

