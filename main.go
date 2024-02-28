// Package tuya-api is an http server for interacting with the Tuya Cloud API.
//
// It provides two endpoints:
//
//  1. /toggle?id=DEVICE_ID
//     Toggles the LED on the device with the given DEVICE_ID.
//
//  2. /set?id=DEVICE_ID&code=CODE&value=VALUE
//     Sets the value of the device with the given DEVICE_ID for the given CODE to the given VALUE.
//
// The server listens on port 8015.
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"github.com/tuya/tuya-cloud-sdk-go/api/common"
	"github.com/tuya/tuya-cloud-sdk-go/api/device"
	"github.com/tuya/tuya-cloud-sdk-go/config"
)

func toggleLedOnOff(id string) {
	dev, _ := device.GetDeviceStatus(id)

	for _, v := range dev.Result {
		if v.Code == "switch_led" {
			cmds := []device.Command{{Code: "switch_led", Value: !v.Value.(bool)}}
			device.PostDeviceCommand(id, cmds)
			break
		}
	}
}

func setDeviceValue(id string, code string, value any) *device.PostDeviceCommandResponse {
	cmds := []device.Command{{Code: code, Value: value}}
	res, err := device.PostDeviceCommand(id, cmds)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return res
}

func getRequiredValue(name string, w http.ResponseWriter, r *http.Request) (val string, err error) {
	val = r.URL.Query().Get(name)
	if val == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("'" + name + "' is required\n"))
		err = fmt.Errorf("'" + name + "' is required")
	}
	return
}

func getRequiredValues(w http.ResponseWriter, r *http.Request, names ...string) (vals []string, err error) {
	vals = make([]string, len(names))
	for i, name := range names {
		vals[i], err = getRequiredValue(name, w, r)
		if err != nil {
			return
		}
	}
	return
}

func someEffortConvert(value string) any {
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	} else if i, err := strconv.Atoi(value); err == nil {
		return i
	} else if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	} else {
		return value
	}
}

func main() {
	godotenv.Load()

	config.SetEnv(common.URLUS, os.Getenv("ACCESS_ID"), os.Getenv("ACCESS_KEY"))

	http.HandleFunc("/toggle", func(w http.ResponseWriter, r *http.Request) {
		deviceId, err := getRequiredValue("id", w, r)
		if err != nil {
			return
		}
		toggleLedOnOff(deviceId)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		vals, err := getRequiredValues(w, r, "id", "code", "value")
		if err != nil {
			return
		}

		deviceId, code, value := vals[0], vals[1], vals[2]

		res := setDeviceValue(deviceId, code, someEffortConvert(value))
		if res.Success {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(res.Msg))
		}
	})

	http.ListenAndServe(":8015", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}
