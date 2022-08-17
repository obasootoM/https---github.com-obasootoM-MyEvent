FROM scratch
COPY myevent /myevent

ENV LISTEN=0.0.0.0:9191
EXPOSE 9191
CMD ["/myevent"]