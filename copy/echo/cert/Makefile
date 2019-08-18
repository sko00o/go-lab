all:
	# generatge self-signed certificate
	@# subj formate is:
	@# "/C=<Country Code>/ST=<State>/L=<City>/O=<Organization>/CN=<Common Name>"
	@openssl req -x509 -nodes -days 365 \
		-subj "/C=CN/ST=ZJ/L=HZ/O=fakeOrg/CN=localhost:8686" \
		-newkey rsa:2048 -keyout server.key \
		-out server.crt
	