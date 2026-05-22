#!/usr/bin/env sh
set -eu

mkdir -p dist

build() {
  os="$1"
  arch="$2"
  ext=""
  if [ "$os" = "windows" ]; then
    ext=".exe"
  fi

  out="dist/p2p_${os}_${arch}${ext}"
  echo "building ${out}"
  GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$out" ./cmd/p2p
}

build darwin arm64
build darwin amd64
build linux arm64
build linux amd64
build windows amd64

(
  cd dist
  shasum -a 256 p2p_* > checksums.txt
)
