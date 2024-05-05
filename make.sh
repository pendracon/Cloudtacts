#!/bin/bash
CMD=$0
TARGET=$1

make_clean () {
	echo "Cleaning up..."

	go clean
	rm -r dist > /dev/null 2>&1
	rm -r build > /dev/null 2>&1
}

make_dist () {
	echo "Setting up dist..."

	mkdir -p dist/Cloudtacts/{bin,config,tools}
	mkdir -p dist/Cloudtacts/deploy
}

make_authrunner () {
	echo "Building auth function runner as dist/Cloudtacts/authrunnerexe..."

	cp cmd/function/*.go build/Cloudtacts/
	pushd build/Cloudtacts > /dev/null
	go mod tidy
	GOOS=linux go build -ldflags="-s -w" -o ../../dist/Cloudtacts/authrunnerexe runner.go
	popd > /dev/null
}

make_deploy () {
	make_clean
	make_prep
	make_authrunner
	make_dist

	echo "Building deployment package dist/Cloudtacts.tgz..."
	cp -r bin/* dist/Cloudtacts/bin/ > /dev/null 2>&1
	cp -r config/* dist/Cloudtacts/config/ > /dev/null 2>&1
	cp -r deploy/* dist/Cloudtacts/deploy/ > /dev/null 2>&1

	cd dist
	chmod a+x Cloudtacts/bin/*
	chmod a+x Cloudtacts/*exe
	tar cvzf Cloudtacts.tgz Cloudtacts/
	cd ..
}

make_prep () {
	mkdir -p build/Cloudtacts/pkg
	cp -r pkg/* build/Cloudtacts/pkg/
	cp go.{mod,sum} build/Cloudtacts/
}

usage () {
	echo "Usage: ${CMD} [auth deploy clean]"
	exit
}

case $TARGET in
	"auth") make_authrunner
		;;
	"deploy") make_deploy
		;;
	"clean") make_clean
		;;
	*) usage
		;;
esac

echo "Done!"
