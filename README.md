#CT50 Thermostat App

This is a simple CLI app to read and control the Radio Thermostat CT50 Thermostat from the CLI.  

## Installation

### Linux
todo

### Mac
todo

### Windows
todo

## Usage

### First run Setup
The thermostat app uses a config file at ~/.config/thermostat/config.json

Run the thermostat command with --new to create a base config file.  
```
thermostat --new
```

You can now edit the newly created config file and enter the IP of your thermostat.  

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


## Manual build
```
go build -o bin/thermostat
```