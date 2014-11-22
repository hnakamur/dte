#!/bin/sh
executable=dte
os_and_arches="darwin_amd64 linux_amd64 linux_386 windows_amd64 windows_386"

build() {
  mkdir -p download/$1
  local IFS=_
  set -- $1
  GOOS=$1 GOARCH=$2 go build -o download/${1}_$2/$executable
}

for os_and_arch in $os_and_arches; do
  build $os_and_arch &
done
wait
