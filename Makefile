VERSION = $(shell cat VERSION | tr -d '\n')
PLATFORM = $(shell uname -s | tr '[A-Z]' '[a-z]')_$(shell uname -m)
TARGET = $${HOME}/.terraform.d/plugins/terraform.foxtel.com/foxtel/looker/$(VERSION)/$(PLATFORM)
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

install : $(EXE)
	mkdir -pv $(TARGET)
	cp -v $(EXE) $(TARGET)/$(EXE)

uninstall :
	rm -f $(TARGET)/$(EXE)

go.mod :
	go mod init github.com/billtrust/terraform-provider-looker

reinit-module : clean-module go.mod

clean-module: clean
	rm -f go.mod

.PHONY : clean clean-module install uninstall
