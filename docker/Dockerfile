ARG repo
FROM ${repo}/library/alpine:amd64-3.12.3

EXPOSE 11408

RUN mkdir -p /opt/govirt && mkdir -p /etc/govirt
COPY govirt /opt/govirt
COPY govirt.yml /etc/govirt.conf
CMD /opt/govirt --config /etc/govirt.conf