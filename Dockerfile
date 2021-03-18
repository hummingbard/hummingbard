# Start with node to build frontend UI
FROM node:14.16.0

WORKDIR /home/node/app
COPY . .
WORKDIR /home/node/app/ui/js

# RUN useradd -m -u 1000 -s /usr/bin/bash node
RUN chown -R node /home/node/

USER node

RUN npm i
RUN npm build

# TODO: Use golang alpine
FROM golang:1.14

# Required for multiline echo into file
SHELL ["/bin/bash", "-c"]

# Copy over built code
WORKDIR /go/src/app
COPY --from=0 /home/node/app .

RUN useradd -m -u 1000 -s /usr/bin/bash go
RUN chown -R go /go/src/app/

# Do building in other user for security
USER go

RUN go get -d -v ./...
RUN make

RUN curl 'https://matrix-client.matrix.org:443/_matrix/client/r0/publicRooms?limit=500' > bigrooms.json

# Install Goose (required by run.sh for DB upgrades)
RUN go get -u github.com/pressly/goose/cmd/goose

# Required for the run script to generate the Goose cmdline
USER root
RUN apt-get update
RUN apt-get install -y python3-toml
# Ensure we keep running in other user
USER go

CMD ["/go/src/app/run.sh"]  

