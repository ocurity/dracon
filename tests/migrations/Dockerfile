FROM ubuntu:latest

COPY ./bin/cmd/draconctl /bin
COPY ./pkg/enrichers/migrations /etc/dracon/

ENTRYPOINT [ "/bin/bash" ]
