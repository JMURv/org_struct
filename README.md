## Configuration
- Create your own `.env` based on `.env.example`:
```shell
cp build/configs/envs/.env.example build/configs/envs/.env.dev && \
cp build/configs/envs/.env.example build/configs/envs/.env.prod
```

### Docker
In the root folder run:
```shell
task dc-dev-build
```

Or run dev with observation containers:
```shell
task dc-dev-obs
```

Observe profile starts svcs like: `prometheus`, `jaeger`, `node-exporter`, `grafana` and etc.

Services are available at:

| Сервис           | Адрес                  |
|------------------|------------------------|
| App (HTTP)       | http://localhost:8080  |
| App (GRPC)       | http://localhost:50050 |
| App (PROMETHEUS) | http://localhost:8085  |
| Prometheus       | http://localhost:9090  |
| Node-exporter    | http://localhost:9100  |
| Jaeger           | http://localhost:16686 |
| Loki             | http://localhost:3100  |
| Grafana          | http://localhost:3000  |

More information could be found inside `compose.yaml`.

---

### Tests

Run `task t`

___

### K8s

- Create your own `configMap` and `secretMap` based on examples:
```shell
cp build/k8s/cfg/cfg.example.yaml build/k8s/cfg/cfg.yaml && \
cp build/k8s/cfg/secret.example.yaml build/k8s/cfg/secret.yaml 
```

Apply manifests:
```shell
task k-up
```

Shutdown manifests:
```shell
task k-down
```
