FROM ubuntu

RUN apt update && apt install -y wget curl jq unzip

RUN wget https://releases.hashicorp.com/terraform/0.12.20/terraform_0.12.20_linux_amd64.zip \
 && unzip terraform_0.12.20_linux_amd64.zip \
 && mv terraform /bin/ \
 && rm -f terraform_0.12.20_linux_amd64.zip

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
 && chmod +x kubectl \
 && mv kubectl /bin/

COPY ./terraform.sh /bin/
COPY ./ /app/
