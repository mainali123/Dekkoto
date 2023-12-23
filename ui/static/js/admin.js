console.log('admin.js loaded');
document.addEventListener('DOMContentLoaded', videoValueValidation);

function videoValueValidation() {
    const title = document.getElementById('title');
    const description = document.getElementById('description');

    function validateTitle(title) {
        // Validate email
        if (title.value === '') {
            // titleErrorMsg.innerHTML = 'Title is required';
            alert('Title is required');
            return false;
        } else {
            // titleErrorMsg.innerHTML = '';
            return true;
        }
    }

    function validateDescription(description) {
        // Validate password
        if (description.value === '') {
            // descriptionErrorMsg.innerHTML = 'Description is required';
            alert('Description is required');
            return false;
        } else {
            // descriptionErrorMsg.innerHTML = '';
            return true;
        }
    }

    form = document.querySelector('.form');

    form.addEventListener('submit', function (event) {
        event.preventDefault();
        console.log("submitting form");
        // print the form data to the console
        console.log("title: " + title.value);
        console.log("description: " + description.value);
        let isValid = true;

        if (!validateTitle(title)) {
            // titleErrorMsg.style.opacity = "1";
            isValid = false;
        }

        else if (!validateDescription(description)) {
            // descriptionErrorMsg.style.opacity = "1";
            isValid = false;
        }
        console.log("isValid: " + isValid);

        if (isValid) {
            // Create an object with the form data
            let data = {
                title: title.value,
                description: description.value
            };

            // Send a POST request
            fetch('/videoDetails', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data),
            })
                .then(function (response) {
                    return response.json();
                })
                .then(function (jsonResponse) {
                    // Display response back to user
                    console.log(jsonResponse);
                    alert(jsonResponse.message);
                    // window.location.href = jsonResponse.redirect;
                })
                .catch(function (error) {
                    console.log(error);
                });
        }
    });
}