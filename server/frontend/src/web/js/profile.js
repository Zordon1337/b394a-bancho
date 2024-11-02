document.addEventListener('DOMContentLoaded', async function() {
    const hash = window.location.href;
    if (window.location.href.includes("/profile/")) {
        const playerId = hash.split("/profile/")[1];
        var fet = await fetch(`/api/v1/findplayer/${playerId}`)
        var response = await fet.json()
        document.getElementById("avatar").src = `/forum/download.php?avatar=${response.UserId}`
        document.getElementById("avatar").width = 128
        document.getElementById("avatar").height = 128
        document.getElementById("Username").textContent = `${response.Username} (#${response.Rank})`
        document.getElementById("RankedScore").textContent = `Ranked Score: ${response.RankedScore}`
        document.getElementById("TotalScore").textContent = `Total Score: ${response.TotalScore}`
        document.getElementById("Accuracy").textContent = `Accuracy: ${response.Accuracy*100}%`
        document.getElementById("PlayCount").textContent = `Play Count: ${response.PlayCount}`
        document.getElementById("JoinDate").textContent = `Join Date: ${response.JoinDate}`
        document.getElementById("LastOnline").textContent = `Last Online: ${response.LastOnline}`
        const topScores = response.topscores;
        const table = document.createElement('table');
        table.className = 'top-scores-table';

        const headerRow = document.createElement('tr');
        const headers = ['Map Name', 'Total Score', 'Accuracy'];
        headers.forEach(headerText => {
            const th = document.createElement('th');
            th.textContent = headerText;
            headerRow.appendChild(th);
        });
        table.appendChild(headerRow);

        topScores.forEach(score => {
            const row = document.createElement('tr');
            const mapNameCell = document.createElement('td');
            mapNameCell.textContent = score.MapName;

            const totalScoreCell = document.createElement('td');
            totalScoreCell.textContent = score.TotalScore.toLocaleString();

            const accuracyCell = document.createElement('td');
            accuracyCell.textContent = `${(score.Accuracy * 100).toFixed(2)}%`; 

            row.appendChild(mapNameCell);
            row.appendChild(totalScoreCell);
            row.appendChild(accuracyCell);
            table.appendChild(row);
        });

        document.getElementById('top5scoresdiv').appendChild(table);
    } else {
        
    }
});
