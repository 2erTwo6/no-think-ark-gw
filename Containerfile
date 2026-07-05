FROM docker.io/golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN go build -o /gateway .

FROM docker.io/alpine:3.21
COPY --from=build /gateway /gateway
EXPOSE 8080
ENTRYPOINT ["/gateway"]
