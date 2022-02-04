FROM golang:1.17-alpine
WORKDIR /app
COPY . /app
RUN go build -mod=vendor -o /bin/app /app/
RUN rm -rf /app/*
CMD ["/bin/app"]
