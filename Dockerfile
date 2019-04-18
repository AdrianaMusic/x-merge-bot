# build stage
FROM golang:alpine AS build-env
COPY . /go/src/xMergeBot/
RUN cd /go/src/xMergeBot/ && go get && go build -o startMergeBot

# final stage
FROM alpine
RUN apk add ca-certificates
WORKDIR /
COPY --from=build-env /go/src/xMergeBot/startMergeBot /
CMD ./startMergeBot