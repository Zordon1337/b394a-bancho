function UpdatePassword() {
    if(document.getElementById('oldpass').value == "" || document.getElementById('newpass').value == "") {
        alert("Please fill the fields")
        return;
    }
    const formData = new URLSearchParams();
    formData.append('oldpass', document.getElementById('oldpass').value);
    formData.append('newpass', document.getElementById('newpass').value);
    fetch('/api/v1/userpanel/UpdatePassword', {
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
