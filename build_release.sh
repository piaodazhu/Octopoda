#!/bin/bash
VERSION=$(cat VERSION)
echo "START BUILD VERSION ${VERSION}"

mkdir -p dist
cd octl

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o octl .
tar -Jcvf octl_${VERSION}_linux_amd64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <linux amd64>"
mv octl_${VERSION}_linux_amd64.tar.xz ../dist/

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o octl .
tar -Jcvf octl_${VERSION}_linux_arm7.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <linux arm7>"
mv octl_${VERSION}_linux_arm7.tar.xz ../dist/

GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o octl .
tar -Jcvf octl_${VERSION}_darwin_arm64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <darwin arm64>"
mv octl_${VERSION}_darwin_arm64.tar.xz ../dist/

GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o octl .
tar -Jcvf octl_${VERSION}_darwin_amd64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <darwin amd64>"
mv octl_${VERSION}_darwin_amd64.tar.xz ../dist/

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o octl .
tar -Jcvf octl_${VERSION}_windows_amd64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <windows amd64>"
mv octl_${VERSION}_windows_amd64.tar.xz ../dist/

cd ../brain

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o brain .
tar -Jcvf brain_${VERSION}_linux_amd64.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux amd64>"
mv brain_${VERSION}_linux_amd64.tar.xz ../dist/

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o brain .
tar -Jcvf brain_${VERSION}_linux_arm7.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux arm7>"
mv brain_${VERSION}_linux_arm7.tar.xz ../dist/

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o brain .
tar -Jcvf brain_${VERSION}_linux_arm64.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux arm64>"
mv brain_${VERSION}_linux_arm64.tar.xz ../dist/

GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o brain .
tar -Jcvf brain_${VERSION}_linux_i386.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux i386>"
mv brain_${VERSION}_linux_i386.tar.xz ../dist/

GOOS=linux GOARCH=mips CGO_ENABLED=0 go build -o brain .
tar -Jcvf brain_${VERSION}_linux_mips.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux mips>"
mv brain_${VERSION}_linux_mips.tar.xz ../dist/

cd ../tentacle

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tentacle .
tar -Jcvf tentacle_${VERSION}_linux_amd64.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux amd64>"
mv tentacle_${VERSION}_linux_amd64.tar.xz ../dist/

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o tentacle .
tar -Jcvf tentacle_${VERSION}_linux_arm7.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux arm7>"
mv tentacle_${VERSION}_linux_arm7.tar.xz ../dist/

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o tentacle .
tar -Jcvf tentacle_${VERSION}_linux_arm64.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux arm64>"
mv tentacle_${VERSION}_linux_arm64.tar.xz ../dist/

GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o tentacle .
tar -Jcvf tentacle_${VERSION}_linux_i386.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux i386>"
mv tentacle_${VERSION}_linux_i386.tar.xz ../dist/

GOOS=linux GOARCH=mips CGO_ENABLED=0 go build -o tentacle .
tar -Jcvf tentacle_${VERSION}_linux_mips.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux mips>"
mv tentacle_${VERSION}_linux_mips.tar.xz ../dist/

echo "ALL DONE!"