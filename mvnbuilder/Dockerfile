FROM maven:3-jdk-8

RUN apt-get install git
WORKDIR /app/

COPY ./entrypoint.sh /app/
RUN chmod +x ./entrypoint.sh
CMD ./entrypoint.sh
