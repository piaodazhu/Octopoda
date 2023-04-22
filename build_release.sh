#!/bin/bash
mkdir -p dist
cd octl

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o octl .
tar -Jcvf octl_1.0_linux_amd64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <linux amd64>"
mv octl_1.0_linux_amd64.tar.xz ../dist/

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o octl .
tar -Jcvf octl_1.0_linux_arm7.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <linux arm7>"
mv octl_1.0_linux_arm7.tar.xz ../dist/

GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o octl .
tar -Jcvf octl_1.0_darwin_arm64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <darwin arm64>"
mv octl_1.0_darwin_arm64.tar.xz ../dist/

GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o octl .
tar -Jcvf octl_1.0_darwin_amd64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <darwin amd64>"
mv octl_1.0_darwin_amd64.tar.xz ../dist/

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o octl .
tar -Jcvf octl_1.0_windows_amd64.tar.xz octl octl.yaml example/helloWorld &>/dev/null && echo "Successfully build octl <windows amd64>"
mv octl_1.0_windows_amd64.tar.xz ../dist/

cd ../brain

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o brain .
tar -Jcvf brain_1.0_linux_amd64.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux amd64>"
mv brain_1.0_linux_amd64.tar.xz ../dist/

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o brain .
tar -Jcvf brain_1.0_linux_arm7.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <linux arm7>"
mv brain_1.0_linux_arm7.tar.xz ../dist/

GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o brain .
tar -Jcvf brain_1.0_darwin_arm64.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <darwin arm64>"
mv brain_1.0_darwin_arm64.tar.xz ../dist/

GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o brain .
tar -Jcvf brain_1.0_darwin_amd64.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <darwin amd64>"
mv brain_1.0_darwin_amd64.tar.xz ../dist/

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o brain .
tar -Jcvf brain_1.0_windows_amd64.tar.xz brain brain.service brain.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build brain <windows amd64>"
mv brain_1.0_windows_amd64.tar.xz ../dist/

cd ../tentacle

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tentacle .
tar -Jcvf tentacle_1.0_linux_amd64.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux amd64>"
mv tentacle_1.0_linux_amd64.tar.xz ../dist/

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o tentacle .
tar -Jcvf tentacle_1.0_linux_arm7.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <linux arm7>"
mv tentacle_1.0_linux_arm7.tar.xz ../dist/

GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o tentacle .
tar -Jcvf tentacle_1.0_darwin_arm64.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <darwin arm64>"
mv tentacle_1.0_darwin_arm64.tar.xz ../dist/

GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o tentacle .
tar -Jcvf tentacle_1.0_darwin_amd64.tar.xz tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh &>/dev/null && echo "Successfully build tentacle <darwin amd64>"
mv tentacle_1.0_darwin_amd64.tar.xz ../dist/


echo "ALL DONE!"