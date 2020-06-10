FROM golang:alpine
WORKDIR /app
ADD .
RUN cd ./src && go build -o ./build/golang-cognito-example
EXPOSE 80
cmd ["/app/src/build/golang-cognito-example 80"]