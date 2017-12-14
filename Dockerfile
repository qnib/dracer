FROM qnib/alplain-golang AS build

WORKDIR /usr/local/src/github.com/qnib/dracer/
COPY . ./
RUN go install

FROM qnib/alplain-init

COPY --from=build /usr/local/bin/dracer /usr/bin/
COPY start.sh /opt/dracer/bin/start.sh
CMD ["/opt/dracer/bin/start.sh"]
