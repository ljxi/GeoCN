#!/bin/sh

# 设置 TERM 信号处理，确保 gunicorn 能被正常关闭
PIDFILE="/tmp/gunicorn.pid"

# trap '...' TERM 会在 shell 收到 SIGTERM 时触发
trap "echo '收到 SIGTERM，正在关闭 Gunicorn...'; kill -TERM \$(cat ${PIDFILE}); wait \$(cat ${PIDFILE})" TERM

# 后台更新进程 
update_loop() {
    sleep 86400
    
    while true; do
        date
        echo "开始更新 GeoLite 数据库..."
        curl  -L -o "GeoLite2-City.mmdb" "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb"
        curl  -L -o "GeoLite2-ASN.mmdb" "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-ASN.mmdb"
        curl  -L -o "GeoCN.mmdb" "http://github.com/ljxi/GeoCN/releases/download/Latest/GeoCN.mmdb"
        
        echo "数据库更新完毕，正在重启worker进程"
        

        # 平滑地重新加载 worker 进程，不会导致服务中断
        kill -HUP $(cat ${PIDFILE})
        
        sleep 86400
    done
}


# 启动时先下载一次数据库
date
echo "正在执行首次数据库下载..."
curl  -L -o "GeoLite2-City.mmdb" "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb"
curl  -L -o "GeoLite2-ASN.mmdb" "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-ASN.mmdb"
curl  -L -o "GeoCN.mmdb" "http://github.com/ljxi/GeoCN/releases/download/Latest/GeoCN.mmdb"

# 在后台启动更新循环
update_loop &

#  在前台启动 Gunicorn workers 设置标准为CPU核心数*2+1
echo "启动 Gunicorn 主进程..."
exec gunicorn main:app \
    --workers 3 \
    --worker-class uvicorn.workers.UvicornWorker \
    --bind '[::]:8080' \
    --pid ${PIDFILE}

