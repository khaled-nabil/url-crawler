# Use the official Golang image as the base image
FROM golang:latest

WORKDIR /app

COPY . .

#RUN mkdir -p /crawlie
#RUN chown crawluser /crawlie
#USER crawluser

USER root
RUN mkdir /crawlie
#RUN chmod 755 /crawlie
#RUN chown -R crawluser:crawluser /crawlie
#USER crawluser

RUN go build -o crawler

EXPOSE 8080

CMD ["./crawler"]
