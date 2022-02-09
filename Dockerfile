FROM golang:latest

WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o /app/main ./cmd/serviceapp

# Build the cli
RUN go build -o /app/cli ./cmd/serviceappcli

# Move to /dist directory as the place for resulting binary folder
# WORKDIR /dist

# Copy binary from build to main folder
#RUN cp main .

EXPOSE 8080

# Command to run when starting the container
CMD ["./main"]

#FROM golang:latest AS builder

#RUN mkdir /app
#ADD . /app
#WORKDIR /app
#RUN go get -d -v
#RUN go mod download
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /app/main .
#RUN go mod download
#RUN go build -o /main .

#FROM scratch
##ENV MONGODB_USERNAME=MONGODB_USERNAME MONGODB_PASSWORD=MONGODB_PASSWORD MONGODB_ENDPOINT=MONGODB_ENDPOINT
#COPY --from=builder /main ./
#ENTRYPOINT ["./main"]
#EXPOSE 8080

#CMD ["/app/main"]