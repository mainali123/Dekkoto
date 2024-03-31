function encryptAndSend(input) {
    let file = input.files[0];

    if (file) {

        let reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = function () {
            let data = reader.result;
            // Extracting base64 string without the prefix
            let base64String = data.split(',')[1]; // Split at the comma and take the second part
            // console.log(base64String);

            fetch('/imageUploadDynamic', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    image: base64String // Send only the base64 string without the prefix
                })
            })
                .then(response => response.json())
                .then(data => {
                    console.log(data);
                })
        }
    } else {
        console.log('No file selected');
    }
}
