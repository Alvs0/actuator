FROM golang:1.18.2-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /actuator-processor

EXPOSE 8082

CMD [ "/actuator-processor" ]