console.log('user_profile.js');

// Add event listener to bx-log-out class
document.querySelector('.logout-now').addEventListener('click', function () {
    fetch('/logout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            window.location.href = '/login';
        }
    });
});

fetch('/userDetails', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    }
})
.then(response => response.json())
.then(data => {
    let email = data.email;
    let uniqueName = email.substring(0, email.indexOf('@'));
    let userName = data.userName;

    // Get the elements from the DOM
    let uniqueNameElement = document.querySelector('.unique-name');
    let userNameElement = document.querySelector('.username');

    // Update the textContent of the elements
    uniqueNameElement.textContent = '@' + uniqueName;
    userNameElement.textContent = userName;
});

fetch('/userProfileImage', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    }
})
.then(response => response.json())
.then(data => {
    console.log(data)
    let success = data.success;
    if (success === false) {
        return;
    } else {
        let image = "../../" + data.imagePath;
        console.log(image);
        let imageElement = document.getElementById('user-profile');
        imageElement.src = image;
    }
})
.catch(error => {
    console.log('Error:', error);
});