FROM scratch

ADD petServer /

EXPOSE 8080

CMD ["/petServer"]

