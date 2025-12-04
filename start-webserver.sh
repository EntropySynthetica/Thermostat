#!/bin/bash
# Start the Thermostat Web Server

# Default port
PORT=8080

# Check if port argument was provided
if [ ! -z "$1" ]; then
    PORT=$1
fi

# Check if binary exists
if [ ! -f "./bin/webserver" ]; then
    echo "Web server binary not found. Building..."
    go build -o bin/webserver ./cmd/webserver
    if [ $? -ne 0 ]; then
        echo "Build failed!"
        exit 1
    fi
fi

# Check if config exists
CONFIG_FILE="$HOME/.config/thermostat/config.json"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Config file not found at $CONFIG_FILE"
    echo "Please run 'thermostat --new' to create a config file first."
    exit 1
fi

# Start the server
echo "Starting Thermostat Web Server on port $PORT..."
./bin/webserver -port $PORT
