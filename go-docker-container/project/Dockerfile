# Start from the latest golang base image
FROM golang:alpine AS builder

# Add Maintainer Info
LABEL maintainer="Joseph Keller <jbkeller@gmail.com>"
#Used to make a bash profile we can enter this container with
RUN apk add --no-cache bash
#MakeDirectory
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./...

#Make the Scratch Production Image
FROM alpine:latest AS production
#Copy the contents of the builder stage into this app directory
COPY --from=builder /app .
#Expose port 8080:8080
EXPOSE 8080
#Create an environment variable for application to work with
ENV TEST_VAR_ENV42069="I'm waking up to ash and dust"
#Environment variables for Mongo
ENV MONGO_URI=""
#Run the executable,(which is whatever name you gave it in the previous step)
CMD ["./main"]