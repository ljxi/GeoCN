FROM python:3.9 as builder
WORKDIR /code

COPY ./requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# ===================================================

FROM python:3.9-slim
ENV TZ=Asia/Shanghai
WORKDIR /code
Expose 80

# 复制依赖和二进制文件到生产镜像中
COPY --from=builder /usr/local/lib/python3.9/site-packages /usr/local/lib/python3.9/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin
# 复制应用程序代码
COPY main.py .
COPY GeoCN.mmdb .
COPY GeoLite2-ASN.mmdb .
COPY GeoLite2-City.mmdb .
# COPY GeoLite2-Country.mmdb .
# 设置时区并安装必要的软件包
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    apt-get update && \
    apt-get install -y --no-install-recommends curl procps && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 使用 ENTRYPOINT 运行应用
ENTRYPOINT ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80", "--no-server-header", "--proxy-headers"]
