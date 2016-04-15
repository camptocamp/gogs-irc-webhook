gogs-irc-webhook:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $@ .
	strip $@
