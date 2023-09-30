#!/bin/bash
BUILD_VERSION=$(cat VERSION)
BUILD_TIME=$(date "+%F@%T")
BUILD_NAME=""
COMMIT_SHA1=$(git rev-parse HEAD )

echo "START BUILDING VERSION ${BUILD_VERSION}"

mkdir -p dist

cd octl

BUILD_OS=linux
BUILD_ARCH=amd64
BUILD_ARM=
BUILD_NAME=octl_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
# use these:
LD_FLAGS="-X main.BuildVersion=${BUILD_VERSION} -X main.BuildName=${BUILD_NAME} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"
# or: 
# VERSION_MSG="something-dev"
# LD_FLAGS="-X main.BuildVersion=${VERSION_MSG} -X main.BuildName=${BUILD_NAME} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r octl octl.yaml example/helloWorld ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=linux
BUILD_ARCH=arm
BUILD_ARM=7
BUILD_NAME=octl_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r octl octl.yaml example/helloWorld ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=darwin
BUILD_ARCH=arm64
BUILD_ARM=
BUILD_NAME=octl_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r octl octl.yaml example/helloWorld ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=darwin
BUILD_ARCH=amd64
BUILD_ARM=
BUILD_NAME=octl_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r octl octl.yaml example/helloWorld ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=windows
BUILD_ARCH=amd64
BUILD_ARM=
BUILD_NAME=octl_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o octl .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r octl octl.yaml example/helloWorld ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

cd ../brain

BUILD_OS=linux
BUILD_ARCH=amd64
BUILD_ARM=
BUILD_NAME=brain_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r brain brain.service brain.yaml setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=linux
BUILD_ARCH=arm
BUILD_ARM=7
BUILD_NAME=brain_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r brain brain.service brain.yaml setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=linux
BUILD_ARCH=arm64
BUILD_ARM=
BUILD_NAME=brain_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r brain brain.service brain.yaml setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=linux
BUILD_ARCH=386
BUILD_ARM=
BUILD_NAME=brain_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o brain .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r brain brain.service brain.yaml setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

cd ../tentacle

BUILD_OS=linux
BUILD_ARCH=amd64
BUILD_ARM=
BUILD_NAME=tentacle_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}


BUILD_OS=linux
BUILD_ARCH=arm
BUILD_ARM=7
BUILD_NAME=tentacle_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=linux
BUILD_ARCH=arm64
BUILD_ARM=
BUILD_NAME=tentacle_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=linux
BUILD_ARCH=386
BUILD_ARM=
BUILD_NAME=tentacle_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o tentacle .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=linux
BUILD_ARCH=mipsle
BUILD_ARM=
BUILD_NAME=tentacle_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"
# This is for Xiaomi 3 Pro
GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 GOMIPS=softfloat go build -ldflags "${LD_FLAGS}" -o tentacle .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r tentacle tentacle.service tentacle.yaml setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

cd ../httpNameServer

BUILD_OS=linux
BUILD_ARCH=amd64
BUILD_ARM=
BUILD_NAME=httpns_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o httpns .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r httpns httpns.service httpns.yaml setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

cd ../pakma

BUILD_OS=linux
BUILD_ARCH=amd64
BUILD_ARM=
BUILD_NAME=pakma_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o pakma .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r pakma pakma_tentacle.service pakma_brain.service setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

BUILD_OS=linux
BUILD_ARCH=arm
BUILD_ARM=7
BUILD_NAME=pakma_${BUILD_VERSION}_${BUILD_OS}_${BUILD_ARCH}
LD_FLAGS="-X main.BuildName=${BUILD_NAME} -X main.BuildVersion=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1}"

GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} GOARM=${BUILD_ARM} CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o pakma .
mkdir ${BUILD_NAME} &>/dev/null && \
cp -r pakma pakma_tentacle.service pakma_brain.service setup.sh uninstall.sh \
      ${BUILD_NAME} &>/dev/null && \
tar -Jcvf ${BUILD_NAME}.tar.xz ${BUILD_NAME} &>/dev/null && \
echo "Successfully build ${BUILD_NAME}" && \
mv ${BUILD_NAME}.tar.xz ../dist/ &>/dev/null && \
rm -rf ${BUILD_NAME}

echo "ALL DONE!"