package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const WebServerVersion = "1.0.0"

// ThermoStats represents the thermostat status
type ThermoStats struct {
	Temp     float64 `json:"temp"`
	Tmode    int     `json:"tmode"`
	Fmode    int     `json:"fmode"`
	Override int     `json:"override"`
	Hold     int     `json:"hold"`
	THeat    float64 `json:"t_heat"`
	TCool    float64 `json:"t_cool"`
	Tstate   int     `json:"tstate"`
	Fstate   int     `json:"fstate"`
	Time     struct {
		Day    int `json:"day"`
		Hour   int `json:"hour"`
		Minute int `json:"minute"`
	} `json:"time"`
	TTypePost int `json:"t_type_post"`
}

// Config represents the application configuration
type Config struct {
	ThermostatIP string `json:"ThermostatIP"`
}

// StatusResponse represents the formatted status for the web UI
type StatusResponse struct {
	CurrentTemp    float64 `json:"currentTemp"`
	TargetTemp     float64 `json:"targetTemp"`
	Mode           string  `json:"mode"`
	ModeCode       int     `json:"modeCode"`
	OperatingState string  `json:"operatingState"`
	Override       string  `json:"override"`
	Hold           string  `json:"hold"`
}

var thermostatIP string

// getStats retrieves the current thermostat status
func getStats(ip string) (*ThermoStats, error) {
	pollURL := "http://" + ip + "/tstat"

	response, err := http.Get(pollURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var stats ThermoStats
	err = json.Unmarshal(responseData, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// formatStats converts raw stats to a user-friendly format
func formatStats(stats *ThermoStats) *StatusResponse {
	status := &StatusResponse{
		CurrentTemp: stats.Temp,
		ModeCode:    stats.Tmode,
	}

	// Determine the Thermostat Mode
	switch stats.Tmode {
	case 0:
		status.Mode = "Off"
	case 1:
		status.Mode = "Heat"
	case 2:
		status.Mode = "Cool"
	case 3:
		status.Mode = "Auto"
	default:
		status.Mode = "Unknown"
	}

	// Get target temp based on mode
	if stats.THeat != 0 {
		status.TargetTemp = stats.THeat
	} else if stats.TCool != 0 {
		status.TargetTemp = stats.TCool
	}

	// Show the Operational State
	switch stats.Tstate {
	case 0:
		status.OperatingState = "Off"
	case 1:
		status.OperatingState = "Heating"
	case 2:
		status.OperatingState = "Cooling"
	default:
		status.OperatingState = "Unknown"
	}

	// Show if an Override is active
	switch stats.Override {
	case 0:
		status.Override = "Off"
	case 1:
		status.Override = "On"
	default:
		status.Override = "Unknown"
	}

	// Show if a manual hold is active
	switch stats.Hold {
	case 0:
		status.Hold = "Off"
	case 1:
		status.Hold = "On"
	default:
		status.Hold = "Unknown"
	}

	return status
}

// setTemp sets the target temperature
func setTemp(ip string, temp int) error {
	pollURL := "http://" + ip + "/tstat"

	// Get current mode
	stats, err := getStats(ip)
	if err != nil {
		return err
	}

	// Craft payload based on current mode
	var query string
	if stats.Tmode == 1 {
		query = fmt.Sprintf("{\"tmode\":1,\"t_heat\":%d}", temp)
	} else if stats.Tmode == 2 {
		query = fmt.Sprintf("{\"tmode\":2,\"t_cool\":%d}", temp)
	} else {
		return fmt.Errorf("thermostat must be in heat or cool mode to set temperature")
	}

	payload := strings.NewReader(query)
	req, err := http.NewRequest("POST", pollURL, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	postResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer postResponse.Body.Close()

	return nil
}

// setMode sets the thermostat operating mode
func setMode(ip string, mode int) error {
	pollURL := "http://" + ip + "/tstat"

	query := fmt.Sprintf("{\"tmode\":%d}", mode)
	payload := strings.NewReader(query)

	req, err := http.NewRequest("POST", pollURL, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	postResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer postResponse.Body.Close()

	return nil
}

// API Handlers

func handleStatus(w http.ResponseWriter, r *http.Request) {
	stats, err := getStats(thermostatIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := formatStats(stats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleSetTemp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Temp int `json:"temp"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Temp < 50 || req.Temp > 90 {
		http.Error(w, "Temperature must be between 50 and 90", http.StatusBadRequest)
		return
	}

	err = setTemp(thermostatIP, req.Temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleSetMode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Mode int `json:"mode"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Mode < 0 || req.Mode > 3 {
		http.Error(w, "Mode must be 0 (Off), 1 (Heat), 2 (Cool), or 3 (Auto)", http.StatusBadRequest)
		return
	}

	err = setMode(thermostatIP, req.Mode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(indexHTML))
	tmpl.Execute(w, nil)
}

const indexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Thermostat Control</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            padding: 40px;
            max-width: 500px;
            width: 100%;
        }
        h1 {
            text-align: center;
            color: #333;
            margin-bottom: 30px;
            font-size: 2em;
        }
        .status-card {
            background: #f8f9fa;
            border-radius: 15px;
            padding: 25px;
            margin-bottom: 30px;
        }
        .temp-display {
            text-align: center;
            font-size: 3.5em;
            font-weight: bold;
            color: #667eea;
            margin: 20px 0;
        }
        .status-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 15px;
            margin-top: 20px;
        }
        .status-item {
            background: white;
            padding: 15px;
            border-radius: 10px;
            text-align: center;
        }
        .status-label {
            font-size: 0.9em;
            color: #666;
            margin-bottom: 5px;
        }
        .status-value {
            font-weight: bold;
            color: #333;
            font-size: 1.1em;
        }
        .control-section {
            margin-bottom: 30px;
        }
        .control-title {
            font-size: 1.2em;
            color: #333;
            margin-bottom: 15px;
            font-weight: 600;
        }
        .temp-control {
            display: flex;
            align-items: center;
            justify-content: space-between;
            background: #f8f9fa;
            padding: 15px;
            border-radius: 10px;
        }
        .temp-input {
            font-size: 2em;
            width: 100px;
            text-align: center;
            border: 2px solid #ddd;
            border-radius: 8px;
            padding: 10px;
        }
        .temp-button {
            background: #667eea;
            color: white;
            border: none;
            border-radius: 50%;
            width: 50px;
            height: 50px;
            font-size: 1.5em;
            cursor: pointer;
            transition: all 0.3s;
        }
        .temp-button:hover {
            background: #5568d3;
            transform: scale(1.1);
        }
        .temp-button:active {
            transform: scale(0.95);
        }
        .mode-buttons {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 10px;
        }
        .mode-button {
            padding: 15px;
            border: 2px solid #ddd;
            background: white;
            border-radius: 10px;
            cursor: pointer;
            font-size: 1em;
            font-weight: 600;
            transition: all 0.3s;
        }
        .mode-button:hover {
            border-color: #667eea;
            transform: translateY(-2px);
        }
        .mode-button.active {
            background: #667eea;
            color: white;
            border-color: #667eea;
        }
        .set-temp-button {
            width: 100%;
            padding: 15px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 10px;
            font-size: 1.1em;
            font-weight: 600;
            cursor: pointer;
            margin-top: 15px;
            transition: all 0.3s;
        }
        .set-temp-button:hover {
            background: #5568d3;
            transform: translateY(-2px);
        }
        .set-temp-button:active {
            transform: translateY(0);
        }
        .message {
            padding: 15px;
            border-radius: 10px;
            margin-top: 15px;
            text-align: center;
            font-weight: 600;
            display: none;
        }
        .message.success {
            background: #d4edda;
            color: #155724;
        }
        .message.error {
            background: #f8d7da;
            color: #721c24;
        }
        .refresh-button {
            position: absolute;
            top: 20px;
            right: 20px;
            background: rgba(255, 255, 255, 0.9);
            border: none;
            border-radius: 50%;
            width: 40px;
            height: 40px;
            cursor: pointer;
            font-size: 1.2em;
            transition: all 0.3s;
        }
        .refresh-button:hover {
            transform: rotate(180deg);
        }
    </style>
</head>
<body>
    <div class="container" style="position: relative;">
        <button class="refresh-button" onclick="loadStatus()">🔄</button>
        <h1>🌡️ Thermostat Control</h1>
        
        <div class="status-card">
            <div class="status-label">Current Temperature</div>
            <div class="temp-display" id="currentTemp">--°F</div>
            
            <div class="status-grid">
                <div class="status-item">
                    <div class="status-label">Target</div>
                    <div class="status-value" id="targetTemp">--°F</div>
                </div>
                <div class="status-item">
                    <div class="status-label">Mode</div>
                    <div class="status-value" id="mode">--</div>
                </div>
                <div class="status-item">
                    <div class="status-label">Status</div>
                    <div class="status-value" id="operatingState">--</div>
                </div>
                <div class="status-item">
                    <div class="status-label">Hold</div>
                    <div class="status-value" id="hold">--</div>
                </div>
            </div>
        </div>

        <div class="control-section">
            <div class="control-title">Set Temperature</div>
            <div class="temp-control">
                <button class="temp-button" onclick="adjustTemp(-1)">−</button>
                <input type="number" id="tempInput" class="temp-input" value="70" min="50" max="90">
                <button class="temp-button" onclick="adjustTemp(1)">+</button>
            </div>
            <button class="set-temp-button" onclick="setTemperature()">Set Temperature</button>
        </div>

        <div class="control-section">
            <div class="control-title">Operating Mode</div>
            <div class="mode-buttons">
                <button class="mode-button" data-mode="0" onclick="setMode(0)">Off</button>
                <button class="mode-button" data-mode="1" onclick="setMode(1)">Heat</button>
                <button class="mode-button" data-mode="2" onclick="setMode(2)">Cool</button>
                <button class="mode-button" data-mode="3" onclick="setMode(3)">Auto</button>
            </div>
        </div>

        <div class="message" id="message"></div>
    </div>

    <script>
        let currentMode = 0;

        function showMessage(text, type) {
            const msg = document.getElementById('message');
            msg.textContent = text;
            msg.className = 'message ' + type;
            msg.style.display = 'block';
            setTimeout(() => {
                msg.style.display = 'none';
            }, 3000);
        }

        function adjustTemp(delta) {
            const input = document.getElementById('tempInput');
            let value = parseInt(input.value) || 70;
            value += delta;
            value = Math.max(50, Math.min(90, value));
            input.value = value;
        }

        async function loadStatus() {
            try {
                const response = await fetch('/api/status');
                if (!response.ok) throw new Error('Failed to load status');
                
                const data = await response.json();
                
                document.getElementById('currentTemp').textContent = data.currentTemp.toFixed(1) + '°F';
                document.getElementById('targetTemp').textContent = data.targetTemp.toFixed(1) + '°F';
                document.getElementById('mode').textContent = data.mode;
                document.getElementById('operatingState').textContent = data.operatingState;
                document.getElementById('hold').textContent = data.hold;
                
                currentMode = data.modeCode;
                updateModeButtons();
                
                if (data.targetTemp > 0) {
                    document.getElementById('tempInput').value = Math.round(data.targetTemp);
                }
            } catch (error) {
                showMessage('Failed to load status: ' + error.message, 'error');
            }
        }

        function updateModeButtons() {
            document.querySelectorAll('.mode-button').forEach(btn => {
                const mode = parseInt(btn.getAttribute('data-mode'));
                if (mode === currentMode) {
                    btn.classList.add('active');
                } else {
                    btn.classList.remove('active');
                }
            });
        }

        async function setTemperature() {
            const temp = parseInt(document.getElementById('tempInput').value);
            
            if (temp < 50 || temp > 90) {
                showMessage('Temperature must be between 50 and 90', 'error');
                return;
            }

            try {
                const response = await fetch('/api/settemp', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ temp: temp })
                });

                if (!response.ok) {
                    const error = await response.text();
                    throw new Error(error);
                }

                showMessage('Temperature set to ' + temp + '°F', 'success');
                setTimeout(loadStatus, 1000);
            } catch (error) {
                showMessage('Failed to set temperature: ' + error.message, 'error');
            }
        }

        async function setMode(mode) {
            try {
                const response = await fetch('/api/setmode', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ mode: mode })
                });

                if (!response.ok) {
                    const error = await response.text();
                    throw new Error(error);
                }

                const modeNames = ['Off', 'Heat', 'Cool', 'Auto'];
                showMessage('Mode set to ' + modeNames[mode], 'success');
                setTimeout(loadStatus, 1000);
            } catch (error) {
                showMessage('Failed to set mode: ' + error.message, 'error');
            }
        }

        // Load initial status
        loadStatus();
        
        // Auto-refresh every 30 seconds
        setInterval(loadStatus, 30000);
    </script>
</body>
</html>
`

func loadConfig(configFile string) (*Config, error) {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	homedir, _ := os.UserHomeDir()
	var configFile string
	var port string
	var thermostatIPFlag string

	// Parse CLI Flags
	flag.StringVar(&configFile, "c", homedir+"/.config/thermostat/config.json", "specify path of config file")
	flag.StringVar(&port, "port", "8080", "port to run the web server on")
	flag.StringVar(&thermostatIPFlag, "ip", "", "thermostat IP address (overrides config file)")
	showVer := flag.Bool("v", false, "Show Version")
	flag.Parse()

	// Print Version of app
	if *showVer {
		fmt.Println("CT50 Thermostat Web Server Version: " + WebServerVersion)
		return
	}

	// Priority: 1. Environment variable, 2. Command line flag, 3. Config file
	thermostatIP = os.Getenv("THERMOSTAT_IP")

	if thermostatIP == "" && thermostatIPFlag != "" {
		thermostatIP = thermostatIPFlag
	}

	if thermostatIP == "" {
		// Load configuration from file
		config, err := loadConfig(configFile)
		if err != nil {
			log.Fatalf("Error loading config file: %v\nPlease set THERMOSTAT_IP environment variable, use -ip flag, or run 'thermostat --new' to create a config file", err)
		}
		thermostatIP = config.ThermostatIP
	}

	if thermostatIP == "" {
		log.Fatal("Thermostat IP not configured. Set THERMOSTAT_IP environment variable, use -ip flag, or configure in config file")
	}

	// Set up HTTP routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/status", handleStatus)
	http.HandleFunc("/api/settemp", handleSetTemp)
	http.HandleFunc("/api/setmode", handleSetMode)

	// Start server
	addr := ":" + port
	fmt.Printf("Starting Thermostat Web Server v%s\n", WebServerVersion)
	fmt.Printf("Thermostat IP: %s\n", thermostatIP)
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	fmt.Println("Press Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(addr, nil))
}
