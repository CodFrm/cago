# 指标

使用普罗米修斯采集

采集器配置：

```
      - job_name: 'remote-cago'
        honor_labels: true

        kubernetes_sd_configs:(如果是远程的话,这样配置)
          - api_server: https://apiserver:6443  # apiserver 地址
            role: pod
            namespaces:
              names:
                - app
            bearer_token_file: /etc/secrets/remote/token
            tls_config:
              insecure_skip_verify: true

        scheme: https
        bearer_token_file: /etc/secrets/remote/token
        tls_config:
          insecure_skip_verify: true 
        # 主要配置是下面这段
        relabel_configs:
          - source_labels: [__meta_kubernetes_pod_label_app_kubernetes_io_name]
            action: keep
            regex: cago
          - action: labelmap
            regex: __meta_kubernetes_pod_label_(.+)
          - source_labels: [__meta_kubernetes_namespace,__meta_kubernetes_pod_name]
            separator: ;
            regex: (.*);(.*)
            target_label: __metrics_path__
            replacement: /api/v1/namespaces/$1/pods/$2:80/proxy/metrics
          - source_labels: [__address__]
            separator: ;
            regex: (.*)
            target_label: __address__
            replacement: remote:6443
            action: replace
```