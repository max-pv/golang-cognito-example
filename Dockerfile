FROM golang:alpine
WORKDIR /app
ADD . .
RUN go build -o ./build/cognito
EXPOSE 80
CMD ["PORT=80 /app/src/build/cognito"]