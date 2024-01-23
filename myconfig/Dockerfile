FROM busybox:latest

COPY ControlDiskIO.sh /root/ControlDiskIO.sh
RUN chmod +x /root/ControlDiskIO.sh

ARG speed=10
ENV speed=$speed

CMD /root/ControlDiskIO.sh ${speed}