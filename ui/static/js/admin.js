// sending video details to server

function sendVideoDetails(title, description) {

    // Create an object with the form data
    let videoDetails = {
        title: title,
        description: description
    };

    // Convert the data to a JSON string
    let jsonData = JSON.stringify(videoDetails);

    // Send a POST request with the data
    fetch('/VideoDetails', {
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

// title and description
let title = document.getElementById("title");
let description = document.getElementById("description");

// Function to validate input against SQL injection and XSS attacks
function validateInput(input) {
    // Regular expression to detect SQL injection and XSS attacks
    const regex = /(\b(FROM|SELECT|INSERT|DELETE|WHERE|DROP|EXEC|UNION|CREATE|ALTER|UPDATE|JOIN)\b)|[<>]/i;

    // Check if input is not empty and doesn't match the regex
    if (input && !regex.test(input)) {
        return true;
    } else {
        return false;
    }
}

// Check if the data is valid
if (validateInput(title.value) && validateInput(description.value)) {
    // If valid, call the sendVideoDetails function
    sendVideoDetails(title.value, description.value);
} else {
    // If not valid, show an error message
    if (!title.value) {
        alert("Title cannot be empty.");
    } else {
        alert("Invalid input. Please avoid using special characters.");
    }
}