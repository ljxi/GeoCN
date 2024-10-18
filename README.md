# GeoCN

- 每月更新中国大陆高精度IPV4+IPV6离线库，部分IP精确到区
- 每月也会自动更新Docker镜像版本，[geo-cn](https://hub.docker.com/r/fc6a1b03/geo-cn)，使用`latest`即可

- 国内区域`mmdb`：[GeoCN.mmdb](https://github.com/ljxi/GeoCN/releases/download/Latest/GeoCN.mmdb)
- 其他区域`mmdb`[MaxMind](https://www.maxmind.com)显示为中国大陆的IP段
- 国内区域编码数据来自[modood/Administrative-divisions-of-China](https://github.com/modood/Administrative-divisions-of-China)

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
