# ---- build stage --------------------------------------------------------
FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" \
    -o /out/screenjson-export ./cmd/screenjson-export

# ---- runtime stage ------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/screenjson-export /usr/local/bin/screenjson-export

WORKDIR /data
USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/screenjson-export"]
CMD ["-h"]
