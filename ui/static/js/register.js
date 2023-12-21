console.log('register.js loaded');

document.addEventListener('DOMContentLoaded', registerValidation);

function registerValidation() {
    const form = document.querySelector('form');
    const name = document.querySelector('#name');
    const email = document.querySelector('#email');
    const password = document.querySelector('#password');
    const confirmPassword = document.querySelector('#Conpassword');
    const nameErrorMsg = document.querySelector('.name-error-msg');
    const emailErrorMsg = document.querySelector('.email-error-msg');
    const passwordErrorMsg = document.querySelector('.password-error-msg');
    const confirmPasswordErrorMsg = document.querySelector('.Conpassword-error-msg');
    const confirmPasswordMatchErrorMsg = document.querySelector('.Conpassword-match-error-msg');

    function validateName(name) {
        if (name.value === '') {
            nameErrorMsg.innerHTML = 'Name is required';
            return false;
        } else {
            nameErrorMsg.innerHTML = '';
            return true;
        }
    }

    function validateEmail(email) {
        if (email.value === '') {
            emailErrorMsg.innerHTML = 'Email is required';
            return false;
        } else if (!/^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/.test(email.value)) {
            emailErrorMsg.innerHTML = 'Please enter a valid email';
            return false;
        } else {
            emailErrorMsg.innerHTML = '';
            return true;
        }
    }

    function validatePassword(password) {
        if (password.value === '') {
            passwordErrorMsg.innerHTML = 'Password is required';
            return false;
        } else if (password.value.length < 8) {
            passwordErrorMsg.innerHTML = 'Your password must be more than 8 characters';
            return false;
        } else {
            passwordErrorMsg.innerHTML = '';
            return true;
        }
    }

    function validateConfirmPassword(password, confirmPassword) {
        if (confirmPassword.value === '') {
            confirmPasswordErrorMsg.innerHTML = 'Confirm password is required';
            return false;
        } else if (confirmPassword.value !== password.value) {
            confirmPasswordMatchErrorMsg.innerHTML = 'Your password must be the same as the password';
            return false;
        } else {
            confirmPasswordErrorMsg.innerHTML = '';
            confirmPasswordMatchErrorMsg.innerHTML = '';
            return true;
        }
    }

    form.addEventListener('submit', function (event) {
        event.preventDefault();
        let isValid = true;

        if (!validateName(name)) {
            nameErrorMsg.style.opacity = "1";
            isValid = false;
        }

        if (!validateEmail(email)) {
            emailErrorMsg.style.opacity = "1";
            isValid = false;
        }

        if (!validatePassword(password)) {
            passwordErrorMsg.style.opacity = "1";
            isValid = false;
        }

        if (!validateConfirmPassword(password, confirmPassword)) {
            confirmPasswordErrorMsg.style.opacity = "1";
            confirmPasswordMatchErrorMsg.style.opacity = "1";
            isValid = false;
        }

        if (isValid) {
            // Create an object with the form data
            let data = {
                name: name.value,
                email: email.value,
                password: password.value,
                confirmPassword: confirmPassword.value
            };

            // Convert the data to a JSON string
            let jsonData = JSON.stringify(data);

            // Send a POST request with the data
            fetch('/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: jsonData
            })
                .then(response => response.json())
                .then(data => {
                    if (data.error) {
                        console.error('Error:', data.error);
                    } else {
                        // Redirect the user to another page or display a success message
                        // console.log('Success:', data);
                        alert("Register successfully. Please login now.")
                        window.location.href = "/login";
                    }
                })
                .catch((error) => {
                    // console.error('Error:', error);
                    if (error === "User already exists") {
                        // show popup message to user
                        alert("User already exists")
                    }
                });
        }
    });
}

