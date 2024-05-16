RUNNER_BIN=authrunnerexe
CC = go build
RUN = go run
CLEAN = go clean
TEST = go test
ODIR = build/Cloudtacts
DDIR = dist/Cloudtacts
FLAGS = -ldflags="-s -w"
GOOS = linux

.PHONY: all authrunner buildir runner clean install localdeploy test

all : clean test buildir prep runner localdeploy

authrunner : buildir prep runner localdeploy

buildir:
	if test -n $(ODIR); then mkdir -p $(ODIR); fi
	if test -n $(DDIR); then mkdir -p $(DDIR); fi

runner:
	cp cmd/auth/function.go $(ODIR)
	cp cmd/runner/runner.go $(ODIR)
	GOOS=$(GOOS) $(CC) $(FLAGS) -o $(DDIR)/$(RUNNER_BIN) $(ODIR)/*.go

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

test:
	$(TEST) ./pkg/auth
	$(TEST) ./pkg/config

clean :
	$(CLEAN)
	rm -rf $(DDIR) 2> /dev/null
	rm -rf $(ODIR) 2> /dev/null
