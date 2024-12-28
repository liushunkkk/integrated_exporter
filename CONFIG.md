# CONFIG

## etc.yaml

- `APP`: app名字

- `syntax[可选]`:  配置文件版本，一般为 V1 即可，暂时无用

- `server`:

  - `port[可选]`:  数字，端口，默认为 `6070`
  - `route[可选]`: 接口路径，默认为 `/metrics`
  - `interval[可选]`:  采集间隔，默认为 `5s`
  - `machineConfig`: 机器指标配置
    - `metrics[可选]`: 字符串数组，需要采集的机器指标，默认为 `[ cpu,memory,disk,process,network ]`
    - `mounts[可选]`: 字符串数组，需要采集指标的磁盘挂载点，默认为 `[ / ]`
    - `processes[可选]`: 字符串数组，需要采集metrics信息的机器进程，进程名，匹配方式采用 strings.Contains(实际程序名，配置值).
  - `gethServices[可选]`: 带metrics的区块链节点
    - `name`: 服务名称，唯一标识，将会在采集到的 metrics 补充`servicename=name`标签
    - `address`: metrics地址，如`http://127.0.0.1:6060/metrics`
    - `token[可选]`: 验证请求头，如 `Bearer 123456`
    - `timeout`: 请求超时时间，超时时间不得超过采集间隔`interval`
  -   `apiServices[可选]`: # 带metrics的后台服务
       - `name`: 服务名称
       - `address`: metrics地址
       - `token`: 验证请求头
       - `timeout`: 请求超时时间

  -   `httpServices[可选]`: 使用 http 探针访问的 http 服务，将会产生 `服务名_live_status` 指标
    - `name`: 服务名称
    - `address`:  请求接口地址，如 `http://127.0.0.1:8001/api/v1/user/list`
    - `method`: 请求方法，如 `GET`
    - `token[可选]`: 验证请求头
    - `body[可选]`: 请求体，如 `'{"message": "hello"}'`
    - `response[可选]`: 字符串，response验证，验证方式：strings.Contains(接口返回体字符串，配置值)
    - `timeout`: 请求超时时间
  -   `grpcServices[可选]`: 使用 grpc 探针访问的 grpc 服务
       - `name`: 服务名称
       - `address`: 请求地址，如 `127.0.0.1:8000`
       - `token`: 验证请求头
       - `rpcMethod`:  请求的服务方法，如 `hellopb.hello/SayHello`
       - `body`: 请求体
       - `response`: 字符串，response验证
       - `timeout`: 请求超时时间
  -   `processServices[可选]`: 需要验证存活状态的进程
       - `name`: 服务名称
       - `target`: 进程名，匹配规则：strings.Contains(实际进程名，配置值)



## .env.yaml

该文件是可选的

```yaml
SYNTAX_ENV: V1
```

编写环境配置，然后在 `etc.yaml` 中使用 `${}` 访问即可



## 区块链节点参考配置

```yaml
APP: integrated-exporter

syntax: V1

server:
  port: 6070
  route: /metrics
  interval: 5s # 替换为想要的采集时间
  machineConfig:
    metrics: [ cpu,memory,disk,process,network ]
    mounts: [ /node/data1,/node/data2 ] # 替换为节点使用的挂载点
  gethServices: # 链节点
    - name: chainNode
      address: http://127.0.0.1:6060/debug/metrics/prometheus # 替换为节点暴露的metrics地址
      timeout: 2s # 替换为想要的超时时间
  apiServices: # 后台服务
    - name: nodeApi
      address: http://127.0.0.1:6061/metrics # 替换为后台服务程序暴露的metrics地址
      timeout: 2s # 替换为想要的超时时间
```

