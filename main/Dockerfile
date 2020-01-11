FROM golang:1.13-alpine3.11 AS neo-seabolt-go

# dependencies
RUN apk add --update --no-cache ca-certificates cmake make g++ libressl-dev openssl-dev git curl pkgconfig

# build seabolt
RUN git clone -b 1.7 https://github.com/neo4j-drivers/seabolt.git /seabolt
WORKDIR /seabolt/build
RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install

#todo: split this image here so it can be used as a base for a multi image build

WORKDIR /src/app
ADD go.mod go.mod
RUN go mod download

ADD . /src/app

RUN CGO_ENABLED=1 go build -o main .
CMD "/src/app/main"