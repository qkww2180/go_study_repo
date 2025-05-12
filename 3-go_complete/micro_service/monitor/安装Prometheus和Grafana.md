# Windows安装Prometheus
1. 下载链接https://github.com/prometheus/prometheus/releases/download/v2.45.0/prometheus-2.45.0.windows-amd64.zip  
2. 解压后编辑prometheus.yml文件,在scrpape_configs域添加自己的job:  
    ```
    - job_name: "golang_upgrading"
        static_configs:
        - targets: ["localhost:5678"]
        metrics_path: "/metrics"
    ```
    prometheus会周期性(高频)性的访问localhost:5678/metrics接口，拉取监控数据，存入数据库  
3. 双击prometheus.exe启动数据库  
# Windows安装Grafana 
1. 下载链接https://grafana.com/grafana/download?platform=windows，下载 Standalone Windows Binaries  
2. 解压后双击bin目录下的grafana-server.exe
3. 然后在浏览器打开http://localhost:3000。初始账户密码为admin:admin，登录后会强制修改admin的密码
