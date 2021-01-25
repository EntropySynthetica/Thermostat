package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "log"
    "os"
    "encoding/json"
    "strconv"
    "github.com/joho/godotenv"
)

type thermo_stats struct {
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

func get_stats(ip string) {
    poll_url := "http://" + ip + "/tstat"

    // Poll the API
    response, err := http.Get(poll_url)

    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    response_data, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    // Parse the JSON
    var response_stats thermo_stats
    json.Unmarshal([]byte(response_data), &response_stats)

    // Determine the Thermostat Mode
    var Tmode_text string
    switch response_stats.Tmode {
        case 0:
            Tmode_text = "Off"
        case 1:
            Tmode_text = "Heat"
        case 2:
            Tmode_text = "Cool"
        case 3:
            Tmode_text = "Auto"
        default:
            Tmode_text = "Unknown"
    }
    fmt.Println("Thermostat Mode = " + Tmode_text)

    // Show current Temp
    fmt.Println("Current Temp = " + strconv.FormatFloat(response_stats.Temp, 'f', -1, 64))

    // The target temp is returned from a diffrent var depending on if the thermostat is in heat or cool mode.  Lets check both and return which ever is not empty
    var target_temp string
    switch {
        case response_stats.THeat != 0:
            target_temp = strconv.FormatFloat(response_stats.THeat, 'f', -1, 64)
        case response_stats.TCool != 0:
            target_temp = strconv.FormatFloat(response_stats.TCool, 'f', -1, 64)
        default:
            target_temp = "Unknown"
    }
    fmt.Println("Target Temp = " + target_temp)

    // Show the Operational State 
    var Tstate_text string
    switch response_stats.Tstate {
        case 0:
            Tstate_text = "Off"
        case 1: 
            Tstate_text = "Heating"
        case 2:
            Tstate_text = "Cooling"
        default: 
            Tstate_text = "Unknown"
    }
    fmt.Println("Operating Status = " + Tstate_text)

    // Show if an Override is active
    var Override_text string
    switch response_stats.Override {
        case 0:
            Override_text = "Off"
        case 1: 
            Override_text = "On"
        default:
            Override_text = "Unknown"
    }
    fmt.Println("Override " + Override_text)

    // Show if a manual hold is active
    var Hold_text string
    switch response_stats.Hold {
        case 0:
            Hold_text = "Off"
        case 1: 
            Hold_text = "On"
        default:
            Hold_text = "Unknown"
    }
    fmt.Println("Manual Hold " + Hold_text)

}

func main() {

    // Get vars from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    thermostat_ip := os.Getenv("THERMOSTAT_IP")

    // Get Thermostat Stats
    get_stats(thermostat_ip)
}