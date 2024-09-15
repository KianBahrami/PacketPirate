/*
    This file provides the function used for applying the filter.
*/

function isValidIP(ip) {
    // IPv4 regex
    const ipv4Regex = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;

    // IPv6 regex (simplified, doesn't cover all edge cases)
    const ipv6Regex = /^(?:[A-F0-9]{1,4}:){7}[A-F0-9]{1,4}$/i;

    return ipv4Regex.test(ip) || ipv6Regex.test(ip);
}

export function FilterButtonFunction(filterOptions, filterPanel) {
    const networkLayerProtocol = document.getElementById('networkLayerProtocolDropdown').value;
    const transportLayerProtocol = document.getElementById('transportLayerProtocolDropdown').value;
    const srcIp = document.getElementById('filter-srcip').value;
    const destIp = document.getElementById('filter-destip').value;
    const minPayloadSize = document.getElementById('filter-minpayloadsize').value;
    // check for correct inputs
    if (srcIp != "" && !isValidIP(srcIp)) {
        alert("Please enter a valid source IP-Address.");
        return
    }
    if (destIp != "" && !isValidIP(destIp)) {
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
    return filterOptions
}
