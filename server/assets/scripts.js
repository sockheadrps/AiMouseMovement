const firstCircle = document.getElementById('firstCircle');
const secondCircle = document.getElementById('secondCircle');
let cursorInFirstCircle = false;
let cursorPath = [];

let start
let end
const pollFreq = 8

const url = '/add_data';

const windowWidth = window.innerWidth;
const windowHeight = window.innerHeight;

function getRandomPosition() {
    const x = Math.random() * (window.innerWidth - 100);
    const y = Math.random() * (window.innerHeight - 100);
    return { x, y };
}

function sendData(data) {
    let json_data = JSON.stringify({
        "window-height": windowHeight,
        "window-width" : window.innerWidth,
        "mouse-array": data['mouse-array'] 
    });

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

function moveCirclesRandomly() {
    const firstPosition = getRandomPosition();
    const secondPosition = getRandomPosition();

    firstCircle.style.top = `${firstPosition.y}px`;
    firstCircle.style.left = `${firstPosition.x}px`;

    secondCircle.style.top = `${secondPosition.y}px`;
    secondCircle.style.left = `${secondPosition.x}px`;

    cursorInFirstCircle = false;
    cursorPath = [];
    firstCircle.style.backgroundColor = "blue";
}

moveCirclesRandomly();

let colorChangeTimeout;
let yellowTimeout;
firstCircle.addEventListener('mouseenter', function () {
    colorChangeTimeout = setTimeout(function () {
        firstCircle.style.backgroundColor = "yellow";
        yellowTimeout = setTimeout(function () {
            firstCircle.style.backgroundColor = "green";
            start = new Date().getTime();
        }, 300);
    }, 300);
});

firstCircle.addEventListener('mouseleave', function () {
    clearTimeout(colorChangeTimeout);
    clearTimeout(yellowTimeout);
    if (firstCircle.style.backgroundColor !== 'green') {
        firstCircle.style.backgroundColor = 'blue';
        cursorPath = [];
    }
});

secondCircle.addEventListener('click', function (event) {
    const rect = event.target.getBoundingClientRect();
    const mouseX = event.clientX - rect.left;
    const mouseY = event.clientY - rect.top;
    if (mouseX >= 0 && mouseX <= 100 && mouseY >= 0 && mouseY <= 100) {
        cursorPath.push({ 
            x: event.clientX / windowWidth,
            y: event.clientY / windowHeight,
            time: new Date().getTime()
        });
        const data = {
            "mouse-array": cursorPath
        }
        sendData(data)
        moveCirclesRandomly();
    }
});

document.addEventListener('mousemove', function (event) {

    if (cursorInFirstCircle && (firstCircle.style.backgroundColor === 'green')) {
        if (new Date().getTime() - start > pollFreq) {
            const mouseX = event.clientX;
            const mouseY = event.clientY;
    
            cursorPath.push({ 
                x: mouseX / windowWidth,
                y: mouseY / windowHeight,
                time: new Date().getTime()
            });
            start = new Date().getTime()
        }
        
        
    } else {
        cursorInFirstCircle = true;
    }

});
