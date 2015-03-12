export REPO = secretd
export NAME = $(REPO)
export GOPATH := $(shell readlink -f ../../../..)
export GOROOT := $(shell readlink -f $(GOPATH)/go)
export PATH := $(GOROOT)/bin:$(PATH)

include $(GOPATH)/Makefile.project

bin/%: build_data
	go install ./cmd/$*

.PHONY: jenkins
jenkins: dep-restore vet build_data test cmds fst-deb

# go-workspace's Makefile.project extracts the version number from version.go
# into the make variable $(VERSION).  We append to it the $(BUILD_NUMBER) that
# Jenkins sets and increments with every build attempt.
#
# Our fst-deb make target here, unlike the deb target in go-workspace's
# Makefile.project, prefixes the package name with fst- and adds the build
# number to the version,

.PHONY: fst-deb
ifeq ($(BUILD_NUMBER),)
fst-deb:
	echo >&2 "ERROR: Jenkin's BUILD_NUMBER is not set in the environment."
	false
else
fst-deb:
	false
	$(eval BUILD := $(shell mktemp -d /tmp/$(NAME)-build.XXXXXX))
	mkdir -p $(BUILD)/bin $(GOPATH)/deb
	cp -a templates $(BUILD)
	cp -a $(GOPATH)/bin/donner $(BUILD)/bin
	( cd $(GOPATH)/deb && /opt/fst-ffpm/bin/ffpm -s dir -t deb -n godonner --url https://github.com/fastly/goDonner -v $(VERSION)-$(BUILD_NUMBER) -C $(BUILD) --prefix /opt/goDonner bin/donner templates )
	rm -rf $(BUILD)
endif
