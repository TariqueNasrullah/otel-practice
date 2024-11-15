FROM golang:1.23.3 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . ./

RUN go build -o /core-app


FROM ubuntu
COPY --from=build /core-app /core-app

ENTRYPOINT ["/core-app"]