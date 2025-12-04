#CT50 Thermostat App

This repo contains two applications for controlling the Radio Thermostat CT50 Thermostat:
1. **thermostat** - A CLI application for command-line control
2. **webserver** - A web-based GUI for browser-based control

## Project Structure
```
.
├── thermostat.go          # CLI application source
├── cmd/
│   └── webserver/
│       └── main.go        # Web server application source
├── start-webserver.sh     # Convenience script to start web server
├── Makefile              # Build automation
├── Dockerfile            # Docker container definition
├── docker-compose.yml    # Docker Compose configuration
├── .dockerignore         # Docker build exclusions
├── bin/
│   ├── thermostat        # Compiled CLI binary
│   └── webserver         # Compiled web server binary
└── README.md
```  

## Installation

### Linux
todo

### Mac
todo

### Windows
todo

## Usage

### First run Setup
Both applications use a config file at ~/.config/thermostat/config.json

Run the thermostat CLI command with --new to create a base config file.  
```
thermostat --new
```

You can now edit the newly created config file and enter the IP of your thermostat.

Alternatively, you can manually create the config file using the example:
```bash
mkdir -p ~/.config/thermostat
cp config.example.json ~/.config/thermostat/config.json
# Edit the file and change the IP to match your thermostat
```

---

## CLI Application (thermostat)

The CLI application provides command-line access to thermostat functions.  

### Get CLI options
```
thermostat --help
```

### Check Current Thermostat settings
```
thermostat
```

### Manually set the temp to 70f
```
thermostat --temp 70
```

### Set mode to Heating
```
thermostat --mode heat
```

### Set mode to Cooling
```
thermostat --mode cool
```

---

## Web Server Application (webserver)

The web server provides a modern, responsive web interface for controlling your thermostat from any browser.

### Starting the Web Server
```bash
# Build the web server
go build -o bin/webserver ./cmd/webserver

# Run the web server (default port 8080)
./bin/webserver

# Run on a custom port
./bin/webserver -port 3000

# Or use the convenience script
./start-webserver.sh

# Start on a different port using the script
./start-webserver.sh 3000
```

### Accessing the Web Interface
Once started, open your browser and navigate to:
```
http://localhost:8080
```

### Features
- **Real-time Status Display**: View current temperature, target temperature, operating mode, and system status
- **Temperature Control**: Adjust target temperature with +/- buttons or direct input
- **Mode Switching**: Easily switch between Off, Heat, Cool, and Auto modes
- **Auto-refresh**: Status updates automatically every 30 seconds
- **Responsive Design**: Works on desktop, tablet, and mobile devices
- **Visual Feedback**: Color-coded status and smooth animations

### Security Note
The web server is designed for use on a local network. If you plan to expose it to the internet, consider adding authentication and using HTTPS.

### Web Server Options
```bash
# Show version
./bin/webserver -v

# Specify custom config file
./bin/webserver -c /path/to/config.json

# Change server port
./bin/webserver -port 3000

# Set thermostat IP via command line
./bin/webserver -ip 192.168.1.100

# Set thermostat IP via environment variable
THERMOSTAT_IP=192.168.1.100 ./bin/webserver
```

**Configuration Priority:**
1. `THERMOSTAT_IP` environment variable (highest priority)
2. `-ip` command line flag
3. Config file at `~/.config/thermostat/config.json` (lowest priority)

### Docker Deployment

#### Using Docker Compose (Recommended)
```bash
# Create environment file from example
cp .env.example .env

# Edit .env and set your thermostat IP
nano .env

# Build and start the container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container
docker-compose down
```

#### Using Docker directly
```bash
# Build the image
docker build -t thermostat-web .

# Run the container
docker run -d \
  -p 8080:8080 \
  -e THERMOSTAT_IP=192.168.1.100 \
  --name thermostat-web \
  thermostat-web

# View logs
docker logs -f thermostat-web

# Stop the container
docker stop thermostat-web
docker rm thermostat-web
```
---

## Manual build

### Using Make (Recommended)
```bash
# Build both applications
make

# Build only CLI
make cli

# Build only web server
make webserver

# Build and run web server
make run-web

# Clean build artifacts
make clean

# Show all available commands
make help
```

### Using Go directly
```bash
# Build CLI application
go build -o bin/thermostat thermostat.go

# Build web server
go build -o bin/webserver ./cmd/webserver

# Build both
go build -o bin/thermostat thermostat.go && go build -o bin/webserver ./cmd/webserver
```