console.log('user_changePassword.js');

document.querySelector('.button').addEventListener('click', function (event) {
    event.preventDefault();

    var oldPassword = document.getElementById('old-password').value;
    var newPassword = document.getElementById('new-password').value;
    var confirmPassword = document.getElementById('confirm-password').value;

    if (newPassword !== confirmPassword) {
        alert('New password and confirm password do not match');
        return;
    }

    fetch('/changePassword', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            oldPassword: oldPassword,
            newPassword: newPassword
        })
    })
        .then(response => response.json())
        .then(data => {
            // console.log(data.message);
            // alert(data.message);
            // handle the 500 error
            if (data.error) {
                alert(data.error);
                return;
            }
        })
        .catch(error => {
            console.log('Error:', error)
            alert('Error:', error);
        });
});