console.log('user_editProfile.js');

import { showNotification } from './notification.js';

fetch('/userDetails', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    }
})
    .then(response => response.json())
    .then(data => {
        let name = data.userName;
        let email = data.email;

        // show the value in the input fields
        document.getElementById('name').value = name;
        document.getElementById('email').value = email;
    });


document.querySelector('.button').addEventListener('click', function (event) {
    event.preventDefault();

    sendData();
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
            let imageElement = document.getElementById('updateImg');
            imageElement.src = image;
        }
    })
    .catch(error => {
        console.log('Error:', error);
    });

function sendData() {
    let name = document.getElementById('name').value;
    let email = document.getElementById('email').value;

    fetch('/editUserProfile', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            userName: name,
            email: email
        })
    })
        // handle the response by showing a message to the user
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                showNotification('error', 'toast-top-right', data.error);
            } else {
                showNotification('success', 'toast-top-right', data.message);
            }
        })
        .catch(error => {
            alert('Error:', error);
        });
}
