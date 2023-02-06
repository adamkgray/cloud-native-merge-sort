FROM golang:1.20
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN go build
ENTRYPOINT [ "./cloud-native-merge-sort" ]