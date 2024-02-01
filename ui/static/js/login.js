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
            return true;
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
            return true;
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

       else if (!validatePassword(password)) {
            passwordErrorMsg.style.opacity = "1";
            isValid = false;
        }
        console.log("isValid: " + isValid);

        if (isValid) {
            // Create an object with the form data
            let data = {
                email: email.value,
                password: password.value
            };

            // Convert the data to a JSON string
            let jsonData = JSON.stringify(data);

            // Send a POST request with the data
          fetch('/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: jsonData
            })
                .then(response => response.json())
                .then(data => {
                    if (data.error) {
                        if (data.error === "User does not exist") {
                            // show popup message to user
                            alert("User does not exist. Please register first.")
                            window.location.href = "/register";
                        }
                    } else {
                        // window.location.href = "/login";
                        console.log("logged in");
                        // window.location.href = "/adminPanel";
                        window.location.href = "/home";
                    }
                })
                .catch((error) => {
                    // console.error('Error:', error);
                    if (error === "User does not exist") {
                        // show popup message to user
                        alert("User does not exist. Please register first.")
                        window.location.href = "/register";
                    }
                });
        }
    });

}


