export function PacketToHTML(packet) {
  return `
            <div class="packet-summary">
                ${packet.networklayer.src} → ${packet.networklayer.dest} (${packet.time
    })
            </div>
            <div class="packet-details" style="display: none;">
                ${createLayerDropdown(
      "Link Layer",
      `
                    <p><strong>Protocol:</strong> ${packet.linklayer.protocol}</p>
                    <p><strong>Source MAC:</strong> ${packet.linklayer.src}</p>
                    <p><strong>Destination MAC:</strong> ${packet.linklayer.dest}</p>
                `
    )}

                ${createLayerDropdown(
      "Network Layer",
      `
                    <p><strong>Protocol:</strong> ${packet.networklayer.protocol}</p>
                    <p><strong>Source IP:</strong> ${packet.networklayer.src}</p>
                    <p><strong>Destination IP:</strong> ${packet.networklayer.dest}</p>
                    <p><strong>TTL/Hop Limit:</strong> ${packet.networklayer.ttl}</p>
                `
    )}

                ${createLayerDropdown(
      "Transport Layer",
      `
                    <p><strong>Protocol:</strong> ${packet.transportlayer.protocol
      }</p>
                    <p><strong>Source Port:</strong> ${packet.transportlayer.src
      }</p>
                    <p><strong>Destination Port:</strong> ${packet.transportlayer.dest
      }</p>
                    ${packet.transportlayer.protocol === "TCP"
        ? `
                        <p><strong>Flags:</strong> ${packet.transportlayer.tcpflags}</p>
                        <p><strong>Sequence Number:</strong> ${packet.transportlayer.tcpseq}</p>
                        <p><strong>Acknowledgment Number:</strong> ${packet.transportlayer.tcpack}</p>
                        <p><strong>Window Size:</strong> ${packet.transportlayer.tcpwindow}</p>
                    `
        : ""
      }
                `
    )}

                ${createLayerDropdown(
      "Application Layer",
      `
                    <p><strong>Protocol:</strong> ${packet.applicationlayer.protocol
      }</p>
                    <p><strong>Payload Size:</strong> ${packet.applicationlayer.payloadsize
      } bytes</p>
                    ${packet.applicationlayer.protocol === "HTTP"
        ? `
                        <p><strong>HTTP Method:</strong> ${packet.applicationlayer.httpmethod}</p>
                        <p><strong>HTTP URL:</strong> ${packet.applicationlayer.httpurl}</p>
                        <p><strong>HTTP Version:</strong> ${packet.applicationlayer.httpversion}</p>
                    `
        : ""
      }
                `
    )}

                ${createLayerDropdown(
      "Raw Packet",
      `
                    <pre>${packet.raw}</pre>
                `
    )}
            </div>
        </div>
    `;
}

function createLayerDropdown(layerName, content) {
  return `
        <div class="layer-dropdown">
            <div class="layer-summary" onclick="toggleLayer(this)">
                <span class="dropdown-arrow"></span> ${layerName}
            </div>
            <div class="layer-details" style="display: none;">
                ${content}
            </div>
        </div>
    `;
}

// toggles a network layers's details
window.toggleLayer = function (element) {
  const details = element.nextElementSibling;
  element.classList.toggle("open");
  if (details.style.display === "none") {
    details.style.display = "block";
  } else {
    details.style.display = "none";
  }
};
