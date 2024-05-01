FROM alpine
RUN mkdir -p /home/webservice && chmod 777 /home/webservice
RUN mkdir -p /var/lib/docker/tmp && chmod 777 /var/lib/docker/tmp

COPY webserviceexe /home/webservice/webserviceexe
RUN chmod a+rx /home/webservice/webserviceexe
WORKDIR /home/webservice

EXPOSE 8888
CMD ["./webserviceexe"]
