package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func main() {
	godotenv.Load()

	config.SetEnv(common.URLUS, os.Getenv("ACCESS_ID"), os.Getenv("ACCESS_KEY"))

	http.HandleFunc("/toggle", func(w http.ResponseWriter, r *http.Request) {
		deviceId := r.URL.Query().Get("id")
		if deviceId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id is required"))
			return
		}
		toggleLedOnOff(deviceId)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		deviceId := r.URL.Query().Get("id")
		if deviceId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id is required"))
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("code is required"))
			return
		}

		value := r.URL.Query().Get("value")
		if value == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("value is required"))
			return
		}

		var val any
		switch strings.ToLower(value) {
		case "true", "false":
			val, _ = strconv.ParseBool(value)
		default:
			tmp, err := strconv.Atoi(value)
			if err != nil {
				val = value
			} else {
				val = tmp
			}
		}

		res := setDeviceValue(deviceId, code, val)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res.Msg))
	})

	http.ListenAndServe(":8015", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}
