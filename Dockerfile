FROM golang:1.23 

WORKDIR /app 

COPY . . 

EXPOSE 8080

CMD ["go", "run", "/app/cmd/web"]