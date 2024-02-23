console.log('register.js loaded');

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

// showAlert("This is a test message");

document.addEventListener('DOMContentLoaded', registerValidation);

function registerValidation() {
    const form = document.querySelector('form');
    const name = document.querySelector('#name');
    const email = document.querySelector('#email');
    const password = document.querySelector('#password');
    const confirmPassword = document.querySelector('#Conpassword');
    // const nameErrorMsg = document.querySelector('.name-error-msg');
    // const emailErrorMsg = document.querySelector('.email-error-msg');
    // const passwordErrorMsg = document.querySelector('.password-error-msg');
    // const confirmPasswordErrorMsg = document.querySelector('.Conpassword-error-msg');
    // const confirmPasswordMatchErrorMsg = document.querySelector('.Conpassword-match-error-msg');

    function validateName(name) {
        if (name.value === '') {
            return false;
        } else {
            return true;
        }
    }

    function validateEmail(email) {
        if (email.value === '') {
            return false;
        } else if (!/^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/.test(email.value)) {
            return false;
        } else {
            return true;
        }
    }

    function validatePassword(password) {
        if (password.value === '') {
            return false;
        } else if (password.value.length < 8) {
            return false;
        } else {
            return true;
        }
    }

    function validateConfirmPassword(password, confirmPassword) {
        if (confirmPassword.value === '') {
            return false;
        } else if (confirmPassword.value !== password.value) {
            return false;
        } else {
            return true;
        }
    }

    form.addEventListener('submit', function (event) {
        event.preventDefault();
        let isValid = true;

        if (!validateName(name)) {
            showAlert("Name is required");
            isValid = false;
        }

        if (!validateEmail(email)) {
            showAlert("Please enter a valid email");
            isValid = false;
        }

        if (!validatePassword(password)) {
            showAlert("Your password must be more than 8 characters");
            isValid = false;
        }

        if (!validateConfirmPassword(password, confirmPassword)) {
            showAlert("Your password must be the same as the password");
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
            console.log(jsonData)

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
                        // alert("Register successfully. Please login now.")
                        window.location.href = "/login";
                    }
                })
                .catch((error) => {
                    // console.error('Error:', error);
                    if (error === "User already exists") {
                        // show popup message to user
                        // alert("User already exists")
                    }
                });
        }
    });
}

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
