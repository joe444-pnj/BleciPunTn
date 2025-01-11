# BleciPunTn

BleciPunTn is a powerful Open Source Intelligence (OSINT) tool written in Go that provides detailed information about IP addresses, including geolocation, hostname, organization, and more. It also features port scanning capabilities for open ports and service identification.

## Features

- **IP Information Retrieval**: Fetches details like hostname, location, timezone, organization, and more for given IP addresses.
- **Reverse DNS Lookup**: Performs reverse DNS lookup to find domain names associated with IP addresses.
- **Port Scanning**: Scans a range of ports to identify open ports and the services running on them.
- **Banner Grabbing**: Attempts to retrieve banners from open ports to identify services.
- **Colored Output**: Uses colors to differentiate various pieces of information for better readability.
- **Logging**: Saves scan results to a `logs` directory for future reference.
- **Interactive Mode**: Allows users to enter IP addresses interactively, with support for multiple commands like `clear` and `help`.

## Getting Started

### Prerequisites

- ![Static Badge](https://img.shields.io/badge/Go-1.23.4-black)

- An internet connection to fetch IP information

### Installation

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/joe444-pnj/BleciPunTn.git
    cd BleciPunTn/cmd
    ```

2. **Install Dependencies**:
    ```bash
    go mod tidy
    ```

### Usage

Run the tool using the following command:
```bash
go run main.go
