FROM golang:latest

RUN mkdir -p /go/src/CarRental

WORKDIR /go/src/CarRental

COPY . /go/src/CarRental

RUN go install CarRental

CMD /go/bin/CarRental

EXPOSE 8080