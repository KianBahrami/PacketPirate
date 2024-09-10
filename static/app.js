import { PacketToHTML } from "./packetToHTML.js";

let socket;
let filterOptions;
let selectedInterface = "";
const startButton = document.getElementById("startButton");
const stopButton = document.getElementById("stopButton");
const packetList = document.getElementById("packetList");
const clearButton = document.getElementById("clearButton");
const interfaceDropdown = document.getElementById("interfaceDropdown");
const filterButton = document.getElementById("filterButton");
const filterPanel = document.getElementById("filterPanel");
const applyFilterButton = document.getElementById("applyFilterButton");
const filterMinPayloadSize = document.getElementById("filter-minpayloadsize");

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

// initiates the websocket
function connectWebSocket() {
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
    const packet = JSON.parse(event.data);
    addPacketToList(packet);
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

// connect to websocket if start button is clicked
startButton.addEventListener("click", function () {
  console.log("Start button clicked");
  if (interfaceDropdown.value == "") {
    alert("Please select an interface!");
    console.log("No interface selected, aborting")
    return
  }
  connectWebSocket();
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

function isValidIP(ip) {
  // IPv4 regex
  const ipv4Regex = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;

  // IPv6 regex (simplified, doesn't cover all edge cases)
  const ipv6Regex = /^(?:[A-F0-9]{1,4}:){7}[A-F0-9]{1,4}$/i;

  return ipv4Regex.test(ip) || ipv6Regex.test(ip);
}

// close filter panel and save filtering options
applyFilterButton.addEventListener("click", () => {
  const networkLayerProtocol = document.getElementById('networkLayerProtocolDropdown').value;
  const transportLayerProtocol = document.getElementById('transportLayerProtocolDropdown').value;
  const srcIp = document.getElementById('filter-srcip').value;
  const destIp = document.getElementById('filter-destip').value;
  const minPayloadSize = document.getElementById('filter-minpayloadsize').value;
  // check for correct inputs
  if  (srcIp != "" && !isValidIP(srcIp)) {
    alert("Please enter a valid source IP-Address.");
    return
  }
  if  (destIp != "" && !isValidIP(destIp)) {
    alert("Please enter a valid destination IP-Address.");
    return
  }
  if (!Number.isInteger(Number(minPayloadSize))) {
    alert("Please enter an integer number as minimal payload size.");
    return
  }
  filterOptions = {
    "networkLayerProtocol": networkLayerProtocol,
    "transportLayerProtocol": transportLayerProtocol,
    "srcIp": srcIp,
    "destIp": destIp,
    "minPayloadSize": minPayloadSize,
  }
  console.log("Filter applied with", networkLayerProtocol, transportLayerProtocol, srcIp, destIp, minPayloadSize);
  filterPanel.classList.remove("open");
});

// also close filter panel if clicked outside
document.addEventListener("click", (event) => {
  if (!filterPanel.contains(event.target) && event.target !== filterButton) {
    filterPanel.classList.remove("open");
  }
});
