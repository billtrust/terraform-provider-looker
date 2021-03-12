VERSION = $(shell cat VERSION | tr -d '\n')
EXE = terraform-provider-looker_v$(VERSION)

all : $(EXE)

$(EXE) : terraform-provider-looker VERSION
	cp -v terraform-provider-looker $(EXE)

terraform-provider-looker : main.go looker/*.go go.mod
	go get
	go build

clean :
	rm -f $(EXE)
	rm -f terraform-provider-looker
	rm -rf build

install : $(EXE) build/TARGET
	mkdir -pv $(shell cat build/TARGET)
	cp -v $(EXE) $(shell cat build/TARGET)/$(EXE)

uninstall : build/TARGET
	rm -f $(shell cat build/TARGET)/$(EXE)

go.mod :
	go mod init github.com/Foxtel-DnA/terraform-provider-looker

reinit-module : clean-module go.mod

clean-module: clean
	rm -f go.mod

build:
	mkdir -pv build

build/PLATFORM: build bin/platform.sh
	./bin/platform.sh > build/PLATFORM

build/TARGET: build/PLATFORM VERSION
	echo $${HOME}/.terraform.d/plugins/terraform.foxtel.com/foxtel/looker/$(VERSION)/$(shell cat build/PLATFORM) > build/TARGET

.PHONY : clean clean-module install uninstall
