FROM ubuntu:bionic

COPY bookingservice /bookingservice
RUN useradd bookingservice
USER bookingservice

ENV LISTEN=0.0.0.0:9191
EXPOSE 9191
CMD ["/bookingservice"]