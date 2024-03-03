console.log("loaded forgetPassword.js");

import { showNotification } from './notification.js';


/*if (registrationSuccessMessage) {
    // toastr.success(registrationSuccessMessage);
    showNotification('success', 'toast-top-right', registrationSuccessMessage);
    localStorage.removeItem('registrationSuccess');
}*/

document.getElementById('resetPassword').addEventListener('submit', function(event) {
    event.preventDefault();
    let email = event.target.elements.email.value;

    console.log(JSON.stringify({email}))

    fetch('/sendEmail', {
        method: 'POST',
        body: JSON.stringify({email}),
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                showNotification('error', 'toast-top-right', data.error);
            } else {
                showNotification('success', 'toast-top-right', data.message);
            }
        })
        .catch((error) => {
            console.error('Error:', error);
        });
});