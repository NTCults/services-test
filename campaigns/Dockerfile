FROM golang:1.9

ADD . /go/src/myapp

WORKDIR /go/src/myapp
        
RUN go get myapp
       
CMD ["go","run","main.go"]