function handleHashRoute() {
    const hash = window.location.hash;
    if (hash.startsWith("#/player/")) {
        const playerId = hash.split("#/player/")[1];
    } else {
        
    }
}

window.addEventListener('hashchange', handleHashRoute);

handleHashRoute();