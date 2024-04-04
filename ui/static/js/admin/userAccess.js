console.log('userAccess.js loaded')

function loadUserAccess() {
    fetch('/userAccess', {
        method: 'POST',
    })
        .then(response => response.json())
        .then(data => {
            // Access the adminAccess array
            if (data.success) {
                console.log(data)
                let adminAccess = data.adminAccess;
                let table = `<table>
                <thead>
                    <tr>
                        <th>UserID</th>
                        <th>Username</th>
                        <th>Email</th>
                        <th>Dashboard</th>
                        <th>Upload</th>
                        <th>Edit & Delete</th>
                        <th>Analytics</th>
                        <th>Server Logs</th>
                        <th>User Promote</th>
                    </tr>
                </thead>
                <tbody>
`;

                // Iterate over each object in the adminAccess array
                adminAccess.forEach(user => {
                    table += `<tr>
                    <td>${user.user_id}</td>
                    <td>${user.user_name}</td>
                    <td>${user.email}</td>
                    <td><button class="${user.dashboard_access === 1 ? 'revoke' : 'provide'}" onclick="handleClick('${user.user_name}',${user.user_id}, 'dashboard', ${user.dashboard_access})">${user.dashboard_access === 1 ? 'Revoke' : 'Provide'}</button></td>
                    <td><button class="${user.upload_access === 1 ? 'revoke' : 'provide'}" onclick="handleClick('${user.user_name}',${user.user_id}, 'upload', ${user.upload_access})">${user.upload_access === 1 ? 'Revoke' : 'Provide'}</button></td>
                    <td><button class="${user.edit_delete_access === 1 ? 'revoke' : 'provide'}" onclick="handleClick('${user.user_name}',${user.user_id}, 'edit_delete', ${user.edit_delete_access})">${user.edit_delete_access === 1 ? 'Revoke' : 'Provide'}</button></td>
                    <td><button class="${user.analytics_access === 1 ? 'revoke' : 'provide'}" onclick="handleClick('${user.user_name}',${user.user_id}, 'analytics', ${user.analytics_access})">${user.analytics_access === 1 ? 'Revoke' : 'Provide'}</button></td>
                    <td><button class="${user.server_log_access === 1 ? 'revoke' : 'provide'}" onclick="handleClick('${user.user_name}',${user.user_id}, 'server_log', ${user.server_log_access})">${user.server_log_access === 1 ? 'Revoke' : 'Provide'}</button></td>
                    <td><button class="${user.user_access === 1 ? 'revoke' : 'provide'}" onclick="handleClick('${user.user_name}',${user.user_id}, 'user', ${user.user_access})">${user.user_access === 1 ? 'Revoke' : 'Provide'}</button></td>
                </tr>`;
                });

                table += `</tbody></table>`;

                // Append table to container
                const container = document.getElementById('container');
                container.innerHTML = table;
            } else {
                console.log(data)
            }
        })
        .catch(error => {
            console.log(error)
        });
}

loadUserAccess();

function handleClick(userName, userId, accessType, accessValue) {
    let username = userName;
    let userID = userId;
    let access_type = accessType;
    let access_value = accessValue;
    console.log(`User ID: ${userId}, Access Type: ${accessType}, Selected Option: ${accessValue === 1 ? 'Revoke' : 'Provide'}`);

    data = {
        userID: userID,
        access: access_value,
        accessValue: access_type
    }

    if (access_value === 1) {
        Swal.fire({
            title: "Are you sure?",
            text: `You will be revoking access of ${access_type} from ${username}!`,
            icon: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3085d6",
            cancelButtonColor: "#d33",
            confirmButtonText: "Yes, revoke access!"
        }).then((result) => {
            if (result.isConfirmed) {
                // Send AJAX request to delete the video
                $.ajax({
                    url: '/adminUserAccessChange',
                    type: 'POST',
                    contentType: 'application/json',
                    data: JSON.stringify(data),
                    success: function (response) {
                        console.log(response)
                        Swal.fire({
                            title: "Access Revoked!",
                            text: "User access changed successfully.",
                            icon: "success"
                        });
                        loadUserAccess();
                    },
                    error: function (error) {
                        console.log(error)
                        Swal.fire({
                            title: "Error!",
                            text: "There was an error deleting the file.",
                            icon: "error"
                        });
                    }
                });
            }
        });
    } else {
        Swal.fire({
            title: "Are you sure?",
            text: `You will be give access of ${access_type} to ${username}!`,
            icon: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3085d6",
            cancelButtonColor: "#d33",
            confirmButtonText: "Yes, give access!"
        }).then((result) => {
            if (result.isConfirmed) {
                // Send AJAX request to delete the video
                $.ajax({
                    url: '/adminUserAccessChange',
                    type: 'POST',
                    contentType: 'application/json',
                    data: JSON.stringify(data),
                    success: function (response) {
                        console.log(response)
                        Swal.fire({
                            title: "Access Given!",
                            text: "User access changed successfully.",
                            icon: "success"
                        });
                        loadUserAccess();
                    },
                    error: function (error) {
                        console.log(error)
                        Swal.fire({
                            title: "Error!",
                            text: "There was an error deleting the file.",
                            icon: "error"
                        });
                    }
                });
            }
        });
    }
}