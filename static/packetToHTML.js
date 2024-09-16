export function PacketToHTML(packet) {
  // Determine if it's an ARP packet
  const isARP = (packet.arplayer.dstip != "None");

  // Create the packet summary
  let summary;
  if (isARP) {
    summary = `${packet.arplayer.srcip} → ${packet.arplayer.dstip} (ARP ${packet.arplayer.operation === 1 ? 'Request' : 'Reply'}) (${packet.time})`;
  } else {
    summary = `${packet.networklayer.src} → ${packet.networklayer.dest} (${packet.time})`;
  }

  return `
    <div class="packet-summary">
        ${summary}
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

        ${isARP ? createARPLayerDropdown(packet.arplayer) : createNetworkLayerDropdown(packet.networklayer)}

        ${!isARP ? createTransportLayerDropdown(packet.transportlayer) : ''}

        ${!isARP ? createApplicationLayerDropdown(packet.applicationlayer) : ''}

        ${createLayerDropdown(
    "Raw Packet",
    `
            <pre>${packet.raw}</pre>
          `
  )}
    </div>
  `;
}

function createARPLayerDropdown(arplayer) {
  return createLayerDropdown(
    "ARP Layer",
    `
      <p><strong>Operation:</strong> ${arplayer.operation === 1 ? 'Request (1)' : 'Reply (2)'}</p>
      <p><strong>Sender MAC:</strong> ${arplayer.srcmac}</p>
      <p><strong>Sender IP:</strong> ${arplayer.srcip}</p>
      <p><strong>Target MAC:</strong> ${arplayer.dstmac}</p>
      <p><strong>Target IP:</strong> ${arplayer.dstip}</p>
    `
  );
}

function createNetworkLayerDropdown(networklayer) {
  return createLayerDropdown(
    "Network Layer",
    `
      <p><strong>Protocol:</strong> ${networklayer.protocol}</p>
      <p><strong>Source IP:</strong> ${networklayer.src}</p>
      <p><strong>Destination IP:</strong> ${networklayer.dest}</p>
      <p><strong>TTL/Hop Limit:</strong> ${networklayer.ttl}</p>
    `
  );
}

function createTransportLayerDropdown(transportlayer) {
  return createLayerDropdown(
    "Transport Layer",
    `
      <p><strong>Protocol:</strong> ${transportlayer.protocol}</p>
      <p><strong>Source Port:</strong> ${transportlayer.src}</p>
      <p><strong>Destination Port:</strong> ${transportlayer.dest}</p>
      ${transportlayer.protocol === "TCP"
      ? `
          <p><strong>Flags:</strong> ${transportlayer.tcpflags}</p>
          <p><strong>Sequence Number:</strong> ${transportlayer.tcpseq}</p>
          <p><strong>Acknowledgment Number:</strong> ${transportlayer.tcpack}</p>
          <p><strong>Window Size:</strong> ${transportlayer.tcpwindow}</p>
        `
      : ""
    }
    `
  );
}

function createApplicationLayerDropdown(applicationlayer) {
  if (applicationlayer.payloadsize == '0') {
    return "";
  }
  return createLayerDropdown(
    "Application Layer",
    `
      <p><strong>Protocol:</strong> ${applicationlayer.protocol}</p>
      <p><strong>Payload Size:</strong> ${applicationlayer.payloadsize} bytes</p>
      ${applicationlayer.protocol === "HTTP"
      ? `
          <p><strong>HTTP Method:</strong> ${applicationlayer.httpmethod}</p>
          <p><strong>HTTP URL:</strong> ${applicationlayer.httpurl}</p>
          <p><strong>HTTP Version:</strong> ${applicationlayer.httpversion}</p>
        `
      : ""
    }
    `
  );
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

// toggles a network layer's details
window.toggleLayer = function (element) {
  const details = element.nextElementSibling;
  element.classList.toggle("open");
  if (details.style.display === "none") {
    details.style.display = "block";
  } else {
    details.style.display = "none";
  }
};

