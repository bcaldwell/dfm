#!/bin/bash

set -e
{
  _dfm_system_detector() {
  local os=$(uname -s | tr '[:upper:]' '[:lower:]')
  local arch
  case $(uname -m) in
    x86_64)
      arch="amd64"
      ;;
    *"arm"*)
      arch="arm"
      ;;
    *)
      arch="unsupported"
      ;;
  esac
  echo "${os}_${arch}"
  }

  LATEST_DFM_VERSION=$(curl -L -s -H 'Accept: application/json' https://github.com/benjamincaldwell/dfm/releases/latest | python -c 'import json,sys;obj=json.load(sys.stdin);print obj["tag_name"]')

  wget "https://github.com/benjamincaldwell/dfm/releases/download/${LATEST_DFM_VERSION}/dfm_$(_dfm_system_detector)" -O /tmp/dfm
  chmod +x /tmp/dfm
  echo "Moving binary to /usr/local/bin/"
  sudo mv /tmp/dfm /usr/local/bin/
}