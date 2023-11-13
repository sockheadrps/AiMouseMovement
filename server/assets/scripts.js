let recordingInterval;
let maxDuration = 3000;
let pollFreq = 8;
const windowWidth = 800;
const windowHeight = 800;
const url = '/add_data';
let recording = false

let pointsElm = document.querySelector('#points')
let points = 0

function sendData(data) {
    points++;
    pointsElm.textContent = points
    let json_data = JSON.stringify({
        "window-height": windowHeight,
        "window-width" : windowWidth,
        "mouse-array": data['mouse-array'] 
    });
    console.log(json_data)

    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: json_data
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Response data:', data);
        })
        .catch(error => {
            console.error('Error:', error);
        });
}


class Cell {
    constructor(x, y, width) {
        this.x = x;
        this.y = y;
        this.width = width;
        this.filled = false;
        this.color = color(211);
    }

    isMouseOver() {
        return mouseX > this.x * this.width && mouseX < (this.x + 1) * this.width &&
               mouseY > this.y * this.width && mouseY < (this.y + 1) * this.width;
    }

    show() {
        let px = this.x * this.width;
        let py = this.y * this.width;
        fill(this.color);
        rect(px, py, this.width, this.width);
        stroke(0);
        noFill();
        rect(px, py, this.width, this.width);
    }
    
}


class Grid {
    constructor(canvasSize, cellCount) {
        this.canvasSize = canvasSize;
        this.cellCount = cellCount;
        this.cellWidth = this.canvasSize / this.cellCount;
        this.grid = [];
        this.initGrid();
        this.recordStartTime = 0;
        this.isRecording = false;
        this.data = {
                'mouse-array': [],
            }
    }

    initGrid() {
        for (let x = 0; x < this.cellCount; x++) {
            for (let y = 0; y < this.cellCount; y++) {
                let cell = new Cell(x, y, this.cellWidth);
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
        this.startCell.color = color(255, 255, 0);
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
                this.grid[targetIndex].color = color(255, 0, 0);
            } else {
                this.grid[targetIndex].color = color(0, 0, 255);
            }
            this.endCellsIndex.push(targetIndex);
        }
    }


    handleMouse() {
        for (let i = 0; i < this.endCellsIndex.length; i++) {
            let index = this.endCellsIndex[i];
            let cell = this.grid[index];
    
            if (cell.isMouseOver()) {
                // If the mouse is over the first end cell and it is clicked, start recording path data
                if (i === 0 && mouseIsPressed) {
                    // Change the color to green
                    cell.color = color(0, 255, 0);
                    this.startRecording()
                }
    
                // If the mouse is over the second end cell (red cell) and it is clicked, end recording and send data
                if (i === 1 && mouseIsPressed && recording) {
                    this.stopRecording()
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
            if (removedIndex1 !== undefined && removedIndex2 !== undefined) {
                // Revert start / target cells to background color
                this.grid[removedIndex1].color = color(211); 
                this.grid[removedIndex2].color = color(211); 
    
                // Check if endCellsIndex has at least two elements before accessing them
                if (this.endCellsIndex.length >= 2) {
                    let newStart = this.endCellsIndex[0];
                    let newTarget = this.endCellsIndex[1];
                    this.grid[newStart].color = color(255, 255, 0);
                    this.grid[newTarget].color = color(255, 0, 0);
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
        if (!this.recordInterval) {
            this.recordStartTime = millis();  
            // push initial position and set time reference
            this.data['mouse-array'].push({ x: relPos.x / windowWidth, y: relPos.y / windowHeight, time: 0 });
            clearData()
            recording = true;
            
            
            // Begin polling mouse position
            this.recordInterval = setInterval(() => {
                let currentTime = millis();
                this.data['mouse-array'].push({ x: relPos.x / windowWidth, y: relPos.y / windowHeight, time: currentTime  - this.recordStartTime });
                

    
                // Check if the recording duration exceeds the maximum allowed duration
                if (currentTime - this.recordStartTime > maxDuration) {
                    this.stopRecording(true);
                }
            }, pollFreq);
    
            recordingInterval = this.recordInterval; // Store the interval reference for later use
        }
    }
    
    stopRecording(timeout = false) {
        recording = false
        if (recordingInterval) {
            clearInterval(recordingInterval);
            this.recordInterval = undefined;
            if (!timeout) {
                sendData(this.data);
            }
        }
        this.data = {
            'mouse-array': [],
        };
        this.step();
    }
        

    show() {
        this.handleMouse();
        for (let cell of this.grid) {
            cell.show();
        }
    }
}


let grid;
let canvasBoundingBox
let myCanvasElement
let relPos

function getMouseCanvasPosition(rect, mouseX, mouseY) {
    const mouseXCanvas = mouseX - rect.left;
    const mouseYCanvas = Math.abs(mouseY - rect.bottom);
    return { x: mouseXCanvas, y: mouseYCanvas };
}

document.addEventListener('mousemove', function(event) {
    const mouseX = event.clientX;
    const mouseY = event.clientY;
    myCanvasElement = document.getElementById('myCanvas');

    canvasBoundingBox = myCanvasElement.getBoundingClientRect();
    relPos = getMouseCanvasPosition(canvasBoundingBox, mouseX, mouseY);
    if (recording) {
        addData(relPos.x, relPos.y)
    }

});



function setup() {
    let canvas = createCanvas(800, 800);
    canvas.id('myCanvas'); // Set an ID for the canvas
    let cellCount = 16;
    grid = new Grid(width, cellCount);
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
        labels: ['Point A', 'Point B', 'Point C', 'Point D', 'Point E', 'Point F'],
        datasets: [{
            label: 'Data Series 1',
            borderColor: 'rgba(75, 192, 192, 1)',
            data: [],
            fill: false,
            borderWidth: 2
        }]
    },
    options: {
        scales: {
            x: {
                type: 'linear',
                position: 'bottom',
                beginAtZero: false,
                suggestedMin: 0, // Set the suggestedMin for x-axis to -800
                suggestedMax: 800,
                ticks: {
                    display: false // Hide x-axis labels
                },
            },
            y: {
                beginAtZero: false,
                suggestedMin: 0, // Set the suggestedMin for y-axis to -800
                suggestedMax: 800,
                ticks: {
                    display: false // Hide x-axis labels
                },
            }
        }
    }
});


// Function to add data to the chart
function addData(x, y) {
    myChart.data.datasets[0].data.push({ x: x, y: y });
    myChart.update();
}

function clearData() {
    myChart.data.datasets[0].data = [];
    myChart.update();
}


document.body.addEventListener('click', function () {
    closeModal();
});

function closeModal() {
    document.getElementById('myModal').style.display = 'none';
}
