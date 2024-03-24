fetch('/serverLogsPost', {
    method: 'POST',
})
    .then(response => response.json())
    .then(data => {
        // Access the logs array
        let logs = data.logs;

        // Sort logs in descending order based on the last LoginTime
        logs.sort((a, b) => {
            const aLastLoginTime = JSON.parse(JSON.parse(a.LoginTime)[JSON.parse(a.LoginTime).length - 1]).login_time;
            const bLastLoginTime = JSON.parse(JSON.parse(b.LoginTime)[JSON.parse(b.LoginTime).length - 1]).login_time;
            return new Date(bLastLoginTime) - new Date(aLastLoginTime);
        });

        // Create table and table header
        const table = document.createElement('table');
        const thead = document.createElement('thead');
        const headers = ["IP", "DeviceType", "DeviceOS", "Browser", "LoginTime", "CountryCode", "CountryName", "RegionName", "CityName", "Latitude", "Longitude", "ZipCode", "TimeZone", "ASN", "AS", "IsProxy"];
        headers.forEach(header => {
            const th = document.createElement('th');
            th.textContent = header;
            thead.appendChild(th);
        });
        table.appendChild(thead);

        // Create table body
        const tbody = document.createElement('tbody');
        logs.forEach(log => {
            const tr = document.createElement('tr');
            headers.forEach(header => {
                const td = document.createElement('td');
                if (header === 'LoginTime') {
                    // Parse the LoginTime string to JSON
                    const loginTimes = JSON.parse(log[header]);
                    // Get the last login time
                    const lastLoginTime = loginTimes[loginTimes.length - 1];
                    // Parse the last login time to JSON
                    const lastLoginTimeObj = JSON.parse(lastLoginTime);
                    // Display the last login time in CSV format
                    td.textContent = lastLoginTimeObj.login_time;
                } else {
                    td.textContent = log[header];
                }
                tr.appendChild(td);
            });
            tbody.appendChild(tr);
        });
        table.appendChild(tbody);

        // Append table to container
        const container = document.getElementById('container'); // Replace 'container' with the id of your container
        container.appendChild(table);
    })
    .catch(error => {
        console.log(error)
    });