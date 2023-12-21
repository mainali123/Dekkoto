console.log('login.js loaded');
document.addEventListener('DOMContentLoaded', loginValidation);

function loginValidation() {
    const form = document.querySelector('.form');
    const email = document.querySelector('#email');
    const password = document.querySelector('#password');
    const emailErrorMsg = document.querySelector('.error-msg');
    const passwordErrorMsg = document.querySelector('.password-error-msg');


    function validateEmail(email) {
        // Validate email
        if (email.value === '') {
            emailErrorMsg.innerHTML = 'Email is required';
            return false;
        } else if (!/^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/.test(email.value)) {
            emailErrorMsg.innerHTML = 'Please enter a valid email';
            return false;
        } else {
            emailErrorMsg.innerHTML = '';
        }
    }

    function validatePassword(password) {
        // Validate password
        if (password.value === '') {
            passwordErrorMsg.innerHTML = 'Password is required';
            return false;
        } else if (password.value.length < 8) {
            passwordErrorMsg.innerHTML = 'Your password must be more than 8 characters';
            return false;
        } else {
            passwordErrorMsg.innerHTML = '';
        }
    }

    form.addEventListener('submit', function (event) {
        event.preventDefault();
        console.log("submitting form");
        let isValid = true;

        if (!validateEmail(email)) {
            emailErrorMsg.style.opacity = "1";
            isValid = false;
        }

        if (!validatePassword(password)) {
            passwordErrorMsg.style.opacity = "1";
            isValid = false;
        }

        if (!isValid) {
            console.log("form is invalid");
        }
    });

}

