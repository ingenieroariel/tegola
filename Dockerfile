FROM golang

RUN cd cmd/tegola; go build -o tegola *.go

COPY tegola /

CMD ["/tegola"]
EXPOSE 8080


## In your Dockerfile you would have: 
#
# FROM terranodo/tegola
# COPY config.toml /
# CMD ["/tegola", "--config=/config.toml"]