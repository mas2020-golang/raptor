# Include this fragment or copy the target into your Makefile.
# Usage:
#   make release
#   make release PUSH=1   # optionally push after tagging

RELEASE_SCRIPT := ./script/release.sh

.PHONY: release
release:
	@echo "==> Running release flow"
	@[ -x $(RELEASE_SCRIPT) ] || (echo "Making $(RELEASE_SCRIPT) executable"; chmod +x $(RELEASE_SCRIPT))
	@$(RELEASE_SCRIPT)
ifeq ($(PUSH),1)
	@git push origin main
	@# Push the newest tag as well
	@git push --follow-tags origin main
endif
