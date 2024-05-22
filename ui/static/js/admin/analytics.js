console.log('analytics.js loaded');

const cards = document.querySelectorAll('.card');
const modal = document.getElementById('myModal');
const modalTitle = document.getElementById('modal-title');
const modalChart = document.getElementById('modalChart').getContext('2d');
let currentChart = null;  // Store the current chart to update or destroy

function showModal(title, chartFunction) {
    modal.style.display = "block";
    modalTitle.textContent = title;
    if (currentChart) {
        currentChart.destroy();  // Destroy the previous chart before drawing a new one
    }
    chartFunction(modalChart);  // Create the new chart in the modal
}

cards.forEach(function(card) {
    card.addEventListener('click', function() {
        console.log('Card clicked:', this.id);
        switch (this.id) {
            case 'views-count-graph':
                showModal('Views Count Graph', createViewsChart);
                break;
            case 'likes-dislikes-graph':
                showModal('Likes and Dislikes Graph', createLikesDislikesChart);
                break;
            case 'video-duration-graph':
                showModal('Video Duration Graph', createDurationChart);
                break;
        }
    });
});

// Close the modal when the user clicks on <span> (x)
document.querySelector('.close').addEventListener('click', function() {
    modal.style.display = "none";
    if (currentChart) {
        currentChart.destroy();  // Ensure to destroy the chart when closing the modal
        currentChart = null;
    }
});

function createViewsChart(ctx) {
    fetch('/mostViewedVideos', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
    })
        .then(response => response.json())
        .then(data => {
            console.log('laude lele mera')
            // Prepare the data for the chart
            const labels = data.map(video => video.Title);
            const views = data.map(video => video.ViewsCount);

            currentChart = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Views',
                        data: views,
                        fill: false,
                        borderColor: 'rgb(75, 192, 192)',
                        tension: 0.1
                    }]
                },
                options: {
                    scales: {
                        y: {
                            beginAtZero: true,
                            ticks: {
                                color: 'white'
                            }
                        },
                        x: {
                            ticks: {
                                color: 'white'
                            }
                        }
                    },
                    plugins: {
                        legend: {
                            labels: {
                                color: 'white' // Add this line to change the color of the labels to white
                            }
                        }
                    }
                }
            });
        })
        .catch(error => {
            console.log(error);
        })
}

function createLikesDislikesChart(ctx) {

    fetch('/likeVsDislike', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => response.json())
        .then(data => {
            const mostLikedVideos = data.mostLikedVideos.slice(0, 5);
            const mostDislikedVideos = data.mostDislikedVideos.slice(0, 5);

            // Combine the two lists into one
            let combinedVideos = mostLikedVideos.concat(mostDislikedVideos);

            // Filter out videos with both likes and dislikes as 0
            combinedVideos = combinedVideos.filter(video => !(video.LikesCount === 0 && video.DislikesCount === 0));


            // Implementation for likes and dislikes chart
            currentChart = new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: combinedVideos.map(video => video.Title),
                    datasets: [{
                        label: 'Likes',
                        data: combinedVideos.map(video => video.LikesCount),
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
                        borderWidth: 1,
                        stack: 'Stack 0'
                    }, {
                        label: 'Dislikes',
                        data: combinedVideos.map(video => video.DislikesCount),
                        backgroundColor: 'rgba(255, 99, 132, 0.2)',
                        borderColor: 'rgba(255, 99, 132, 1)',
                        borderWidth: 1,
                        stack: 'Stack 0'
                    }]
                },
                options: {
                    scales: {
                        y: {
                            beginAtZero: true,
                            ticks: {
                                color: 'white'
                            }
                        },
                        x: {
                            ticks: {
                                color: 'white'
                            }
                        }
                    },
                    plugins: {
                        legend: {
                            labels: {
                                color: 'white' // Add this line to change the color of the labels to white
                            }
                        }
                    }
                }
            });

        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function createDurationChart(ctx) {
    // Implementation for duration distribution chart

    fetch('/duration', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            // Prepare the data for the chart
            const durations = data.map(video => parseFloat(video.Duration));

            // Define the bin edges for the histogram
            const binEdges = [0, 50, 100, 150, 200, 250, 300];

            // Create the bins for the histogram
            const bins = new Array(binEdges.length - 1).fill(0);
            for (const duration of durations) {
                for (let i = 0; i < binEdges.length - 1; i++) {
                    if (binEdges[i] <= duration && duration < binEdges[i + 1]) {
                        bins[i]++;
                        break;
                    }
                }
            }

            // Create the labels for the chart
            const labels = binEdges.slice(0, -1).map((edge, i) => `${edge} - ${binEdges[i + 1]}`);


            currentChart = new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Duration',
                        data: bins,
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
                        borderWidth: 1
                    }]
                },
                options: {
                    scales: {
                        y: {
                            beginAtZero: true,
                            ticks: {
                                color: 'white'
                            }
                        },
                        x: {
                            ticks: {
                                color: 'white'
                            }
                        }
                    }
                }
            });
        })
        .catch(error => {
            console.log(error);
        });
}



/*Line chart:
Config:
const config = {
    type: 'line',
    data: data,
};

Setup:
    const labels = Utils.months({count: 7});
const data = {
    labels: labels,
    datasets: [{
        label: 'My First Dataset',
        data: [65, 59, 80, 81, 56, 55, 40],
        fill: false,
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.1
    }]
};*/

function mostViewedVideos() {
    fetch('/mostViewedVideos', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
    })
        .then(response => response.json())
        .then(data => {
            // Prepare the data for the chart
            const labels = data.map(video => video.Title);
            const views = data.map(video => video.ViewsCount);

            // Create the chart
            const ctx = document.getElementById('viewsChart').getContext('2d');
            new Chart(ctx, {
                type: 'line', // Changed from 'bar' to 'line'
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Views',
                        data: views,
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
                        borderWidth: 1
                    }]
                },
                options: {
                    scales: {
                        y: {
                            beginAtZero: true,
                            // change the y-axis label to white
                            ticks: {
                                color: 'white'
                            }
                        },
                        x: {
                            ticks: {
                                color: 'white'
                            }
                        }
                    },
                }
            });
        })
        .catch(error => {
            console.log(error);
        });
}

function likesVsDislikes() {
    fetch('/likeVsDislike', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => response.json())
        .then(data => {
            const mostLikedVideos = data.mostLikedVideos.slice(0, 5);
            const mostDislikedVideos = data.mostDislikedVideos.slice(0, 5);

            // Combine the two lists into one
            let combinedVideos = mostLikedVideos.concat(mostDislikedVideos);

            // Filter out videos with both likes and dislikes as 0
            combinedVideos = combinedVideos.filter(video => !(video.LikesCount === 0 && video.DislikesCount === 0));

            const ctx = document.getElementById('likesDislikesChart').getContext('2d');
            new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: combinedVideos.map(video => video.Title),
                    datasets: [{
                        label: 'Likes',
                        data: combinedVideos.map(video => video.LikesCount),
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
                        borderWidth: 1,
                        stack: 'Stack 0'
                    }, {
                        label: 'Dislikes',
                        data: combinedVideos.map(video => video.DislikesCount),
                        backgroundColor: 'rgba(255, 99, 132, 0.2)',
                        borderColor: 'rgba(255, 99, 132, 1)',
                        borderWidth: 1,
                        stack: 'Stack 0'
                    }]
                },
                options: {
                    scales: {
                        y: {
                            beginAtZero: true,
                            ticks: {
                                color: 'white'
                            }
                        },
                        x: {
                            ticks: {
                                color: 'white'
                            }
                        }
                    }
                }
            });
        })
        .catch(error => console.error('Error:', error));
}


function durationDistribution() {
    fetch('/duration', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            // Prepare the data for the chart
            const durations = data.map(video => parseFloat(video.Duration));

            // Define the bin edges for the histogram
            const binEdges = [0, 50, 100, 150, 200, 250, 300];

            // Create the bins for the histogram
            const bins = new Array(binEdges.length - 1).fill(0);
            for (const duration of durations) {
                for (let i = 0; i < binEdges.length - 1; i++) {
                    if (binEdges[i] <= duration && duration < binEdges[i + 1]) {
                        bins[i]++;
                        break;
                    }
                }
            }

            // Create the labels for the chart
            const labels = binEdges.slice(0, -1).map((edge, i) => `${edge} - ${binEdges[i + 1]}`);

            // Create the chart
            const ctx = document.getElementById('durationChart').getContext('2d');
            new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Duration',
                        data: bins,
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
                        borderWidth: 1
                    }]
                },
                options: {
                    scales: {
                        y: {
                            beginAtZero: true,
                            ticks: {
                                color: 'white'
                            }
                        },
                        x: {
                            ticks: {
                                color: 'white'
                            }
                        }
                    }
                }
            });
        })
        .catch(error => {
            console.log(error);
        });
}


durationDistribution();
likesVsDislikes();
mostViewedVideos();