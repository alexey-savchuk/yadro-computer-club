image_name = "computer-club:alexeysavchuk"

build:
	@docker build -t $(image_name) .

run:
	@if [ -z "$(FILE)" ]; then \
		echo "No file provided. Please specify a file path using FILE=\"/path/to/file\""; \
		exit 1; \
	elif [ ! -f "$(FILE)" ]; then \
		echo "The specified path does not exist or is not a file: $(FILE)"; \
		exit 1; \
	fi; \
	docker run --rm --volume $(realpath $(FILE)):/app/file:ro $(image_name); \

clean:
	@docker rmi $(image_name)
