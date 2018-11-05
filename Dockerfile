FROM index.alauda.cn/alaudaorg/alaudabase:alpine-supervisor-migrate-1
LABEL maintainer="Alauda Trouble Shooting"

COPY alauda_trouble_shooting /bin/alauda_trouble_shooting

EXPOSE      6667
ENTRYPOINT  [ "/bin/alauda_trouble_shooting" ]
