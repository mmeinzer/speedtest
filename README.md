# Speedtest

Speedtest performs a [speed test](https://github.com/kylegrantlucas/speedtest) and sends the data via a POST request to the specified endpoint.

## How to use
`./speedtest -e ENDPOINT_URL -u USERNAME -p PASSWORD`

## To compile for the Raspberry Pi
`env GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build`
