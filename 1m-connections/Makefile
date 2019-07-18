
MAKEFLAGS += --silent

# tcp
.PHONY: tcp-server-v1 tcp-server-v2 tcp-client

tcp-server-v1:
	go run tcp/server/v1/.

tcp-server-v2:
	cd tcp/server/v2/ && go run .

tcp-client:
	go run tcp/client/main.go -conn 500

# websocket

.PHONY: ws-server ws-client prom

ws-server:
	cd ws/server/ && go run .

ws-client:
	go run ws/client/main.go -conn 500

prom:
	prometheus
