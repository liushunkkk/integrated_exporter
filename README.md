# Integrated Exporter

对标`Prometheus`的集成监控程序：

- 在项目内引入，轻松创建各类监控指标：
  - Counter
  - Gauge
  - Summary
  - Histogram
- 支持`machine`基础监控指标：
  - CPU：CPU核心数、CPU使用率...
  - Memory：主存以及交换内存的使用情况...
  - Disk：指定磁盘挂载点的使用情况...
  - Process：进程数、指定进程的详细信息...
  - Network：网络连接数、网络读写信息...
  


- 配置服务列表，集成多服务指标，支持如下服务探测：
  - 无监控指标的`HTTP`服务：配置`Rest`探测接口，并生成`live_status`指标
  - 无监控指标的`gRPC`服务：配置`gRPC`探测接口，并生成`live_status`指标
  - 有监控指标的`API`服务：配置指标地址，生成`live_status`指标，并合并`API`服务的指标
  - 有监控指标的`Geth`服务：配置指标地址，生成`live_status`指标，并合并`Geth`服务的指标

> **Attention：如果不配置服务列表，默认情况下是个简单版本的`Node Exporter`.**



## 用法

### 在项目中使用

1、引入包

```shell
go get github.com/liushunking/integrated_exporter@v1.0.0
```

2、创建配置，并启动程序

```go
cfg := config.ServerConfig{
    Port: 6070,
    Interval: "5s",
    Route: "/metrics",
}
server.Run(cfg, nil, nil) // 传入 nil 则使用默认的 IRegistry 和 MetricsHandler
```

3、编写指标

```go
// counter
c := metricx.GetOrRegisterICounter(&metricx.IOpts{
    Namespace: "request",
    Name:      "success_count",
    Labels:    map[string]string{"method": "userlist"},
}, nil)
c.Inc()
// gauge
g := metricx.GetOrRegisterICounter(&metricx.IOpts{
    Namespace: "subscription",
    Name:      "transaction_stream",
    Labels:    map[string]string{"ip": "0.0.0.0"},
}, nil)
g.Inc()
g.Dec()
g.Set(10)
// summary
opts := metricx.NewRecommendISummaryOpts(&metricx.IOpts{
    Namespace: "a",
    Name:      "test_summary",
})
summary := metricx.GetOrRegisterISummary(opts, nil)
summary.Observe(100)
// histogram
opts := &metricx.IHistogramOpts{
  &IOpts{
      Namespace: "b",
      Name:      "test_histogram",
  },
  []float64{10, 20, 30, 40},
}
histogram := metricx.GetOrRegisterIHistogram(opts, nil)
histogram.Observe(15)
```

### 作为程序使用

#### 有 go 环境

1、下载程序

```shell
go install github.com/liushunking/integrated_exporter@v1.0.0
```

2、创建配置文件（如果需要）

```shell
# 进入你的执行目录
cd running_dir
# 创建配置文件夹
mkdir etc
# 编写配置文件
touch etc/etc.yaml
# 如果需要，可以编写环境变量配置文件
touch etc/.env.yaml
```

3、启动程序

```shell
integrated_exporter server # 部分参数也支持 flags 传入
```



#### 无 go 环境

可以根据系统自行下载[对应的 release 包](https://github.com/liushunking/integrated_exporter/releases)。



#### Docker启动

使用默认配置

```sh
docker pull ghcr.io/liushunking/integrated_exporter:v1.0.0
docker run -d -p 6070:6070 integrated_exporter
```

使用自定义配置

```sh
# 进入你的执行目录
cd running_dir
# 创建配置文件夹
mkdir etc
# 编写配置文件
touch etc/etc.yaml
# 如果需要，可以编写环境变量配置文件
touch etc/.env.yaml
docker pull ghcr.io/liushunking/integrated_exporter:v1.0.0
docker run -d -p 6070:6070 -v running_dir/etc:/etc integrated_exporter
```



## Cookbook

### 配置项

样例配置可以前往[`etc/etc.yaml`](https://github.com/liushunking/integrated_exporter/blob/main/etc/etc.yaml)文件查看.

配置默认值如下：

- Port: `6070`
- Route: `/metrics`
- Interval: `5s`
- MachineConfig:
  - Metrics: `[ cpu,memory,disk,process,network ]`
  - Mounts: `[ / ]`
- GethServices: `nil`
- ApiServices: `nil`
- HttpServices: `nil`
- GrpcServices: `nil`
- ProcessServices: `nil`



### 环境变量

- 可以在 `etc/.env.yaml` 中创建环境变量，并在 `etc/etc.yaml` 中用 `${}` 使用.
- 在机器环境变量中设置以 `{EnvPrefix}_` 开头的变量，如 `{ENV_PREFIX}_A_B` 将映射到 `config.C.a.b`.
- 默认`EnvPrefix`为`INTEGRATEDEXPORTER`

