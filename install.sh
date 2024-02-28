#!/bin/bash

sudo mkdir -p /www/tuya-api
go build . && sudo cp tuya-api .env /www/tuya-api/

sudo cp TuyaAPI.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable TuyaAPI
sudo systemctl start TuyaAPI
