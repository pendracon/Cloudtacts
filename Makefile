RUNNER_BIN=authrunnerexe
CC = go build
RUN = go run
ODIR = build/Cloudtacts
DDIR = dist/Cloudtacts
FLAGS = -ldflags="-s -w"
GOOS = linux

.PHONY: all authrunner buildir runner clean install localdeploy

all : buildir $(TARGET)

authrunner : buildir prep runner localdeploy

buildir:
	if test -n $(ODIR); then mkdir -p $(ODIR); fi
	if test -n $(DDIR); then mkdir -p $(DDIR); fi

runner:
	cp cmd/function/runner.go $(ODIR)
	cp cmd/function/function.go $(ODIR)
	GOOS=$(GOOS) $(CC) $(FLAGS) -o $(DDIR)/$(RUNNER_BIN) $(ODIR)/runner.go

localdeploy:
	cp -r config $(DDIR)
	#cd $(ODIR); $(RUN) runner.go
	cd $(DDIR); chmod +x $(RUNNER_BIN); ./$(RUNNER_BIN)

install : $(TARGET)
	install $(TARGET) ~/bin

prep:
	go mod tidy
	cp go.mod $(DDIR)
	cp go.sum $(DDIR)
	cp -r pkg $(ODIR)

clean :
	rm -rf $(DDIR) 2> /dev/null
	rm -rf $(ODIR) 2> /dev/null
