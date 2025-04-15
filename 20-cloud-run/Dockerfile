FROM golang:1.21-slim AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cloudrun

FROM scratch
WORKDIR /app
COPY --from=build /app/cloudrun .
ENTRYPOINT [ "./cloudrun" ]
