FROM postgres:12-alpine

ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=Iseestars

COPY ../providers/fixtures/pg_hba.conf postgresql.conf /etc/postgresql/config/
COPY ../providers/fixtures/root.crt server.crt server.key /etc/postgres/security/

EXPOSE 8014

CMD ["postgres", "-c", "config_file=/etc/postgresql/config/postgresql.conf", "-c", "hba_file=/etc/postgresql/config/pg_hba.conf"]