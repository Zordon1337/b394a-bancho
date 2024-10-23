document.addEventListener('DOMContentLoaded', async function() {
    const hash = window.location.href;
    if (window.location.href.includes("/profile/")) {
        const playerId = hash.split("/profile/")[1];
        var fet = await fetch(`/api/v1/findplayer/${playerId}`)
        var response = await fet.json()
        document.getElementById("Username").textContent = `${response.Username} (#${response.Rank})`
        document.getElementById("RankedScore").textContent = `Ranked Score: ${response.RankedScore}`
        document.getElementById("TotalScore").textContent = `Total Score: ${response.TotalScore}`
        document.getElementById("Accuracy").textContent = `Accuracy: ${response.Accuracy*100}%`
        document.getElementById("PlayCount").textContent = `Play Count: ${response.PlayCount}`
        document.getElementById("JoinDate").textContent = `Join Date: ${response.JoinDate}`
        document.getElementById("LastOnline").textContent = `Last Online: ${response.LastOnline}`
    } else {
        
    }
});
