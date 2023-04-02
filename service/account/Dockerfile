FROM golang:1.18.3-stretch

WORKDIR /app

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /actuator-account .

EXPOSE 8082

CMD [ "/actuator-account" ]