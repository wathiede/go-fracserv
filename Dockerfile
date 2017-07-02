FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download ./cmd/fracserv  # "go get -d -v ./..."
RUN go-wrapper install ./cmd/fracserv   # "go install -v ./..."

CMD [ "/go/bin/fracserv", "-staticDir", "/go/src/app/static", "-templateDir", "/go/src/app/templates" ]
