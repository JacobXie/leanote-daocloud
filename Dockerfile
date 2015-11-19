FROM ubuntu:14.04
MAINTAINER JacobXie "xieyuehong2010@gmail.com"

# add golang repository
RUN apt-get update
RUN apt-get install -y python-software-properties software-properties-common

RUN apt-get -y install golang
RUN export GOPATH=/gopath
RUN export PATH=$PATH:$GOPATH/bin

# install git and openssh
RUN apt-get install -y git-core mercurial openssh-server openssh-client

ADD src/ $GOPATH/src/

RUN go install github.com/revel/cmd/revel

EXPOSE 9000

CMD revel run github.com/Jacobxie/leanote-daocloud daocloud
