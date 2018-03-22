FROM alpine
RUN apk --no-cache add openssl
ADD kubecred-ca /
ENTRYPOINT ["/kubecred-ca"]
