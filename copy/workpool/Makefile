run:
	go build -o workpool
	d=$$(date +%s); \
		./workpool && \
		echo took $$(($$(date +%s)-d)) seconds.
