console.log('user_changePassword.js');

import {showNotification} from "./notification.js";


document.querySelector('.button').addEventListener('click', function (event) {
    event.preventDefault();

    let oldPassword = document.getElementById('old-password').value;
    let newPassword = document.getElementById('new-password').value;
    let confirmPassword = document.getElementById('confirm-password').value;

    console.log('oldPassword:', oldPassword);
    console.log('newPassword:', newPassword);
    console.log('confirmPassword)', confirmPassword);

    if (newPassword !== confirmPassword) {
        // alert('New password and confirm password do not match');
        showNotification('error', 'toast-top-right', 'New password and confirm password do not match');
        return;
    }

    fetch('/changePassword', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            oldPassword: oldPassword,
            newPassword: newPassword
        })
    })
        .then(response => response.json())
        .then(data => {
            // console.log(data.message);
            // alert(data.message);
            // handle the 500 error
            if (data.error) {
                showNotification('error', 'toast-top-right', data.error);
            } else {
                showNotification('success', 'toast-top-right', data.message);
            }
        })
        .catch(error => {
            console.log('Error:', error)
            alert('Error:', error);
        });
});