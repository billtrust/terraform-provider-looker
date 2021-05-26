FROM golang

WORKDIR /go/src/github.com/billtrust/terraform-provider-looker

RUN apt-get update && \
    apt-get install unzip

RUN wget https://releases.hashicorp.com/terraform/0.11.14/terraform_0.11.14_linux_amd64.zip && \
    unzip terraform_0.11.14_linux_amd64.zip && \
    chmod +x terraform && \
    mv terraform /usr/local/bin

COPY ./ .

RUN go get -v -insecure ./...
# For whatever reason this doesnt get installed with the previous call
RUN go get github.com/gruntwork-io/terratest/modules/terraform