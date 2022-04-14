FROM golang
ENV MYPASS="pass", DBHOST="localhost"
WORKDIR /go/src/app
COPY . .
RUN go get github.com/go-sql-driver/mysql
RUN go build -o main .
CMD ["go","run","main.go"]
