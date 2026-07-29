.PHONY: all build test lint clean coverage examples examples-clean examples-test utilities utilities-clean install api-docs help

# All utilities. Mirrors gofi's UTILITIES := gofimac gofinet gofips one-for-one.
UTILITIES := goglmac goglnet goglps

# All examples
EXAMPLES := basic list reservations

# Install destination. Defaults to the user's ~/bin so no sudo is needed.
# Override with: make install INSTALL_DIR=/somewhere/else
INSTALL_DIR ?= $(HOME)/bin

.DEFAULT_GOAL := help

all: lint test build

# Build the module (compile-checks the whole tree) and the runnable utilities
# into bin/. Nothing is ever written to the repo root.
build: utilities
	go build ./...

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run ./...

clean: examples-clean utilities-clean
	go clean ./...
	rm -rf coverage.out coverage.html

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# === Examples ===

examples:
	@mkdir -p bin
	@for ex in $(EXAMPLES); do \
		echo "Building $$ex..."; \
		go build -o bin/$$ex ./examples/$$ex; \
	done
	@echo "All examples built in bin/"

examples-clean:
	rm -rf bin/

examples-test:
	@echo "Testing examples compile..."
	@for ex in $(EXAMPLES); do \
		echo "  Checking $$ex..."; \
		go build -o /dev/null ./examples/$$ex || exit 1; \
	done
	@echo "All examples compile successfully."

# === Utilities ===

utilities:
	@mkdir -p bin
	@for util in $(UTILITIES); do \
		echo "Building $$util..."; \
		go build -o bin/$$util ./utilities/$$util; \
	done
	@echo "All utilities built in bin/"

utilities-clean:
	rm -rf bin/

install: utilities
	@mkdir -p $(INSTALL_DIR)
	@for util in $(UTILITIES); do \
		echo "Installing $$util to $(INSTALL_DIR)/$$util"; \
		install -m 755 bin/$$util $(INSTALL_DIR)/$$util; \
	done
	@echo "All utilities installed to $(INSTALL_DIR)."
	@case ":$$PATH:" in *":$(INSTALL_DIR):"*) ;; *) echo "Note: $(INSTALL_DIR) is not on your PATH.";; esac

# Regenerate docs/api/ from GL.iNet's API description. The description is fetched
# to /tmp rather than vendored: it comes from a GPL-3.0 package and this project is
# MIT.
api-docs:
	./scripts/fetch-api-description.sh /tmp/gl-api-description.json
	python3 scripts/generate-api-docs.py /tmp/gl-api-description.json

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Main targets:"
	@echo "  all           Run lint, test, and build"
	@echo "  build         Build the module and utilities into bin/"
	@echo "  test          Run all tests"
	@echo "  lint          Run linter"
	@echo "  clean         Clean all build artifacts"
	@echo "  coverage      Generate coverage report"
	@echo ""
	@echo "Example targets:"
	@echo "  examples        Build all examples to bin/"
	@echo "  examples-clean  Remove example binaries"
	@echo "  examples-test   Verify all examples compile"
	@echo ""
	@echo "Utility targets:"
	@echo "  utilities       Build all utilities to bin/"
	@echo "  utilities-clean Remove utility binaries"
	@echo "  install         Build and install utilities to ~/bin (override INSTALL_DIR)"
	@echo ""
	@echo "Documentation targets:"
	@echo "  api-docs        Regenerate docs/api/ from GL.iNet's API description"
