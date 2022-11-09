FROM  zhhanscn/builder:golang-1.19.2-alpine-3.16 as  builder
ARG COMPRESS='false'
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -a -ldflags "-linkmode external -w -s -extldflags  -static " -gcflags "-l=4 -m=2 "  --tags netgo -o dist/serverd github.com/orangebees/go-layout/cmd/serverd  \
    &&  if [  ${COMPRESS} = 'true' ]; then upx -9 dist/serverd; fi




FROM zhhanscn/runtime-env:scratch-shanghai
EXPOSE 80/tcp
ENV ATREUGO_ADDR="0.0.0.0:80"
COPY  --from=builder  /app/dist/ /app/
CMD ["/app/serverd"]

