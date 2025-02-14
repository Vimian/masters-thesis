FROM minio/mc:RELEASE.2025-02-08T19-14-21Z

COPY ./persistence/minio/init.sh /minio/init.sh
RUN chmod +x /minio/init.sh

ENTRYPOINT ["/bin/sh", "/minio/init.sh"]