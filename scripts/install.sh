#!/bin/bash

package_name=$1
if [[ -z "${package_name}" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

version=$2
if [[ -z "${version}" ]]; then
  version="1.0.0"
fi

os="`uname`"
case ${os} in
  'Linux')
    os='linux'
    ;;
  'FreeBSD')
    os='freebsd'
    ;;
  'WindowsNT')
    os='windows'
    ;;
  'Darwin')
    os='darwin'
    ;;
  'Sunos')
    os='solaris'
    ;;
  *)
    echo "Unsupported os: ${os}"
    exit 1
    ;;
esac

arch="`uname -m`"
case ${arch} in
  'x86_64')
    arch='amd64'
    ;;
  'i?86')
    arch='386'
    ;;
  *)
    echo "Unsupported architecture: ${arch}"
    exit 1
    ;;
esac

plugin_dir=~/.terraform.d/plugins/${os}_${arch}
mkdir -p "${plugin_dir}"

cp "bin/${os}_${arch}/${package_name}_v${version}" "${plugin_dir}/${package_name}_v${version}"

