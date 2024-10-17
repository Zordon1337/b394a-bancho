function Login(username, password) {
    fetch(`api/v1/login?u=${username}&p=${password}`)
    .then(response => response.text())
    .then(data => {
       if(data.startsWith("ERR")) {
            var dat = data.split('\n')
            alert(dat[1])
       } else {
        alert("Successfully logged in!")
       }
    })
    .catch(error => console.error('Error logging in', error));
}