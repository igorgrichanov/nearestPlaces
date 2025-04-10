FROM golang AS build
WORKDIR /go/src/nearestPlaces
COPY go.mod go.sum ./
RUN go mod download

COPY config config/
COPY cmd cmd/
COPY datasets datasets/
COPY internal internal/
COPY templates templates/

RUN go build -o /bin/main ./cmd/server/server.go
CMD ["/bin/main"]