console.log('video list Loaded successfully')


function convertToMinutesSeconds(timeInSeconds) {
    let minutes = Math.floor(timeInSeconds / 60);
    let seconds = Math.floor(timeInSeconds % 60);
    return `${minutes}:${seconds < 10 ? '0' + seconds : seconds}`;
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}

// Fetch the data from the server
function showVideos() {
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
        .then(videoData => {
            console.log(videoData);
            let count = 1;
            let videoList = document.getElementById('videoList');
            // videoList.innerHTML = '';

            videoData.videos.forEach(video => {
                let videoItem = document.createElement('div');
                videoItem.className = 'videoList';
                let genreHTML = video.CategoryName.split(',').map(category => `<p>${capitalizeFirstLetter(category.trim())}</p>`).join('');
                videoItem.innerHTML = `
                    <div class="thumbnail-duration">
                    <img src="${video.ThumbnailURL}" alt="${video.Title}" alt="">
                    <p class="duration">${convertToMinutesSeconds(video.Duration)}</p>
                        <div class="title-desc">
                            <h3>${video.Title}</h3>
                            <div class="date-category">
                                <p class="uploaded-date">Uploaded : ${new Date(video.UploadDate).toDateString()}</p>|
                                <p class="category">${capitalizeFirstLetter(video.GenreName)}</p>
                            </div>
                        </div>
                    </div>
                    <div class="video-stats">
                        <p class="views">
                            <i class="fa-solid fa-eye" style="color: #919191;"></i>${video.ViewsCount}
                        </p>.
                        <p class="likes">
                            <i class="fa-solid fa-thumbs-up" style="color: #919191;"></i>${video.LikesCount}
                        </p>.
                        <p class="dislikes">
                            <i class="fa-solid fa-thumbs-down" style="color: #919191;"></i>${video.DislikesCount}
                        </p>
                    </div>
                    <div class="genre">
                        ${genreHTML}
                    </div>
                    <div class="btn-container">
                        <button class="edit-btn" id="editBtn${count}" metaData='${btoa(encodeURIComponent(JSON.stringify(video)))}'><i class="fa-solid fa-pen-to-square"></i></button>
                        <button class="delete-btn" id="deleteBtn${count}" metaData='${btoa(encodeURIComponent(JSON.stringify(video)))}'><i class="fa-solid fa-trash"></i></button>
                    </div>
                `;
                videoList.appendChild(videoItem);

                // if delete button is clicked
                document.getElementById(`deleteBtn${count}`).addEventListener('click', function () {
                    let videoData = JSON.parse(decodeURIComponent(atob(this.getAttribute('metaData'))));
                    console.log(videoData);
                    deleteVideo(videoData.VideoID); // pass the VideoID to the deleteVideo function
                });
                // if edit button is clicked
                document.getElementById(`editBtn${count}`).addEventListener('click', function () {
                    let videoData = JSON.parse(decodeURIComponent(atob(this.getAttribute('metaData'))));
                    console.log(videoData);
                    editVideo(videoData); // pass the VideoID to the editVideo function
                });

                count++;
            });
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
        });
}

function deleteVideo(videoId) {
    Swal.fire({
        title: "Are you sure?",
        text: "You won't be able to revert this!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes, delete it!"
    }).then((result) => {
        if (result.isConfirmed) {
            // Send AJAX request to delete the video
            $.ajax({
                url: '/deleteVideo',
                type: 'POST',
                contentType: 'application/json',
                data: JSON.stringify({videoID: videoId}),
                success: function (response) {
                    Swal.fire({
                        title: "Deleted!",
                        text: "Your file has been deleted.",
                        icon: "success"
                    });
                    showVideos();
                },
                error: function (error) {
                    Swal.fire({
                        title: "Error!",
                        text: "There was an error deleting the file.",
                        icon: "error"
                    });
                }
            });
        }
    });
}

let editVideoModal = document.getElementById('editDialog');


function editVideo(videoData) {
    editVideoModal.showModal();

    //add external html to edit dialog using fetch
    fetch('/adminPanel/addVideo')
        .then(response => response.text())
        .then(html => {
            editVideoModal.innerHTML = html;

            // remove its action
            let form = document.querySelector('.upload-video form');
            form.removeAttribute('action');
            form.removeAttribute('method');
            form.removeAttribute('enctype');

            document.querySelector(".upload-area").style.display = "none";

            //remove submit button and add update button
            let submitBtn = document.querySelector('#submitBtn');
            submitBtn.style.display = 'none';
            let updateBtn = document.createElement('button');
            updateBtn.id = 'updateBtn';
            updateBtn.className = 'primary-btn';
            updateBtn.innerHTML = 'Update';
            submitBtn.parentElement.appendChild(updateBtn);

            //remove reset btn
            let resetBtn = document.querySelector('#resetBtn');
            resetBtn.style.display = 'none';

            //populate data
            document.getElementById('title').value = videoData.Title;
            document.getElementById('description').value = videoData.Description;
            document.querySelector(".previewImgThumbnail").src = videoData.ThumbnailURL

            document.querySelector(".previewImgBanner").src = videoData.ThumbnailURL.replace("thumbnails", "banners");

            //check the radio button
            document.getElementById(videoData.GenreName).checked = true;

            //check genre checkboxes
            let genreList = videoData.CategoryName.split(',');
            genreList.forEach(genre => {
                document.getElementById(genre.trim()).checked = true;
                console.log(genre.trim());
            });

            //create close button
            let closeBtn = document.createElement('button');
            closeBtn.id = 'closeBtn';
            closeBtn.className = 'primary-btn';
            closeBtn.innerHTML = '<i class="fa-solid fa-xmark"></i> ';
            editVideoModal.appendChild(closeBtn);

            //add event listener to update button


            //add event listener to close button
            document.getElementById('closeBtn').addEventListener('click', function () {
                editVideoModal.close();
            });

            //add event listener to update button
            updateBtn.addEventListener('click', function (event) {
                event.preventDefault();
                console.log(videoData)

                console.log('update button clicked')
                let title = document.getElementById('title').value;
                let description = document.getElementById('description').value;
                let thumbnail = document.querySelector(".previewImgThumbnail").src;
                let banner = document.querySelector(".previewImgBanner").src;
                let genre= document.querySelector('input[name="type"]:checked').value;
                const checkboxes = document.querySelectorAll('input[name="genre"]:checked');
                const category = Array.from(checkboxes).map(checkbox => checkbox.value).join(',');
                const fileName = videoData.ThumbnailURL.split('/').pop();

                blobToBase64(thumbnail, function (base64) {
                    console.log('Base64 encoded image:', base64);
                    if (base64.startsWith('data:image/')) {
                        thumbnail = base64;
                    } else {
                        thumbnail = 'Same Image';
                    }

                    blobToBase64(banner, function (base64) {
                        console.log('Base64 encoded image:', base64);
                        if (base64.startsWith('data:image/')) {
                            banner = base64;
                        } else {
                            banner = 'Same Image';
                        }
                    });
                });

                // Convert both images to base64
                Promise.all([
                    blobToBase64(thumbnail),
                    blobToBase64(banner)
                ]).then(([base64Thumbnail, base64Banner]) => {
                    // Both conversions are done here
                    let values = {
                        title: title,
                        description: description,
                        thumbnail: base64Thumbnail,
                        banner: base64Banner,
                        category: category,
                        genre: genre,
                        videoID: videoData.VideoID,
                        fileName: fileName
                    };
                    updateVideo(values);
                }).catch(error => {
                    console.error('There was a problem with the blob to base64 conversion:', error);
                });

            });
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
        });


}

showVideos();

function updateVideo(values) {
    console.log(values.videoID);
    $.ajax({
        url: '/editVideo',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(values),
        success: function (response) {
            Swal.fire({
                title: "Updated!",
                text: "Your file has been updated.",
                icon: "success"
            });
            showVideos();
            editVideoModal.close();
        },
        error: function (error) {
            Swal.fire({
                title: "Error!",
                text: "There was an error updating the file.",
                icon: "error"
            });
        }
    });
}

function blobToBase64(blobUrl) {
    return new Promise((resolve, reject) => {
        if (!blobUrl.startsWith('blob:')) {
            resolve('Same Image');
        }
        let xhr = new XMLHttpRequest();
        xhr.onload = function () {
            let reader = new FileReader();
            reader.onloadend = function () {
                resolve(reader.result);
            }
            reader.readAsDataURL(xhr.response);
        };
        xhr.onerror = reject;
        xhr.open('GET', blobUrl);
        xhr.responseType = 'blob';
        xhr.send();
    });
}