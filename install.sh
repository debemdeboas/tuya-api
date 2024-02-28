#!/bin/bash

sudo cp dist/TuyaAPI.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable TuyaAPI
sudo systemctl start TuyaAPI
