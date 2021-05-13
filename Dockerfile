FROM golang:latest

RUN apt-get -y update && apt-get -y upgrade && \
    apt-get install -y git \
    make openssh-client

RUN mkdir /app
COPY . /app
WORKDIR /app

CMD ["go", "run", "golive"]
