FROM golang:1.20
COPY . .
RUN go mod download
RUN go build
ENTRYPOINT [ "./cloud-native-merge-sort" ]