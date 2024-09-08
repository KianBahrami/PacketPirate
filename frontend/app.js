import { PacketToHTML } from './packetToHTML.js';

let socket;
let selectedInterface = ''
const startButton = document.getElementById('startButton');
const stopButton = document.getElementById('stopButton');
const packetList = document.getElementById('packetList');
const clearButton = document.getElementById('clearButton');
clearButton.disabled = false;
const interfaceDropdown = document.getElementById('interfaceDropdown');

// add eventlistener for loading the webpage
document.addEventListener('DOMContentLoaded', function() {
    fetch('http://192.168.0.14:8080/get-interfaces')    // fetch interface api
        .then(response => response.json())  // bring response into json
        .then(data => {
            data.interfaces.forEach(iface => {  // loop over received interfaces
                // append name to dropdown
                const interfaceOption = document.createElement('option');
                interfaceOption.value = iface.name;
                interfaceOption.textContent = iface.name
                interfaceDropdown.appendChild(interfaceOption);
            });
        })
        .catch(error => console.error('Error:', error));
});

// initiates the websocket
function connectWebSocket() {
    console.log("Attempting to connect WebSocket...");
    socket = new WebSocket('ws://localhost:8080/ws');

    socket.onopen = function(e) {
        console.log("WebSocket connection established");
        // signal backend to start sending packages occuring at the specified interface
        socket.send(JSON.stringify({command: "start", interface: selectedInterface}));
    };

    // if a message arrives parse the packet into json and add it to list
    socket.onmessage = function(event) {
        console.log("Received message:", event.data);
        const packet = JSON.parse(event.data);
        addPacketToList(packet);
    };

    // close connection if desired and reset buttons
    socket.onclose = function(event) {
        console.log("WebSocket connection closed", event);
        startButton.disabled = false;
        stopButton.disabled = true;
    };

    // log possible errors
    socket.onerror = function(error) {
        console.log("WebSocket error:", error);
    };
}

// adds package to list
function addPacketToList(packet) {
    const packetElement = document.createElement('div');
    packetElement.className = 'packet-item';
    packetElement.innerHTML = PacketToHTML(packet);
    
    // Add click event to toggle details
    packetElement.querySelector('.packet-summary').addEventListener('click', function() {
        const details = packetElement.querySelector('.packet-details');
        details.style.display = details.style.display === 'none' ? 'block' : 'none';
    });
    
    packetList.appendChild(packetElement);
    packetList.scrollTop = packetList.scrollHeight;
}

// connect to websocket if start button is clicked
startButton.addEventListener('click', function() {
    console.log("Start button clicked");
    connectWebSocket();
    this.disabled = true;
    stopButton.disabled = false;
    interfaceDropdown.disabled = true;
});

// close websocket if stop button is clicked
stopButton.addEventListener('click', function() {
    console.log("Stop button clicked");
    if (socket) {
        socket.send("stop");    // signal backend to stop sending
        socket.close();
    }
    this.disabled = true;
    startButton.disabled = false;
    interfaceDropdown.disabled = false;
});

// get selected interface
interfaceDropdown.addEventListener('change', function() {
    selectedInterface = this.value;
    console.log("Selected interface:", selectedInterface);
});

clearButton.addEventListener('click', function() {
    console.log("Clear button clicked");
    packetList.innerHTML = '';
});