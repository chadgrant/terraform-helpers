FROM golang:latest

WORKDIR /go/src/github.com/chadgrant/terraform-helpers/

RUN go get -u github.com/hashicorp/terraform && \
    go get -u github.com/aws/aws-sdk-go/aws && \
    go get -u github.com/mitchellh/gox && \
    go get -u github.com/tcnksm/ghr

COPY go-build.sh .

COPY . .

RUN chmod +x go-build.sh && \
    mkdir -p /out/

CMD ["/go/src/github.com/chadgrant/terraform-helpers/go-build.sh"]
