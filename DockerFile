FROM golang:1.5.1-onbuild

MAINTAINER JacobXie "xieyuehong2010@gmail.com"

RUN apt-get update

RUN export GOPATH=/gopath
RUN export PATH=$PATH:$GOPATH/bin

ADD src/ $GOPATH/src/

RUN go install github.com/revel/cmd/revel

EXPOSE 9000

CMD revel run github.com/Jacobxie/leanote-daocloud daocloud
