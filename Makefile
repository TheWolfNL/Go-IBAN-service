dev: ## Start service for development
	go build -o service && . ./env-vars && ./service

help: # Help Command to output info on these make commands
	@printf "\033[33m%-30s\033[0m %s\n" 'Make commands:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
