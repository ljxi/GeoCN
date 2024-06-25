# GeoCN

每周更新中国大陆高精度IPV4+IPV6离线库，部分IP精确到区

下载地址：[GeoCN.mmdb](https://github.com/ljxi/GeoCN/releases/download/Latest/GeoCN.mmdb)

本项目仅包含[MaxMind](https://github.com/P3TERX/GeoLite.mmdb)显示为中国大陆的IP段

区域编码数据来自[modood/Administrative-divisions-of-China](https://github.com/modood/Administrative-divisions-of-China)

在本项目中，直辖市被视为省，直辖区被视为市，缺省信息为空字符串或数字0

### 字段示列

| 字段 | 示列 |
| --- | --- |
| isp | 中国移动 |
| net | 宽带 |
| province | 四川省 |
| provinceCode | 51 |
| city | 成都市 |
| cityCode | 5101 |
| districts | 武侯区 |
| districtsCode | 510107 |

### 在线演示

查询自己ipv4：`https://ipv4.netart.cn/`

查询自己ipv6：`https://ipv6.netart.cn/`

查询双栈：`https://ipvx.netart.cn/`

查询其他ipv4：`https://ipvx.netart.cn/222.137.1.1` `https://ipvx.netart.cn/?ip=222.137.1.1`

查询其他ipv6：`https://ipvx.netart.cn/240e:476::` `https://ipvx.netart.cn/?ip=240e:476::`

### Docker部署

`docker run -d -p 8000:80 netart/ipapi`

海外数据来自MaxMind，每天会自动拉取数据库

