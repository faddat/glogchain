IMAGE golang
RUN apt-get update && apt-get install -y \
  curl \
  mysql \
  wget \
RUN curl https://glide.sh/get | sh
WORKDIR $GOPATH/src/github.com/baabeetaa/glogchain
RUN glide install
RUN go build .


