FROM golang:1.20
WORKDIR /app
COPY . .
RUN go mod download
RUN go build
ENTRYPOINT ["/app/cloud-native-merge-sort"]