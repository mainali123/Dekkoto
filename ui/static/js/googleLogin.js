import {initializeApp} from "https://www.gstatic.com/firebasejs/10.9.0/firebase-app.js";
import {getAuth, GoogleAuthProvider, signInWithPopup} from "https://www.gstatic.com/firebasejs/10.9.0/firebase-auth.js";
import {showNotification} from "./notification.js";

const firebaseConfig = {
    apiKey: "AIzaSyCPtoBmQaJPtEBWbPcKyP6wFYEvwkS0lo8",
    authDomain: "dekkoto-74fce.firebaseapp.com",
    projectId: "dekkoto-74fce",
    storageBucket: "dekkoto-74fce.appspot.com",
    messagingSenderId: "352710129050",
    appId: "1:352710129050:web:7bd14dcf29be1ad8123623"
};


$('.jquery__register').click(function () {
    registerUser();
});

$('.jquery__login').click(function () {
    loginUser();
});

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const auth = getAuth();
auth.languageCode = 'en';

const provider = new GoogleAuthProvider();

const googleSignupButton = document.getElementById('google-register-button');

// googleSignupButton.addEventListener('click', (event) => {
function registerUser() {
    signInWithPopup(auth, provider)
        .then((result) => {
            const credential = GoogleAuthProvider.credentialFromResult(result);
            const user = result.user;
            // console.log(JSON.stringify(user))
            const uname = user.displayName;
            const email = user.email;
            const imageURL = user.photoURL;
            const emailVerified = user.emailVerified;

            console.log(uname)
            console.log(email)
            console.log(imageURL)
            console.log(emailVerified)
            const password = 'GoogleRegisteredUser'

            if (emailVerified) {
                sendData(uname, email, imageURL, password)
            }

            console.log()
        }).catch((error) => {
        const errorCode = error.code;
        const errorMessage = error.message;
    });
}

// })

function sendData(uname, email, imageURL, password) {
    const bodyJson = JSON.stringify({
        name: uname,
        email: email,
        password: password,
        confirmPassword: password
    })
    fetch('/register', {
        method: 'POST',
        body: bodyJson
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            window.location.href = "/login"
        })
        .then(error => {
            console.log(error)
        })
}

const googleLoginButton = document.getElementById('google-register-button');

// googleLoginButton.addEventListener('click', (event) => {
function loginUser() {
    signInWithPopup(auth, provider)
        .then((result) => {
            const credential = GoogleAuthProvider.credentialFromResult(result);
            const user = result.user;
            // console.log(JSON.stringify(user))
            const uname = user.displayName;
            const email = user.email;
            const imageURL = user.photoURL;
            const emailVerified = user.emailVerified;

            console.log(email)
            console.log(emailVerified)
            const password = 'GoogleRegisteredUser'

            if (emailVerified) {
                login(email, password)
            }

            console.log()
        }).catch((error) => {
        const errorCode = error.code;
        const errorMessage = error.message;
    });
}

// })

function login(email, password) {
    const bodyJson = JSON.stringify({
        email: email,
        password: password,
    })
    fetch('/login', {
        method: 'POST',
        body: bodyJson
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                if (data.error === "User does not exist") {
                    // show popup message to user
                    showNotification('error', 'toast-top-right', "User does not exist. Please register first.")
                    window.location.href = "/register";
                }
                showNotification('error', 'toast-top-right', data.error);
            } else {
                // window.location.href = "/login";
                console.log("logged in");
                // window.location.href = "/adminPanel";
                localStorage.setItem('loginSuccess', 'Login successful.');
                window.location.href = "/home";
            }
        })
        .then(error => {
            console.log(error)
        })
}