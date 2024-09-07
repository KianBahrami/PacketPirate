# Navigate to the project root directory
# Adjust this path if the script is not in the project root
Set-Location -Path $PSScriptRoot

# Function to get IP address
function Get-IPAddress {
    $ip = (Get-NetIPAddress -AddressFamily IPv4 -PrefixOrigin Dhcp).IPAddress
    if (-not $ip) {
        $ip = "localhost"
    }
    return $ip
}

# Get the IP address
$IP_ADDRESS = Get-IPAddress

# Start the backend
Write-Host "Starting the backend..."
$backendJob = Start-Job -ScriptBlock {
    Set-Location -Path "$using:PSScriptRoot\backend"
    go run .
}

# Wait a moment to ensure the backend has started
Start-Sleep -Seconds 2

# Start the frontend
Write-Host "Starting the frontend..."
$frontendJob = Start-Job -ScriptBlock {
    Set-Location -Path "$using:PSScriptRoot\frontend"
    http-server -p 8000
}

Write-Host "Both servers are now running."
Write-Host "Backend is running at: http://$($IP_ADDRESS):8080"
Write-Host "Frontend is running at: http://$($IP_ADDRESS):8000"
Write-Host "Press Ctrl+C to stop both servers."

try {
    # Wait for jobs to complete (which they won't unless there's an error)
    Wait-Job -Job $backendJob, $frontendJob
} finally {
    # This block will execute when Ctrl+C is pressed
    Write-Host "Stopping servers..."
    Stop-Job -Job $backendJob, $frontendJob
    Remove-Job -Job $backendJob, $frontendJob
    Write-Host "Servers stopped."
}