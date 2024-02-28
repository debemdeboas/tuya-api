# [Tuya API](https://github.com/debemdeboas/tuya-api.git)

Control your Tuya devices using this API and a Tuya account.

## Motivation

I have a few smart home devices and I wanted to control them using the iOS Shortcuts app.
Namely, I have a collection of "Positivo" and "Gaya Smart" devices whose apps don't expose an API to the Shortcuts system.

First I tried to use the Alexa API to control them, but the was getting too complex and clunky for this use-case, so I switched gears.
After some research I found that these devices are actually rebranded Tuya devices and that Tuya has open-source SDKs for controlling them.
With the Tuya SDK and a Tuya Developer account in hand, I set out to create this simple API to control my home.

**tl;dr**: this project is a Go HTTP server that uses the Tuya Go SDK to control Tuya devices.

## Environment set up

1. Follow [this guide](https://github.com/codetheweb/tuyapi/blob/master/docs/SETUP.md#linking-a-tuya-device-with-smart-link) to create a Tuya Developer account and link your devices to it.
2. Copy the `ACCESS_ID` and `ACCESS_KEY` values from the Tuya Developer account to a `.env` file.
You'll also need the IDs of the devices that you wish to control.

Example `.env` file:

```sh
ACCESS_ID=your_access_id
ACCESS_KEY=your_access_key
```

## Docker

You can run this API using Docker.

> [!TIP]
> This example uses a `.env` file in the current directory to set the environment variables.
>
> Alternatively, you could set the environment variables directly in the `docker run` command via the `-e` option.

```sh
$ docker run -it --name tuya-api --rm -p 8015:8015 -v `pwd`/.env:/app/.env debemdeboas/tuya-api:latest
```
