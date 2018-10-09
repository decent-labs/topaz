FROM postgres:10.5

RUN apt-get update && apt-get install -y postgresql-contrib

ENV POSTGRES_DB postgres
COPY schema.sql /docker-entrypoint-initdb.d/

ADD install.sh /docker-entrypoint-initdb.d/
RUN chmod 755 /docker-entrypoint-initdb.d/install.sh
