let dataId;
let dataPoint;
let containerElm = document.getElementById('container');
let numDocs;

function getUUID() {
  return localStorage.getItem('verificationUUID');
}

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
  containerElm.style.display = 'show';

  dataId = data.randomDocument._id;
  dataPoint = data.randomDocument;

  updateChart(data.randomDocument['mouse-array']);
}

const ctx = document.getElementById('dataGraph');
const myChart = new Chart(ctx, {
  type: 'line',
  data: {
    labels: [], 
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

function updateChart(mouseArray) {
  const labels = [];
  const dataPoints = [];

  // Iterate through the mouse array and extract x, y values
  mouseArray.forEach((point, index) => {
    labels.push(`Point ${String.fromCharCode(65 + index)}`);
    dataPoints.push({ x: point.x * 800, y: point.y * 800 });
  });

  myChart.data.labels = labels;
  myChart.data.datasets[0].data = dataPoints;

  myChart.update();
}
function consumeData(approved) {
  let jsonData;
  let url;
  if (approved) {
    jsonData = dataPoint;
    jsonData.uuid = getUUID();
    url = '/approve_data';
  } else {
    url = '/remove_data';
    jsonData = {
      _id: dataId,
      uuid: getUUID(),
    };
  }

  fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(jsonData),
  })
    .then((response) => {
      if (response.redirected) {
        window.location.href = response.url;
      }
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
    })
    .then((data) => {
      getDataSet();
    })
    .catch((error) => {
      // Handle errors
      console.error('Error:', error);
    });

  if (numDocs === 1) {
    setTimeout(() => {
      getDataSet();
    }, 1000);
  } 
}

document
  .getElementById('removeButton')
  .addEventListener('click', function () {
    consumeData(false);
  });

document
  .getElementById('approveButton')
  .addEventListener('click', function () {
    consumeData(true);
  });

document.addEventListener('DOMContentLoaded', sendUUID);
