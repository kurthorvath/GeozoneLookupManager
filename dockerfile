FROM ubuntu:latest
EXPOSE 7000
CMD  mkdir /APP
COPY ./main /APP/
WORKDIR /APP
CMD chmod 777 main
CMD ./main