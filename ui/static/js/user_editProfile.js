console.log('user_editProfile.js');

fetch('/userDetails', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    }
})
    .then(response => response.json())
    .then(data => {
        let name = data.userName;
        let email = data.email;

        // show the value in the input fields
        document.getElementById('name').value = name;
        document.getElementById('email').value = email;
    });


document.querySelector('.button').addEventListener('click', function (event) {
    event.preventDefault();

    sendData();
});

function sendData() {
    var name = document.getElementById('name').value;
    var email = document.getElementById('email').value;

    fetch('/editUserProfile', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            userName: name,
            email: email
        })
    })
        // handle the response by showing a message to the user
        .then(response => response.json())
        .then(data => {
            alert(data.message);
        })
        .catch(error => {
            alert('Error:', error);
        });
}