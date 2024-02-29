# loki

通过API将日志提交到loki, 不推荐该方式, 推荐使用`promtail`来抓取日志

```yaml
# 日志组件配置
logger:
  level: info
  logfile:
    enable: true
    errorfilename: ./runtime/logs/cago.err.log
    filename: ./runtime/logs/cago.log
  # 不推荐该方式, 推荐使用`promtail`来抓取日志
  loki:
    level: info
    url: http://127.0.0.1:3100/loki/api/v1/push
    username: ""
    password: ""
```

## k8s helm promtail配置

拉取promtail的helm chart, `scrapeConfigs`添加job

```yaml
      - job_name: cago
        pipeline_stages:
          - cri: { }
        kubernetes_sd_configs:
          - role: pod
        relabel_configs:
          - source_labels:
              - __meta_kubernetes_pod_label_app_kubernetes_io_name
            action: keep
            regex: cago
          - source_labels:
              - __meta_kubernetes_pod_label_app_cago_io_name
            action: replace
            target_label: app
          - source_labels:
              - __meta_kubernetes_pod_label_app_cago_io_version
            action: replace
            target_label: version
          - source_labels:
              - __meta_kubernetes_pod_label_app_cago_io_environment
            action: replace
            target_label: env
          { { - toYaml .Values.config.snippets.common | nindent 4 } }
          { { - with .Values.config.snippets.extraRelabelConfigs } }
          { { - toYaml . | nindent 4 } }
          { { - end } }
```