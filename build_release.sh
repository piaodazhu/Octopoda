#!/bin/bash
BUILD_VERSION=$(cat VERSION)
BUILD_TIME=$(date "+%F@%T")
BUILD_NAME="OCTOPODA"
COMMIT_SHA1=$(git rev-parse HEAD )

LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

echo "START BUILDING VERSION ${BUILD_VERSION}"

mkdir -p dist

cd octl

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
mkdir octl_${BUILD_VERSION}_linux_amd64 &>/dev/null && \
cp -r octl octl.yaml example/helloWorld octl_${BUILD_VERSION}_linux_amd64 &>/dev/null && \
tar -Jcvf octl_${BUILD_VERSION}_linux_amd64.tar.xz octl_${BUILD_VERSION}_linux_amd64 &>/dev/null && \
echo "Successfully build octl <linux amd64>" && \
mv octl_${BUILD_VERSION}_linux_amd64.tar.xz ../dist/ &>/dev/null && \
rm -rf octl_${BUILD_VERSION}_linux_amd64

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
mkdir octl_${BUILD_VERSION}_linux_arm7 &>/dev/null && \
cp -r octl octl.yaml example/helloWorld octl_${BUILD_VERSION}_linux_arm7 &>/dev/null && \
tar -Jcvf octl_${BUILD_VERSION}_linux_arm7.tar.xz octl_${BUILD_VERSION}_linux_arm7 &>/dev/null && \
echo "Successfully build octl <linux arm7>" && \
mv octl_${BUILD_VERSION}_linux_arm7.tar.xz ../dist/ &>/dev/null && \
rm -rf octl_${BUILD_VERSION}_linux_arm7

GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
mkdir octl_${BUILD_VERSION}_darwin_arm64 &>/dev/null && \
cp -r octl octl.yaml example/helloWorld octl_${BUILD_VERSION}_darwin_arm64 &>/dev/null && \
tar -Jcvf octl_${BUILD_VERSION}_darwin_arm64.tar.xz octl_${BUILD_VERSION}_darwin_arm64 &>/dev/null && \
echo "Successfully build octl <darwin arm64>" && \
mv octl_${BUILD_VERSION}_darwin_arm64.tar.xz ../dist/ &>/dev/null && \
rm -rf octl_${BUILD_VERSION}_darwin_arm64

GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
mkdir octl_${BUILD_VERSION}_darwin_amd64 &>/dev/null && \
cp -r octl octl.yaml example/helloWorld octl_${BUILD_VERSION}_darwin_amd64 &>/dev/null && \
tar -Jcvf octl_${BUILD_VERSION}_darwin_amd64.tar.xz octl_${BUILD_VERSION}_darwin_amd64 &>/dev/null && \
echo "Successfully build octl <darwin amd64>" && \
mv octl_${BUILD_VERSION}_darwin_amd64.tar.xz ../dist/ &>/dev/null && \
rm -rf octl_${BUILD_VERSION}_darwin_amd64

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl.exe .
mkdir octl_${BUILD_VERSION}_windows_amd64 &>/dev/null && \
cp -r octl.exe octl.yaml example/helloWorld octl_${BUILD_VERSION}_windows_amd64 &>/dev/null && \
tar -Jcvf octl_${BUILD_VERSION}_windows_amd64.tar.xz octl_${BUILD_VERSION}_windows_amd64 &>/dev/null && \
echo "Successfully build octl <windows amd64>" && \
mv octl_${BUILD_VERSION}_windows_amd64.tar.xz ../dist/ &>/dev/null && \
rm -rf octl_${BUILD_VERSION}_windows_amd64

cd ../brain

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
mkdir brain_${BUILD_VERSION}_linux_amd64 &>/dev/null && \
cp -r brain brain.service brain.yaml setup.sh uninstall.sh \
      brain_${BUILD_VERSION}_linux_amd64 &>/dev/null && \
tar -Jcvf brain_${BUILD_VERSION}_linux_amd64.tar.xz brain_${BUILD_VERSION}_linux_amd64 &>/dev/null && \
echo "Successfully build brain <linux amd64>" && \
mv brain_${BUILD_VERSION}_linux_amd64.tar.xz ../dist/ &>/dev/null && \
rm -rf brain_${BUILD_VERSION}_linux_amd64

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
mkdir brain_${BUILD_VERSION}_linux_arm7 &>/dev/null && \
cp -r brain brain.service brain.yaml setup.sh uninstall.sh \
      brain_${BUILD_VERSION}_linux_arm7 &>/dev/null && \
tar -Jcvf brain_${BUILD_VERSION}_linux_arm7.tar.xz brain_${BUILD_VERSION}_linux_arm7 &>/dev/null && \
echo "Successfully build brain <linux arm7>" && \
mv brain_${BUILD_VERSION}_linux_arm7.tar.xz ../dist/ &>/dev/null && \
rm -rf brain_${BUILD_VERSION}_linux_arm7

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
mkdir brain_${BUILD_VERSION}_linux_arm64 &>/dev/null && \
cp -r brain brain.service brain.yaml setup.sh uninstall.sh \
      brain_${BUILD_VERSION}_linux_arm64 &>/dev/null && \
tar -Jcvf brain_${BUILD_VERSION}_linux_arm64.tar.xz brain_${BUILD_VERSION}_linux_arm64 &>/dev/null && \
echo "Successfully build brain <linux arm64>" && \
mv brain_${BUILD_VERSION}_linux_arm64.tar.xz ../dist/ &>/dev/null && \
rm -rf brain_${BUILD_VERSION}_linux_arm64

GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
mkdir brain_${BUILD_VERSION}_linux_i386 &>/dev/null && \
cp -r brain brain.service brain.yaml setup.sh uninstall.sh \
      brain_${BUILD_VERSION}_linux_i386 &>/dev/null && \
tar -Jcvf brain_${BUILD_VERSION}_linux_i386.tar.xz brain_${BUILD_VERSION}_linux_i386 &>/dev/null && \
echo "Successfully build brain <linux i386>" && \
mv brain_${BUILD_VERSION}_linux_i386.tar.xz ../dist/ &>/dev/null && \
rm -rf brain_${BUILD_VERSION}_linux_i386

GOOS=linux GOARCH=mips CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
mkdir brain_${BUILD_VERSION}_linux_mips &>/dev/null && \
cp -r brain brain.service brain.yaml setup.sh uninstall.sh \
      brain_${BUILD_VERSION}_linux_mips &>/dev/null && \
tar -Jcvf brain_${BUILD_VERSION}_linux_mips.tar.xz brain_${BUILD_VERSION}_linux_mips &>/dev/null && \
echo "Successfully build brain <linux mips>" && \
mv brain_${BUILD_VERSION}_linux_mips.tar.xz ../dist/ &>/dev/null && \
rm -rf brain_${BUILD_VERSION}_linux_mips

cd ../tentacle

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
mkdir tentacle_${BUILD_VERSION}_linux_amd64 &>/dev/null && \
cp -r tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh \
      tentacle_${BUILD_VERSION}_linux_amd64 &>/dev/null && \
tar -Jcvf tentacle_${BUILD_VERSION}_linux_amd64.tar.xz tentacle_${BUILD_VERSION}_linux_amd64 &>/dev/null && \
echo "Successfully build tentacle <linux amd64>" && \
mv tentacle_${BUILD_VERSION}_linux_amd64.tar.xz ../dist/ &>/dev/null && \
rm -rf tentacle_${BUILD_VERSION}_linux_amd64

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
mkdir tentacle_${BUILD_VERSION}_linux_arm7 &>/dev/null && \
cp -r tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh \
      tentacle_${BUILD_VERSION}_linux_arm7 &>/dev/null && \
tar -Jcvf tentacle_${BUILD_VERSION}_linux_arm7.tar.xz tentacle_${BUILD_VERSION}_linux_arm7 &>/dev/null && \
echo "Successfully build tentacle <linux arm7>" && \
mv tentacle_${BUILD_VERSION}_linux_arm7.tar.xz ../dist/ &>/dev/null && \
rm -rf tentacle_${BUILD_VERSION}_linux_arm7

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
mkdir tentacle_${BUILD_VERSION}_linux_arm64 &>/dev/null && \
cp -r tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh \
      tentacle_${BUILD_VERSION}_linux_arm64 &>/dev/null && \
tar -Jcvf tentacle_${BUILD_VERSION}_linux_arm64.tar.xz tentacle_${BUILD_VERSION}_linux_arm64 &>/dev/null && \
echo "Successfully build tentacle <linux arm64>" && \
mv tentacle_${BUILD_VERSION}_linux_arm64.tar.xz ../dist/ &>/dev/null && \
rm -rf tentacle_${BUILD_VERSION}_linux_arm64

GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
mkdir tentacle_${BUILD_VERSION}_linux_i386 &>/dev/null && \
cp -r tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh \
      tentacle_${BUILD_VERSION}_linux_i386 &>/dev/null && \
tar -Jcvf tentacle_${BUILD_VERSION}_linux_i386.tar.xz tentacle_${BUILD_VERSION}_linux_i386 &>/dev/null && \
echo "Successfully build tentacle <linux i386>" && \
mv tentacle_${BUILD_VERSION}_linux_i386.tar.xz ../dist/ &>/dev/null && \
rm -rf tentacle_${BUILD_VERSION}_linux_i386

GOOS=linux GOARCH=mips CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
mkdir tentacle_${BUILD_VERSION}_linux_mips &>/dev/null && \
cp -r tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh \
      tentacle_${BUILD_VERSION}_linux_mips &>/dev/null && \
tar -Jcvf tentacle_${BUILD_VERSION}_linux_mips.tar.xz tentacle_${BUILD_VERSION}_linux_mips &>/dev/null && \
echo "Successfully build tentacle <linux mips>" && \
mv tentacle_${BUILD_VERSION}_linux_mips.tar.xz ../dist/ &>/dev/null && \
rm -rf tentacle_${BUILD_VERSION}_linux_mips


echo "ALL DONE!"