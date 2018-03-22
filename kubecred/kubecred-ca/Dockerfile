# docker run -ti --rm -v /etc/kubernetes/pki/:/etc/kubernetes/pki/ kubecred-ca -days 30 $CSR
FROM alpine
RUN apk --no-cache add openssl
ADD kubecred-ca /
ENTRYPOINT ["/kubecred-ca"]
