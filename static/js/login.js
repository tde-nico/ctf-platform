function signIn(event){
    alert(event);
    var email = document.getElementById("login-email").value;
    var password = document.getElementById("login-password").value;
    var data = {
        email: email,
        password: password,
    };
    fetch('/login', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(function(response){
        if(response.status === 200){
            window.location.href = '/home';
        }else{
            alert('Login failed. Please try again');
        }
    });
}

function signUp(){
    var email = document.getElementById("register-email").value;
    var password = document.getElementById("register-password").value;
    var confirmPassword = document.getElementById("register-confirm-password").value;
    if(password !== confirmPassword){
        alert('Passwords do not match');
        return;
    }
    var data = {
        email: email,
        password: password,
    };
    fetch('/signup', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(function(response){
        if(response.status === 200){
            window.location.href = '/home';
        }else{
            alert('Signup failed. Please try again');
        }
    });
}