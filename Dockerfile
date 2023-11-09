FROM golang:1.19
WORKDIR /workspace
COPY . /workspace

ENTRYPOINT ["go"]