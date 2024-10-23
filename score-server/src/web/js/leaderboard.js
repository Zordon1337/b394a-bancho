document.addEventListener('DOMContentLoaded', function() {
    fetch('api/v1/gettop')
    .then(response => response.json())  
    .then(data => {
        const leaderboardTable = document.querySelector('.content table');

        data.forEach(entry => {
            const row = document.createElement('tr');

            const rankCell = document.createElement('td');
            rankCell.textContent = entry.rank;

            const usernameCell = document.createElement('td');
            usernameCell.innerHTML = `<a href='/profile/${entry.Username}'>${ entry.Username}</a>`;

            const rankedScoreCell = document.createElement('td');
            rankedScoreCell.textContent = entry.RankedScore;

            const totalScoreCell = document.createElement('td');
            totalScoreCell.textContent = entry.TotalScore;

            const accuracyCell = document.createElement('td');
            accuracyCell.textContent = entry.Accuracy*100;

            row.appendChild(rankCell);
            row.appendChild(usernameCell);
            row.appendChild(rankedScoreCell);
            row.appendChild(totalScoreCell);
            row.appendChild(accuracyCell);

            leaderboardTable.appendChild(row);
        });
    })
    .catch(error => console.error('Error fetching leaderboard data:', error));
});