# GeoCN

随缘更新中国大陆高精度IPV4+IPV6离线库，部分IP精确到区

本项目数据来源于二分算法扫描公开免费高质量的API接口，并对数据进行清洗整理得到最终数据，更新周期取决于作者收集到的上游API的精度与更新周期，如果你有合适的API，欢迎提供

下载地址：[GeoCN.mmdb](https://github.com/ljxi/GeoCN/releases/latest/download/GeoCN.mmdb)

本项目仅包含[MaxMind](https://github.com/P3TERX/GeoLite.mmdb)显示为中国大陆的IP段

区域编码数据来自[xiangyuecn/AreaCity-JsSpider-StatsGov](https://github.com/xiangyuecn/AreaCity-JsSpider-StatsGov)

数据字段参照releases中的数据说明

### 在线演示

查询自己ipv4：`https://ipv4.netart.cn/`

查询自己ipv6：`https://ipv6.netart.cn/`

查询双栈：`https://ip.netart.cn/`

查询其他ipv4：`https://ip.netart.cn/222.137.1.1` `https://ip.netart.cn/?ip=222.137.1.1`

查询其他ipv6：`https://ip.netart.cn/240e:476::` `https://ip.netart.cn/?ip=240e:476::`

### Docker部署

`docker run -d -p 8000:80 netart/ipapi`

docker容器内置了所需的数据，不会自动更新，重新拉取镜像即可更新数据
