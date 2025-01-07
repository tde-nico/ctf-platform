function loginUser() {
    var usernameField = document.getElementById("LoginUsername");
    var passwordField = document.getElementById("LoginPassword");

    var username = usernameField.value;
    var password = passwordField.value;
    var formData = new FormData();
    formData.append("username", username);
    formData.append("password", password);

    fetch("/login", {
        method: "POST",
        body: formData
    })
    .then( response => {
        if (response.status == 200) {
            window.location.href = "/challenges";
        }
        else {
            window.location.reload();
        }
    })
    .catch(error => console.error("Error:", error));
}

