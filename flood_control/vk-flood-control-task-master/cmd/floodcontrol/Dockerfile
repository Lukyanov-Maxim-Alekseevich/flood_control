FROM ubuntu:20.04

RUN mkdir -p /opt/app

COPY cmd/floodcontrol/build /opt/app

WORKDIR /opt/app

CMD [ "/opt/app/main", "-config.files", "container.yaml", "-env.vars.file", "application.env" ]