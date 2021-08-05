FROM golang:latest

RUN apt-get update && apt-get install -y clang

RUN go get golang.org/dl/gotip && gotip download dev.fuzz

RUN GOPROXY=direct go get -u github.com/stevenjohnstone/go114-fuzz-build
RUN GOPROXY=direct go get -u github.com/magefile/mage
COPY . /fuzztests
WORKDIR /fuzztests


