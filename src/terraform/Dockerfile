FROM ubuntu

RUN apt update && apt install -y wget jq unzip

RUN wget https://releases.hashicorp.com/terraform/0.12.20/terraform_0.12.20_linux_amd64.zip \
 && unzip terraform_0.12.20_linux_amd64.zip \
 && mv terraform /bin/ \
 && rm -f terraform_0.12.20_linux_amd64.zip

COPY ./terraform.sh /bin/
COPY ./ /app/
