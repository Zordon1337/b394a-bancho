function UpdateUsername() {
    if(document.getElementById('newusername').value == "") {
        alert("Please fill the username field")
        return;
    }
    const formData = new URLSearchParams();
    formData.append('newusername', document.getElementById('newusername').value);
    fetch('/api/v1/userpanel/UpdateUsername', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: formData.toString()
    })
    .then(response => response.text())
    .then(data => {
        alert(data);
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}
