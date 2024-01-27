FROM node:20

WORKDIR /app

COPY website/ website/

RUN cd  website && yarn
RUN cd  website && yarn build


FROM 		golang:1.21-alpine

RUN         adduser -h /home/role-finder -u 1000 -s /sbin/nologin -D role-finder && apk add ca-certificates

WORKDIR		/app
ADD         go.mod go.sum  /app/
RUN         go mod download

ADD         . /app/
RUN	        CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o role-finder .

FROM 		scratch
ARG		    VERSION


COPY --from=1 /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=1 /etc/passwd /etc/passwd
COPY --from=1 /etc/group /etc/group
COPY --from=1 /home/role-finder /home/role-finder


COPY --from=0 /app/website/dist /app/website/dist
COPY --from=1 /app/data /data
COPY --from=1 /app/role-finder /

USER role-finder


EXPOSE 8080
ENTRYPOINT ["/role-finder"]
CMD ["serve", "--from-file"]
