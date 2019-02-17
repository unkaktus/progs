# docker run -d --restart=always --cap-add=SYS_TIME nogoegst/fixtime -i 1s
FROM scratch
COPY fixtime /
ENTRYPOINT ["/fixtime"]
