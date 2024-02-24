console.log('login.js loaded');


import { showNotification } from './notification.js';


// If register is success
document.addEventListener('DOMContentLoaded', (event) => {
    const registrationSuccessMessage = localStorage.getItem('registrationSuccess');
    if (registrationSuccessMessage) {
        // toastr.success(registrationSuccessMessage);
        showNotification('success', 'toast-top-right', registrationSuccessMessage);
        localStorage.removeItem('registrationSuccess');
    }
});

showNotification('info', 'toast-top-right', 'Enter login credentials to continue.');


let showPassword = document.querySelectorAll(".password-show");

showPassword.forEach((show) => {
	show.addEventListener("click", () => {
		const password = show.parentElement.children[1];
		if (password.type === "password") {
			password.type = "text";
			show.innerHTML = '<i class="fa-solid fa-eye-slash"></i>';
		} else {
			password.type = "password";
			show.innerHTML = '<i class="fa-solid fa-eye"></i>';
		}
	});
});

const alert = document.querySelector(".alert-box-dialog");
const alertMessage = document.querySelector(".alert-message-dialog");

function showAlert(message) {
	alertMessage.innerHTML = message;
	alert.style.visibility = "visible";
	alert.style.opacity = "1";

	setTimeout(() => {
		alert.style.visibility = "hidden";
		alert.style.opacity = "0";
	}, 5000);
}

function loginValidation() {
    const form = document.querySelector('.form');
    const email = document.querySelector('#email');
    const password = document.querySelector('#password');


    function validateEmail(email) {
        // Validate email
        if (email.value === '') {
            showAlert("Email is required");
            return false;
        } else if (!/^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/.test(email.value)) {
            showAlert("Please enter valid email");
            return false;
        } else {
            return true;
        }
    }

    function validatePassword(password) {
        // Validate password
        if (password.value === '') {
            showAlert("Password is required");
            return false;
        } else if (password.value.length < 8) {
        showAlert("Your password must be more than 8 characters");
            return false;
        } else {
            return true;
        }
    }

    form.addEventListener('submit', function (event) {
        event.preventDefault();
        console.log("submitting form");
        let isValid = true;

        if (!validateEmail(email)) {
            isValid = false;
        }

       else if (!validatePassword(password)) {
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
            console.log(jsonData);

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
                        localStorage.setItem('loginSuccess', 'Login successful.');
                        window.location.href = "/home";
                    }
                })
                .catch((error) => {
                    // console.error('Error:', error);
                    if (error === "User does not exist") {
                        // show popup message to user
                        console.log("User does not exist. Please register first.")
                        alert("User does not exist. Please register first.")
                        window.location.href = "/register";
                    }
                });
        }
    });

}

loginValidation();
