FROM registry.z.xinu.tv/golang/alpine/onbuild AS base
FROM alpine

COPY --from=base /go/bin/app /bin/fracserv
COPY --from=base /go/src/app/static/ /data/static/
COPY --from=base /go/src/app/templates/ /data/templates/
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD [ "/go/bin/fracserv", "-staticDir", "/go/src/app/static", "-templateDir", "/go/src/app/templates" ]
