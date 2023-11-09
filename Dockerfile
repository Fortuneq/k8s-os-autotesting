FROM dzo.sw.sbc.space/sbt/ci90000051_synai/golang:1.19
WORKDIR /workspace
COPY . /workspace

ENTRYPOINT ["go"]