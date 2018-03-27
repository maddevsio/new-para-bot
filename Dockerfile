FROM golang:1.8 
RUN mkdir -p /src/github.com/maddevsio/new-para-bot
ADD . /src/github.com/maddevsio/new-para-bot
ENV GOPATH /
WORKDIR /src/github.com/maddevsio/new-para-bot
RUN go build -o new-para-bot main/main.go 
CMD ["/app/main"]