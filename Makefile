build:
	@docker build -t computer-club:alexeysavchuk .

run:
	@if [ -z "$(FILE)" ]; then \
		echo "No file provided. Please specify a file path using PATH=\"/path/to/file\""; \
		exit 1; \
	elif [ ! -f "$(FILE)" ]; then \
		echo "The specified path does not exist or is not a file: $(FILE)"; \
		exit 1; \
	fi; \
	docker run --rm --volume $(realpath $(FILE)):/app/data/file:ro computer-club:alexeysavchuk; \

clean:
	@docker rmi computer-club:alexeysavchuk
