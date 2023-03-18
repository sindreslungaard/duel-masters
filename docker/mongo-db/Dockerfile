FROM mongo:3.6

RUN apt-get update \
  && apt-get install -y \
    vim bash-completion

COPY ./init /docker-entrypoint-initdb.d/
