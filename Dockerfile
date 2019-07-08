FROM golang as builder
COPY ./ /go/src/github.com/farhansajid1/go-docker
WORKDIR /go/src/github.com/farhansajid1/go-docker

RUN go get ./ # build the dependencies
RUN go build -o main

RUN chmod +x main
CMD ./main



# FROM alpine:latest
# WORKDIR /app
# COPY --from=builder /go/src/github.com/farhansajid1/go-docker/ /app/
# CMD ./main

# EXPOSE 8080