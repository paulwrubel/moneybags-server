# lightweight base image
FROM alpine:3.12

# add timezone data
# this allows the user to specify a timezone as an environment variable: TZ
RUN apk add --no-cache tzdata

# copy go binary into container
COPY moneybags /app/moneybags

# expose port for API access
EXPOSE 8080

# set go binary as entrypoint
CMD ["/app/moneybags"]