.DEFAULT_GOAL = help
PROJECT	= mirea_ics

$(PROJECT):
	go build

.PHONY: build ## Build project
build: $(PROJECT)

.PHONY: all.ics
all.ics: build
	./$(PROJECT) -links ${PATTERN} \
	  | xargs -I % ./$(PROJECT) -file % -merge -all

%.ics: build
	./$(PROJECT) -links $(basename $@) \
	  | xargs -I % ./$(PROJECT) -file % $(basename $@)
	./$(PROJECT) -merge $(basename $@)

.PHONY: update
update: ${GROUP}.ics
	git add -f *.ics
	git stash
	git checkout --orphan ${GROUP} || git switch ${GROUP} 
	git checkout stash -- .
	git stash drop
	git reset
	git add -f *.ics

	git diff --staged --color \
	  | perl -nle 'print if /\e\[3[12]m/' \
	  | grep -vE 'DTSTAMP|CREATED|LAST-MODIFIED|UID|@@' \
	  > /dev/null

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: help
help: ## Display this
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	  | sort \
	  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[0;32m%-30s\033[0m %s\n", $$1, $$2}'
