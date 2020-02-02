FROM scratch

ADD ciGatingServer /

EXPOSE 8080

CMD ["/webServer"]

