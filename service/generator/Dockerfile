FROM golang:1.18.3-stretch

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /actuator-generator .

EXPOSE 8080

CMD [ "/actuator-generator" ]