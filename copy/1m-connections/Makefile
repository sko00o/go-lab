
MAKEFLAGS += --silent

# tcp
.PHONY: tcp-server-v1 tcp-server-v2 tcp-client

tcp-server-v1:
	go run tcp/server/v1/main.go

tcp-server-v2:
	go run tcp/server/v2/main.go -c 5

tcp-server-v3:
	go run tcp/server/v3/main.go -c 5

tcp-client:
	go run tcp/client/main.go -conn 500 -c 10

# websocket

.PHONY: ws-server ws-client prom

ws-server:
	cd ws/server/ && go run .

ws-client:
	go run ws/client/main.go -conn 500

prom:
	prometheus
