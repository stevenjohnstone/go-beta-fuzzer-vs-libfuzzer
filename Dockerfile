FROM golang:latest

RUN apt-get update && apt-get install -y clang gcc-aarch64-linux-gnu binutils-aarch64-linux-gnu

RUN go get golang.org/dl/gotip && gotip download dev.fuzz

RUN GOPROXY=direct go get -u github.com/stevenjohnstone/go114-fuzz-build
RUN GOPROXY=direct go get -u github.com/magefile/mage
COPY . /fuzztests
WORKDIR /fuzztests
RUN mage -compile /usr/local/bin/run


