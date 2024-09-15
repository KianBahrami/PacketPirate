/*
    This file provides the establishment of a websocket connection and its handler to the backend.
*/

import { PacketToHTML } from "./packetToHTML.js";

// initiates the websocket
export function ConnectWebSocket(selectedInterface, filterOptions, socket, startButton, stopButton) {
    console.log("Attempting to connect to WebSocket...");
    socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = function (e) {
        console.log("WebSocket connection established");
        // signal backend to start sending packages occuring at the specified interface
        socket.send(
            JSON.stringify({ command: "start", interface: selectedInterface, filteroptions: filterOptions })
        );
    };

    // if a message arrives parse the packet into json and add it to list
    socket.onmessage = function (event) {
        console.log("Received message:", event.data);
        const data = JSON.parse(event.data);
        if (data.bps !== undefined) {
            updateBPSDisplay(data);
        } else {
            addPacketToList(data);
        }
    };

    // close connection if desired and reset buttons
    socket.onclose = function (event) {
        console.log("WebSocket connection closed", event);
        startButton.disabled = false;
        stopButton.disabled = true;
    };

    // log possible errors
    socket.onerror = function (error) {
        console.log("WebSocket error:", error);
    };

    return socket;
}

// adds package to list
function addPacketToList(packet) {
    const packetElement = document.createElement("div");
    packetElement.className = "packet-item";
    packetElement.innerHTML = PacketToHTML(packet);

    // Add click event to toggle details
    packetElement
        .querySelector(".packet-summary")
        .addEventListener("click", function () {
            const details = packetElement.querySelector(".packet-details");
            details.style.display =
                details.style.display === "none" ? "block" : "none";
        });

    packetList.appendChild(packetElement);
    packetList.scrollTop = packetList.scrollHeight;
}

function updateBPSDisplay(data) {
    bpsDisplay.innerHTML = `
        Datarate: ${(data.bps / 1024).toFixed(2)} KB/s
    `;
}