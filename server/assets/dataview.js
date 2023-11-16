let dataId;
let dataPoint;
let containerElm = document.getElementById('container');
let numDocs;

function getUUID() {
  return localStorage.getItem('verificationUUID');
}

// Function to send the UUID to /auth/uuid
function sendUUID() {
  const savedUUID = getUUID();

  fetch('/auth/uuid', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ uuid: savedUUID }),
  })
    .then((response) => {
      if (!response.ok) {
        // Check for HTTP error status
        if (response.status === 401) {
          window.location.href = '/view-data';
        } else {
          throw new Error('Network response was not ok');
        }
      }
      return response.json();
    })
    .then((data) => {
      // Handle the response from the server
      console.log(data);
      getDataSet();
    })
    .catch((error) => {
      // Handle errors
      console.error('Error:', error.message);
      window.location.href = '/validate';
    });
}

function getDataSet() {
  fetch('/get-data-point', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  })
    .then((response) => {
      if (!response.ok) {
        console.log('53', response);

        throw new Error('Network response was not ok');
      }
      return response.json();
    })
    .then((data) => {
      // Assuming you have a function to process the data (replace with your logic)
      if ('error' in data) {
        containerElm.style.display = 'none';
        return;
      }
      if (data.documents === 1) {
        numDocs = data.documents;
      }
      console.log(data);
      processData(data);
    })
    .catch((error) => {
      // Handle errors
      console.error('Error:', error);
    });
}

function processData(data) {
  // Process and use the data as needed
  containerElm.style.display = 'show';

  dataId = data.randomDocument._id;
  dataPoint = data.randomDocument;

  updateChart(data.randomDocument['mouse-array']);
}

const ctx = document.getElementById('dataGraph');
const myChart = new Chart(ctx, {
  type: 'line',
  data: {
    labels: [], // Add labels dynamically based on mouse array length
    datasets: [
      {
        label: 'Data Series 1',
        borderColor: 'rgba(75, 192, 192, 1)',
        data: [],
        fill: false,
        borderWidth: 2,
      },
    ],
  },
  options: {
    scales: {
      x: {
        type: 'linear',
        position: 'bottom',
        beginAtZero: true,
        suggestedMax: 800,
        ticks: {
          display: false,
        },
      },
      y: {
        beginAtZero: true,
        suggestedMax: 800,
        reverse: true,
        ticks: {
          display: false,
        },
      },
    },
  },
});

// Function to update the chart with mouse array data
function updateChart(mouseArray) {
  const labels = [];
  const dataPoints = [];

  // Iterate through the mouse array and extract x, y values
  mouseArray.forEach((point, index) => {
    labels.push(`Point ${String.fromCharCode(65 + index)}`);
    dataPoints.push({ x: point.x * 800, y: point.y * 800 });
  });

  // Update chart data
  myChart.data.labels = labels;
  myChart.data.datasets[0].data = dataPoints;

  // Update the chart
  myChart.update();
}
function consumeData(approved = true) {
  let jsonData;
  let url;
  if (dataPoint !== undefined && dataId !== undefined) {
    if (approved) {
      jsonData = dataPoint;
      jsonData.uuid = getUUID();
      url = '/approve_data';
      console.log(jsonData);
    } else {
      url = '/remove_data';
      jsonData = {
        _id: dataId,
        uuid: getUUID(),
      };
    }
    dataPoint = undefined;
    dataId = undefined;


    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(jsonData),
    })
      .then((response) => {
        console.log(response);
        if (response.redirected) {
          window.location.href = response.url;
        }
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
      })
      .then((data) => {
        // Handle the response from the server
      })
      .catch((error) => {
        // Handle errors
        console.error('Error:', error);
      });
      
    if (numDocs === 1) {
      setTimeout(() => {
        getDataSet();
      }, 1000);
    } else {
      getDataSet();
    }
  } else {
    getDataSet();
  }
}

// Event listener for the Remove button
document
  .getElementById('removeButton')
  .addEventListener('click', function () {
    consumeData(false); // Pass false as the 'approved' parameter
  });

// Event listener for the Approve button
document
  .getElementById('approveButton')
  .addEventListener('click', function () {
    consumeData(true); // Pass true as the 'approved' parameter
  });

// Example of how to use sendUUID
document.addEventListener('DOMContentLoaded', sendUUID);
