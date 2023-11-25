let recordingInterval;
let maxDuration = 3000;
let pollFreq = 2;
const windowWidth = 800;
const windowHeight = 800;
let recording = false;

let debug = false;
let docsElm = document.querySelector('#documents');

let docReady = false;

function getNumOfDocs() {
  fetch('document_count', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      docsElm.textContent = data.count;
    })
    .catch((error) => {
      console.error('Error:', error);
    });
}
getNumOfDocs();

function sendData(data) {
  let json_data = JSON.stringify({
    'window-height': windowHeight,
    'window-width': windowWidth,
    'mouse-array': data['mouse-array'],
  });

  if (!debug) {
    fetch('add_data', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: json_data,
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        getNumOfDocs();
      })
      .catch((error) => {
        console.error('Error:', error);
      });
  } else {
    console.log(data);
  }
}

class Cell {
  constructor(x, y, width, margin) {
    this.x = x;
    this.y = y;
    this.margin = margin;
    this.width = width;
    this.filled = false;
    this.margin = 5;
    (this.color = 51), 51, 51, 50;
  }

  isMouseOver() {
    return (
      mouseX > this.x * (this.width + this.margin) &&
      mouseX + Math.round(this.margin / 2) <
        (this.x + 1) * (this.width + this.margin) &&
      mouseY > this.y * (this.width + this.margin) &&
      mouseY + Math.round(this.margin / 2) <
        (this.y + 1) * (this.width + this.margin)
    );
  }

  show() {
    noStroke(); // Remove the stroke to hide the outer border

    let px = this.x * (this.width + this.margin);
    let py = this.y * (this.width + this.margin);
    fill(this.color);
    rect(px, py, this.width, this.width, 5); // Adjust the radius value
  }
}

let touch = false;
document.addEventListener('touchstart', (evt) => {
  touch = true;
});

class Grid {
  constructor(canvasSize, cellCount, cellWidth, margin) {
    this.canvasSize = canvasSize;
    this.cellCount = cellCount;
    this.margin = margin;
    this.cellWidth = cellWidth;
    this.grid = [];
    this.initGrid();
    this.recordStartTime = 0;
    this.isRecording = false;
    this.timeoutId = undefined;
    this.data = {
      'mouse-array': [],
    };
  }

  initGrid() {
    for (let x = 0; x < this.cellCount; x++) {
      for (let y = 0; y < this.cellCount; y++) {
        let cell = new Cell(x, y, this.cellWidth, this.margin);
        this.grid.push(cell);
      }
    }
    this.initCells();
  }

  initCells() {
    this.endCellsIndex = []; // Initialize endCellsIndex here

    // Get random cell in grid for start cell
    this.startIndex = Math.floor(Math.random() * this.grid.length);
    this.startCell = this.grid[this.startIndex];
    this.startCell.color = color(255, 255, 100);
    this.endCellsIndex.push(this.startIndex);
    let numTargets;
    do {
      numTargets = Math.floor(Math.random() * (11 - 6 + 1)) + 6;
    } while (numTargets % 2 !== 0);

    // Initialize the target cells
    for (let i = 0; i <= numTargets; i++) {
      let targetIndex;
      do {
        targetIndex = Math.floor(Math.random() * this.grid.length);
      } while (this.endCellsIndex.includes(targetIndex));
      if (i === 0) {
        this.grid[targetIndex].color = color(255, 100, 255);
      } else {
        this.grid[targetIndex].color = color(100, 100, 255);
      }
      this.endCellsIndex.push(targetIndex);
    }
  }

  handleMouse() {
    if (!touch) {
      for (let i = 0; i < this.endCellsIndex.length; i++) {
        let index = this.endCellsIndex[i];
        let cell = this.grid[index];

        if (cell.isMouseOver()) {
          // If the mouse is over the first end cell and it is clicked, start recording path data
          if (i === 0 && mouseIsPressed && !recording) {
            this.startRecording();

            // Change the color to green
            cell.color = color(100, 255, 100);
            recording = true;
          }

          // If the mouse is over the second end cell (red cell) and it is clicked, end recording and send data
          if (i === 1 && mouseIsPressed && recording) {
            this.stopRecording();
          }
        }
      }
    }
  }

  step() {
    if (this.endCellsIndex.length >= 2) {
      // Remove the first two target cells
      let removedIndex1 = this.endCellsIndex.shift();
      let removedIndex2 = this.endCellsIndex.shift();

      // Make sure removedIndex1 and removedIndex2 are valid indices
      if (
        removedIndex1 !== undefined &&
        removedIndex2 !== undefined
      ) {
        // Revert start / target cells to background color
        this.grid[removedIndex1].color = color(51, 51, 51, 50);
        this.grid[removedIndex2].color = color(51, 51, 51, 50);

        // Check if endCellsIndex has at least two elements before accessing them
        if (this.endCellsIndex.length >= 2) {
          let newStart = this.endCellsIndex[0];
          let newTarget = this.endCellsIndex[1];
          this.grid[newStart].color = color(255, 255, 100);
          this.grid[newTarget].color = color(255, 100, 1000);
        } else {
          // Restart the game by initializing cells
          this.initCells();
        }
      }
    } else {
      // Restart the game by initializing cells
      this.initCells();
    }
  }

  startRecording() {
    clearData();
    myChart.update();
    recording = true;

    this.timeoutId = setTimeout(() => {
      this.stopRecording(true);
    }, maxDuration);
  }

  stopRecording(timeout = false) {
    if (!timeout) {
      clearTimeout(this.timeoutId);
      sendData(data);
      // Adds all datapoints to the frontend graph
      data['mouse-array'].forEach((entry) => {
        addData(entry.x * windowWidth, entry.y * windowHeight);
      });
      myChart.update();
    }
    this.step();
    recording = false;

    // reset data array
    data = {
      'mouse-array': [],
    };
  }

  show() {
    this.handleMouse();
    for (let cell of this.grid) {
      cell.show();
    }
  }
}

let grid;
let canvasBoundingBox;
let myCanvasElement;
let relPos;
let startTime;

function getMouseCanvasPosition(rect, mouseX, mouseY) {
  const mouseXCanvas = mouseX - rect.left;
  const mouseYCanvas = mouseY - rect.top;
  return { x: mouseXCanvas, y: mouseYCanvas };
}

let data = {
  'mouse-array': [],
};

document.addEventListener('DOMContentLoaded', function () {
  docReady = true;
});

document.addEventListener('pointermove', function (event) {
  if (docReady) {
    const mouseX = event.clientX;
    const mouseY = event.clientY;

    myCanvasElement = document.getElementById('myCanvas');
    canvasBoundingBox = myCanvasElement.getBoundingClientRect();
    relPos = getMouseCanvasPosition(
      canvasBoundingBox,
      mouseX,
      mouseY
    );

    if (recording) {
      // First entry set time to 0
      if (data['mouse-array'].length === 0) {
        let startTime = 0;
        data['mouse-array'].push({
          x: relPos.x / windowWidth,
          y: relPos.y / windowHeight,
          time: startTime,
        });
        startTime = performance.now();
      } else {
        data['mouse-array'].push({
          x: relPos.x / windowWidth,
          y: relPos.y / windowHeight,
          time:
            data['mouse-array'][data['mouse-array'].length - 1] -
            startTime,
        });
        startTime = performance.now();
      }
    }
  }
});

function setup() {
  let cellCount = 16;
  let margin = 5;
  let canvasSize = 800;
  let cellWidth = canvasSize / cellCount;

  let w = cellCount * (cellWidth + margin);
  let h = cellCount * (cellWidth + margin);

  let canvas = createCanvas(w, h);
  canvas.id('myCanvas'); // Set an ID for the canvas
  grid = new Grid(width, cellCount, cellWidth, margin);
  myCanvasElement = document.getElementById('myCanvas');
  canvasBoundingBox = myCanvasElement.getBoundingClientRect();

  // Append the canvas to the 'dataMap' div
  let dataMapDiv = document.getElementById('dataMap');
  dataMapDiv.appendChild(document.getElementById('myCanvas'));
}

function draw() {
  grid.show();
}

const ctx = document.getElementById('myChart');

// Set the canvas size
ctx.width = 400;
ctx.height = 400;

const myChart = new Chart(ctx, {
  type: 'line',
  data: {
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
        beginAtZero: true, // Start x-axis at zero
        suggestedMax: 800,
        ticks: {
          display: false, // Hide x-axis labels
        },
      },
      y: {
        beginAtZero: true, // Start y-axis at zero
        suggestedMax: 800,
        reverse: true, // Invert y-axis
        ticks: {
          display: false, // Hide y-axis labels
        },
      },
    },
  },
});

// Function to add data to the chart
function addData(x, y) {
  myChart.data.datasets[0].data.push({ x: x, y: y });
}

function clearData() {
  myChart.data.datasets[0].data = [];
  console.log('clearing data', (myChart.data.datasets[0].data = []));
  myChart.update();
}

document.body.addEventListener('click', function () {
  if (touch) {
    document.querySelector(
      '.modal-content'
    ).innerHTML = `<h1>This application is for cursor movement path data only! Touchscreens disable interaction.</h1>
    <h1><a href="https://github.com/sockheadrps/AiMouseMovement">View the code on github</a></h1>`;
    document.querySelector('.modal-content').style.backgroundColor =
      'rgba(255, 102, 102, 0.8)';
  } else {
    closeModal()
  }
});

function closeModal() {
  document.getElementById('myModal').style.display = 'none';
}

function openModel() {
  document.getElementById('myModal').style.display = 'block';
}
