FROM alpine
RUN mkdir -p /home/firebase && chmod 777 /home/firebase
RUN mkdir -p /var/lib/docker/tmp && chmod 777 /var/lib/docker/tmp

WORKDIR /home/firebase
RUN curl -sL firebase.tools | bash

EXPOSE 8888
CMD ["./webserviceexe"]
