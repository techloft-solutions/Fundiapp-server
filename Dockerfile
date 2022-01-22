FROM golang:latest

WORKDIR /app

# Copy and download dependency using go mod
#COPY go.mod .
#COPY go.sum .
#RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main ./cmd/hudumaapp

# Move to /dist directory as the place for resulting binary folder
# WORKDIR /dist

# Copy binary from build to main folder
#RUN cp main .

EXPOSE 8080

# Command to run when starting the container
CMD [ "./main" ]