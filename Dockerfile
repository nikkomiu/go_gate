FROM alpine

COPY bin/* /usr/bin
COPY lib/* /etc/go_gate/lib
COPY config/* /etc/go_gate

ENTRYPOINT [ "/usr/bin/go_gate" ]
CMD [ "/usr/bin/go_gate", "--config", "/etc/go_gate/app.yaml" ]
