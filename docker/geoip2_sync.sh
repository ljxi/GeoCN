#!/bin/sh

# 检查是否提供了 LICENSE_KEY
if [ -z "$LICENSE_KEY" ]; then
  echo "错误：未设置LICENSE_KEY。请将其导出为环境变量。"
  exit 1
fi

# 设置 MaxMind 的下载地址
BASE_URL="https://download.maxmind.com/app/geoip_download"

echo "下载 GeoLite2-City.mmdb..."
curl -L -o "GeoLite2-City.mmdb.tar.gz" \
  "${BASE_URL}?edition_id=GeoLite2-City&license_key=${LICENSE_KEY}&suffix=tar.gz"

if [ $? -eq 0 ]; then
  tar -xvzf GeoLite2-City.mmdb.tar.gz --strip-components=1 && rm GeoLite2-City.mmdb.tar.gz
  echo "GeoLite2-City.mmdb 已成功下载并提取。"
else
  echo "GeoLite2-City.mmdb 下载失败。"
  exit 1
fi

echo "下载 GeoLite2-ASN.mmdb..."
curl -L -o "GeoLite2-ASN.mmdb.tar.gz" \
  "${BASE_URL}?edition_id=GeoLite2-ASN&license_key=${LICENSE_KEY}&suffix=tar.gz"

if [ $? -eq 0 ]; then
  tar -xvzf GeoLite2-ASN.mmdb.tar.gz --strip-components=1 && rm GeoLite2-ASN.mmdb.tar.gz
  echo "GeoLite2-ASN.mmdb 已成功下载并提取。"
else
  echo "GeoLite2-ASN.mmdb 下载失败。"
  exit 1
fi

echo "下载 GeoLite2-Country.mmdb..."
curl -L -o "GeoLite2-Country.mmdb.tar.gz" \
  "${BASE_URL}?edition_id=GeoLite2-Country&license_key=${LICENSE_KEY}&suffix=tar.gz"

if [ $? -eq 0 ]; then
  tar -xvzf GeoLite2-Country.mmdb.tar.gz --strip-components=1 && rm GeoLite2-Country.mmdb.tar.gz
  echo "GeoLite2-Country.mmdb 已成功下载并提取。"
else
  echo "GeoLite2-Country.mmdb 下载失败。"
  exit 1
fi

echo "下载 GeoCN.mmdb..."
curl -L -o "GeoCN.mmdb" \
  "http://github.com/ljxi/GeoCN/releases/download/Latest/GeoCN.mmdb"

if [ $? -eq 0 ]; then
  echo "GeoCN.mmdb 已成功下载。"
else
  echo "GeoCN.mmdb 下载失败。"
  exit 1
fi

echo "所有下载已完成。"
