FROM scratch
ADD timezone-webhook /
EXPOSE 8080
CMD ["./timezone-webhook"]

