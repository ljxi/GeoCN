FROM python:3.9 as py-builder
WORKDIR /code

COPY ./requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

FROM python:3.9-slim
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && apt-get install -y curl procps && rm -rf /var/lib/apt/lists/*

WORKDIR /code

COPY --from=py-builder /usr/local/lib/python3.9/site-packages /usr/local/lib/python3.9/site-packages
COPY --from=py-builder /usr/local/bin /usr/local/bin

COPY ./main.py .
COPY ./update_and_restart.sh .

CMD sh update_and_restart.sh