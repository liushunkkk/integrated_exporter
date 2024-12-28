# Integrated Exporter

对标`Prometheus`的集成监控程序，集机器监控，服务监控指标合并，服务探测于一体，使多服务的指标监控更简单：

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

> **Attention：如果不配置服务列表，默认情况下是个简易版本的`Node Exporter`.**



## 用法

### 在项目中使用

1、引入包

```shell
go get github.com/liushunkkk/integrated_exporter@latest
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
go install github.com/liushunkkk/integrated_exporter@latest
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
# 也支持传入制定的文件
integrated_exporter server --config=/path/to/your/etc.yaml
```

3、启动程序

```shell
integrated_exporter server # 部分参数也支持 flags 传入
```



#### 无 go 环境

可以根据系统自行下载[对应的 release 包](https://github.com/liushunkkk/integrated_exporter/releases)，然后执行（下面是mac的样例）：

```sh
# 下载包
curl -L -o integrated_exporter.tar.gz https://github.com/liushunkkk/integrated_exporter/releases/download/v0.1.2/integrated_exporter_Darwin_arm64.tar.gz
# 解压
tar -xzvf release-package.tar.gz
# 运行
./integrated_exporter server
# 同样的，如果有需要可以设置配置文件，或者传入 flags，比如：
./integrated_exporter server --port=6666 --route=/prometheus/metrics
```



#### Docker启动

使用默认配置

```sh
docker pull ghcr.io/liushunkkk/integrated_exporter:latest
docker run -d -p 6070:6070 --name integrated_exporter ghcr.io/liushunkkk/integrated_exporter
# or 
docker pull liushun311/integrated_exporter:latest
docker run -d -p 6070:6070 --name integrated_exporter liushun311/integrated_exporter
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
docker pull ghcr.io/liushunkkk/integrated_exporter:latest
docker run -d -p 6070:6070 -v running_dir/etc:/app/etc --name integrated_exporter ghcr.io/liushunkkk/integrated_exporter
# or 
docker pull liushun311/integrated_exporter:latest
docker run -d -p 6070:6070 -v running_dir/etc:/app/etc --name integrated_exporter liushun311/integrated_exporter
```



## Cookbook

### 配置项

样例配置可以前往[`etc/etc.example.yaml`](https://github.com/liushunkkk/integrated_exporter/blob/main/etc/etc.example.yaml)文件查看，也可以选择查看详细文档[CONFIG.md](https://github.com/liushunkkk/integrated_exporter/blob/main/CONFIG.md)。

[这里](https://github.com/liushunkkk/integrated_exporter/blob/main/CONFIG.md#%E5%8C%BA%E5%9D%97%E9%93%BE%E8%8A%82%E7%82%B9%E5%8F%82%E8%80%83%E9%85%8D%E7%BD%AE)给出了一个区块链节点的参考配置。

> 注意：
> 服务名不支持 `-` 拼接单词, 推荐使用驼峰形式。且为了防止指标冲突，请保证服务名是唯一的。
>
> ```yaml
> apiServices:
>   - name: portalApi // 推荐
> apiServices:
>   - name: portal-api // 不支持
> ```



### 环境变量

- 可以在 `etc/.env.yaml` 中创建环境变量，并在 `etc/etc.yaml` 中用 `${}` 使用.
- 在机器环境变量中设置以 `{EnvPrefix}_` 开头的变量，如 `{ENV_PREFIX}_A_B` 将映射到 `config.C.a.b`.
- 默认`EnvPrefix`为`INTEGRATEDEXPORTER`



### 指标设计

- 机器指标均以`machine`开头，如`machine_cpu_used`，`machine_mem_total`
- 探针存活服务状态暴露为`服务名_live_status`
- 请求到的其他程序的metrics，将会添加`servicename=服务名`标签以作区分，并且也会添加`服务名_live_status`作为服务状态指标
