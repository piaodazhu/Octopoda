#!/bin/bash
BUILD_VERSION=$(cat VERSION)
BUILD_TIME=$(date "+%F@%T")
BUILD_NAME="OCTOPODA"
COMMIT_SHA1=$(git rev-parse HEAD )

LD_FLAGS="-X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.BuildName=${BUILD_NAME} -X main.CommitID=${COMMIT_SHA1}"

echo "START BUILD VERSION ${BUILD_VERSION}"

mkdir -p dist

cd octl

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
tar -Jcvf octl_${BUILD_VERSION}_linux_amd64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <linux amd64>"
mv octl_${BUILD_VERSION}_linux_amd64.tar.xz ../dist/

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
tar -Jcvf octl_${BUILD_VERSION}_linux_arm7.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <linux arm7>"
mv octl_${BUILD_VERSION}_linux_arm7.tar.xz ../dist/

GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
tar -Jcvf octl_${BUILD_VERSION}_darwin_arm64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <darwin arm64>"
mv octl_${BUILD_VERSION}_darwin_arm64.tar.xz ../dist/

GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
tar -Jcvf octl_${BUILD_VERSION}_darwin_amd64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <darwin amd64>"
mv octl_${BUILD_VERSION}_darwin_amd64.tar.xz ../dist/

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
tar -Jcvf octl_${BUILD_VERSION}_windows_amd64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <windows amd64>"
mv octl_${BUILD_VERSION}_windows_amd64.tar.xz ../dist/

cd ../brain

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
tar -Jcvf brain_${BUILD_VERSION}_linux_amd64.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux amd64>"
mv brain_${BUILD_VERSION}_linux_amd64.tar.xz ../dist/

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
tar -Jcvf brain_${BUILD_VERSION}_linux_arm7.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux arm7>"
mv brain_${BUILD_VERSION}_linux_arm7.tar.xz ../dist/

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
tar -Jcvf brain_${BUILD_VERSION}_linux_arm64.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux arm64>"
mv brain_${BUILD_VERSION}_linux_arm64.tar.xz ../dist/

GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
tar -Jcvf brain_${BUILD_VERSION}_linux_i386.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux i386>"
mv brain_${BUILD_VERSION}_linux_i386.tar.xz ../dist/

GOOS=linux GOARCH=mips CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
tar -Jcvf brain_${BUILD_VERSION}_linux_mips.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux mips>"
mv brain_${BUILD_VERSION}_linux_mips.tar.xz ../dist/

cd ../tentacle

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
tar -Jcvf tentacle_${BUILD_VERSION}_linux_amd64.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux amd64>"
mv tentacle_${BUILD_VERSION}_linux_amd64.tar.xz ../dist/

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
tar -Jcvf tentacle_${BUILD_VERSION}_linux_arm7.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux arm7>"
mv tentacle_${BUILD_VERSION}_linux_arm7.tar.xz ../dist/

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
tar -Jcvf tentacle_${BUILD_VERSION}_linux_arm64.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux arm64>"
mv tentacle_${BUILD_VERSION}_linux_arm64.tar.xz ../dist/

GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
tar -Jcvf tentacle_${BUILD_VERSION}_linux_i386.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux i386>"
mv tentacle_${BUILD_VERSION}_linux_i386.tar.xz ../dist/

GOOS=linux GOARCH=mips CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
tar -Jcvf tentacle_${BUILD_VERSION}_linux_mips.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux mips>"
mv tentacle_${BUILD_VERSION}_linux_mips.tar.xz ../dist/

echo "ALL DONE!"