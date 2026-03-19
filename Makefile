.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: arch-install
arch-install: ## Build Arch Linux package and clean up afterwards.
	$(MAKE) clean
	cd distribution/arch-linux && makepkg -sri
	$(MAKE) clean

.PHONY: arch-install-clean
clean: ## Remove all files except PKGBUILD.
	@rm -rf distribution/arch-linux/src/ distribution/arch-linux/pkg/ distribution/arch-linux/*.pkg.tar.zst distribution/arch-linux/*.pkg.tar.xz distribution/arch-linux/*.tar.gz
