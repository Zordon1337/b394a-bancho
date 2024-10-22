document.addEventListener('DOMContentLoaded', function() {
    fetch('/api/v1/GetUser')
    .then(response => response.json())
    .then(data => {
        if (data.loggedIn) {
            const navbarRight = document.getElementById('navbar-right');
            navbarRight.innerHTML = `
                <p>Logged in as<a href='/profile/${data.username}'>${data.username}</a> <a href='/logout'>(logout)</a></p>
            `;
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });

});