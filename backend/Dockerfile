FROM golang:1.26.4-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/gobooker .

FROM alpine:3.21

WORKDIR /app
COPY --from=build /out/gobooker ./gobooker

EXPOSE 8081

CMD ["./gobooker"]