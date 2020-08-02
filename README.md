# Drone Navigation Service (DNS) - Atlas corp

DNS helps robots to locate databank to upload gathered data from space exploration.

- Each observed sector of the galaxy has unique numeric SectorID assigned to it
- DNS can be deployed multiple times specifying the SectorID at runtime.
- Robots can calculate the databank positions via a JSON API.

![data collectors robots](./docs/atlas-dns.png)

## Project architecture:
- The project follows recommendations from [https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- Go 1.11+ (Go modules support required)
- Gin gonic 1.6 [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- Logrus for logs [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)
- Docker is recommended for building and packaging [docker.com](https://www.docker.com/)

## Available runtime configuration (via environment variables):
- DNS_SECTOR_ID:
    Type: integer
    Default value: 1
    Description: represents the Sector of the galaxy where DNS is deployed.
- DNS_PORT:
    Type: integer
    Default value: 8080
    Description: represents the TCP port where the server will be listening for connections.
- DNS_LOG_LEVEL:
    Type: text
    Default value: INFO in production, DEBUG locally.
    Description: set the log level (see [logrus documentation](https://github.com/sirupsen/logrus)).
- DNS_ENVIRONMENT:
    Type: text
    Default value: PRODUCTION
    Description: is used to load the app configurations.

## JSON API:
#### Get nearest databank location

```bash
POST: /calculate-databank-location
```

Request body example (content-type: application/json):
```json
{
	"x": "123.12",
	"y": "456.56",
	"z": "789.89",
	"vel": "20.0"
}
```

Response example (content-type: application/json):
```json
{
    "loc": 1389.5700000000002
}
```

#### Ping health check request

```bash
GET: /ping
```

Response example (content-type: text/plain):
```json
"pong"
```

## Run Atlas DNS:

#### - Docker (recommended):

The app container can be built with the provided make script:

Build a docker image:
```bash
make build-docker
```

Run the docker image overriding config env variables:
```bash
docker run -p 8080:8081 \
    --env DNS_SECTOR_ID=2 \
    --env DNS_PORT=8081 \
    --env DNS_LOG_LEVEL=INFO \
    --env DNS_ENVIRONMENT=PRODUCTION \
    jegutierrez/atlas-dns
```

#### - Locally:

The app can be run using the provided make script, in port 8080 by default and SectorID will be 1 by default.
```bash
make run
```

## Usage:

Examples are shown using [curl](https://es.wikipedia.org/wiki/CURL) as http-client.

```bash
# Get nearest databank location
curl --request POST \
  --url http://localhost:8080/calculate-databank-location \
  --header 'content-type: application/json' \
  --data '{
	"x": "123.12",
	"y": "456.56",
	"z": "789.89",
	"vel": "20.0"
}'

# Response
HTTP/1.1 200 OK
Content-Length: 26
Content-Type: application/json; charset=utf-8
Date: Sun, 02 Aug 2020 15:04:43 GMT

{
    "loc": 1389.5700000000002
}
```

```bash
# Ping health check request
curl --request GET \
  --url http://localhost:8080/ping

#Response
HTTP/1.1 200 OK
Content-Length: 4
Content-Type: text/plain; charset=utf-8
Date: Sun, 02 Aug 2020 15:09:32 GMT

pong
```

## TODO:
- Evaluate the replacement of float64 to do complicated math operations.
