function Register(username, password) {
    fetch(`api/v1/register?u=${username}&p=${password}`)
    .then(response => response.text())
    .then(data => {
       if(data.startsWith("ERR")) {
            var dat = data.split('\n')
            alert(dat[1])
       } else {
        document.location.href = `/profile/${username}`
       }
    })
    .catch(error => console.error('Error registering', error));
}