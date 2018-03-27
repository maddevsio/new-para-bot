FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o new-para-bot main/main.go 
CMD ["/app/main"]