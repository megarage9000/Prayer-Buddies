# This is a comment
# Use a lightweight debian os
# as the base image
FROM debian:stable-slim

COPY Prayer-Buddies /bin/Prayer-Buddies

CMD ["/bin/Prayer-Buddies"]
