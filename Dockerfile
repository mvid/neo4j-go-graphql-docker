FROM golang:1.12-alpine3.9 AS neo-seabolt-go

# dependencies
RUN apk add --update --no-cache ca-certificates cmake make g++ openssl-dev git curl pkgconfig

# build seabolt
RUN git clone -b 1.7 https://github.com/neo4j-drivers/seabolt.git /seabolt
WORKDIR /seabolt/build
RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install

WORKDIR /src/app
ADD go.mod go.mod
RUN go mod download

ADD . /src/app

RUN CGO_ENABLED=1 go build -o main .
CMD "/src/app/main"