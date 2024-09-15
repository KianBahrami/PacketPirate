/*
  This file only assignes the eventhandlers to each button and initializes the website.
  Besides, it provides the global variables that are used accross functions.
*/

import { ConnectWebSocket } from "./websocketHandler.js";
import { FilterButtonFunction } from "./filterHandler.js";

let socket;
let filterOptions;  // used in ConnectWebsSocket
let selectedInterface = "";
const startButton = document.getElementById("startButton");
const stopButton = document.getElementById("stopButton");
const packetList = document.getElementById("packetList");
const clearButton = document.getElementById("clearButton");
const interfaceDropdown = document.getElementById("interfaceDropdown");
const filterButton = document.getElementById("filterButton");
const filterPanel = document.getElementById("filterPanel");
const applyFilterButton = document.getElementById("applyFilterButton");

// add eventlistener for loading the webpage
document.addEventListener("DOMContentLoaded", function () {
  fetch("http://localhost:8080/get-interfaces") // fetch interface api
    .then((response) => response.json()) // bring response into json
    .then((data) => {
      data.interfaces.forEach((iface) => {
        // loop over received interfaces
        // append name to dropdown
        const interfaceOption = document.createElement("option");
        interfaceOption.value = iface.name;
        interfaceOption.textContent = iface.name;
        interfaceOption.title = `${iface.description || "No description available"
          }`;
        interfaceDropdown.appendChild(interfaceOption);
      });
    })
    .catch((error) => console.error("Error:", error));
});

// connect to websocket if start button is clicked
startButton.addEventListener("click", function () {
  console.log("Start button clicked");
  if (interfaceDropdown.value == "") {
    alert("Please select an interface!");
    console.log("No interface selected, aborting")
    return
  }
  socket = ConnectWebSocket(selectedInterface, filterOptions, socket, startButton, stopButton);
  this.disabled = true;
  stopButton.disabled = false;
  interfaceDropdown.disabled = true;
});

// close websocket if stop button is clicked
stopButton.addEventListener("click", function () {
  console.log("Stop button clicked");
  if (socket) {
    socket.send("stop"); // signal backend to stop sending
    socket.close();
  }
  this.disabled = true;
  startButton.disabled = false;
  interfaceDropdown.disabled = false;
  bpsDisplay.innerHTML = `Datarate: 0 KB/s`;
});

// get selected interface
interfaceDropdown.addEventListener("change", function () {
  selectedInterface = this.value;
  console.log("Selected interface:", selectedInterface);
});

// clear packet list
clearButton.addEventListener("click", function () {
  packetList.innerHTML = "";
});

// open filter panel
filterButton.addEventListener("click", () => {
  filterPanel.classList.add("open");
});

// add handling og filter inputs
applyFilterButton.addEventListener("click", () => {
  filterOptions = FilterButtonFunction(filterOptions, filterPanel);
});

// also close filter panel if clicked outside
document.addEventListener("click", (event) => {
  if (!filterPanel.contains(event.target) && event.target !== filterButton) {
    filterPanel.classList.remove("open");
  }
});
