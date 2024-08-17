FROM golang:1.23 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /embed-api

FROM gcr.io/distroless/base-debian11 AS runner

WORKDIR /

COPY --from=build /embed-api /embed-api

EXPOSE 3000

ENTRYPOINT ["/embed-api"]
