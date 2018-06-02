FROM alpine:3.7

LABEL MAINTAINER "Aurelien PERRIER <a.perrier89@gmail.com>"
LABEL APP "tfstate"

ENV TERRAFORM_PATH /root/.tfversion/bin
ENV PATH "$PATH:${TERRAFORM_PATH}"

# Copy binary
COPY bin/tfstate /usr/bin

EXPOSE 8000

ENTRYPOINT [ "tfstate", "serve" ]