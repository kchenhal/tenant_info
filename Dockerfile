FROM debian:stable-slim


ADD tenant_info /root/

WORKDIR /root

CMD ["/root/tenant_info"]




