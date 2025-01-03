function loginUser() {
    var usernameField = document.getElementById("LoginUsername");
    var passwordField = document.getElementById("LoginPassword");

    var username = usernameField.value;
    var password = passwordField.value;
    var data = {
        username: username,
        password: password
    };
    var formData = new FormData();
    formData.append("username", username);
    formData.append("password", password);

    fetch("/login", {
        method: "POST",
        body: formData
    })
    .then(response => {
        if (response.status == 200) {
            window.location.href = "/";
        }
        else {
            alert("Invalid username or password");
        }
    })
    .catch(error => console.error("Error:", error));
}

