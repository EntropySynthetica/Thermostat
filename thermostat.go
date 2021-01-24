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

    // Show the data we got
    // fmt.Println(string(response_data))

    var response_stats thermo_stats
    json.Unmarshal([]byte(response_data), &response_stats)

    // Determine the Thermostat Mode
    if response_stats.Tmode == 0 {
        fmt.Println("Thermostat Mode = Off")
    } else if response_stats.Tmode == 1 {
        fmt.Println("Thermostat Mode = Heat")
    } else if response_stats.Tmode == 2 {
        fmt.Println("Thermostat Mode = Cool")
    } else if response_stats.Tmode == 3 {
        fmt.Println("Thermostat Mode = Auto")
    }

    // Show current Temp
    fmt.Println("Current Temp = " + strconv.FormatFloat(response_stats.Temp, 'f', -1, 64))

    // The target temp is returned from a diffrent var depending on if the thermostat is in heat or cool mode.  Lets check both and return which ever is not empty
    var target_temp string

    if response_stats.THeat != 0 {
        target_temp = strconv.FormatFloat(response_stats.THeat, 'f', -1, 64)
    } else if response_stats.TCool != 0 {
        target_temp = strconv.FormatFloat(response_stats.TCool, 'f', -1, 64)
    }

    fmt.Println("Target Temp = " + target_temp)

    // Show the Operational State 
    if response_stats.Tstate == 0 {
        fmt.Println("Operating Status = Off")
    } else if response_stats.Tstate == 1 {
        fmt.Println("Operating Status = Heating")
    } else if response_stats.Tstate == 2 {
        fmt.Println("Operating Status = Cooling")
    }

    // Show if an Override is active
    if response_stats.Override == 0 {
        fmt.Println("Override Off")
    } else if response_stats.Override == 1 {
        fmt.Println("Override On")
    }

    // Show if a manual hold is active
    if response_stats.Hold == 0 {
        fmt.Println("Manual Hold Off")
    } else if response_stats.Hold ==1 {
        fmt.Println("Manual Hold On")
    }

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