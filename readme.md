# PacketPirate 🏴‍☠️

PacketPirate is a network packet capture and analysis tool built with Go and JavaScript. It provides a web-based interface for real-time packet inspection and filtering.

## What PacketPirate Does

PacketPirate allows users to:

1. Capture network packets in real-time from selected network interfaces.
2. Filter packets based on various criteria:
3. View detailed information about captured packets, including:
   - Source and destination MAC addresses
   - Source and destination IP addresses
   - Protocol information
   - Packet payload
4. Start and stop packet capture sessions dynamically through the web interface.

The tool uses Go for backend packet capturing and processing, leveraging the `gopacket` library for low-level packet handling. The frontend is built with JavaScript, providing an interactive user interface for controlling the capture process and viewing packet data.

PacketPirate caters to individuals seeking a straightforward, customizable network analysis tool that can be easily adapted to meet specific requirements.

*Note*: Running PacketPirate may require administrative privileges due to its need to access network interfaces at a low level.
## Getting Started

### Dependencies

* go version go1.23.0 windows/amd64
    * needed go modules in `./backend/go.mod`
* node version v16.13.2

### Installing

* Clone this repository with `git clone https://github.com/KianBahrami/PacketPirate.git`

### Executing program

* Running `./start.ps1` will start both the backend and frontend
* open the frontend on `http://localhost:8000/` or the address shown in the console

### Help

Feel free to open any issues.

### ToDos
* more complex filtering options: include min/max payload size, text in payload,...
* incoming/outgoing data rate including graphs
* better GUI 

## Authors

Contributors names and contact info

- Kian Bahrami, kian.bahrami@rwth-aachen.de
