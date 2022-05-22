FROM scratch

ENV APP=/usr/local/bin/vwap-engine \
    USER_UID=1101 \
    USER_NAME=vwap-engine

# install engine binary
COPY build/_output/bin/vwap-engine ${APP}
USER ${USER_UID}:${USER_UID}
CMD [ "./${APP}" ]
