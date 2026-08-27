FROM golang:1.26
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=0
WORKDIR /workspace
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .
RUN go build -mod=vendor -o /ammonialoop ./cmd/ammonialoop
CMD ["/ammonialoop"]
