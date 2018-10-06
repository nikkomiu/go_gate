FROM alpine

COPY bin/* /go_gate/
COPY lib/* /go_gate/lib/
COPY config/* /go_gate/

ENTRYPOINT [ "/go_gate/go_gate" ]
CMD [ "/go_gate/go_gate" ]
