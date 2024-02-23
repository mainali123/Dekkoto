console.log('user_profile.js');

fetch('/userDetails', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    }
})
.then(response => response.json())
.then(data => {
    let email = data.email;
    let uniqueName = email.substring(0, email.indexOf('@'));
    let userName = data.userName;

    // Get the elements from the DOM
    let uniqueNameElement = document.querySelector('.unique-name');
    let userNameElement = document.querySelector('.username');

    // Update the textContent of the elements
    uniqueNameElement.textContent = '@' + uniqueName;
    userNameElement.textContent = userName;
});