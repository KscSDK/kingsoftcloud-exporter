# 金山云监控 Exporter

通过kingsoft cloud exporter将云监控支持的产品监控指标自动批量导出  

## 一、支持的产品列表
产品     | 命名空间 |支持的指标|
--------|---------|----------
云服务器 | KEC |[指标详情](https://docs.ksyun.com/documents/26#one)
弹性IP  | EIP |[指标详情](https://docs.ksyun.com/documents/26#two)
NAT    | NAT |[指标详情](https://docs.ksyun.com/documents/26#three)
负载均衡 | SLB |[指标详情](https://docs.ksyun.com/documents/26#six)
4层监听器 | LISTENER |[指标详情](https://docs.ksyun.com/documents/26#seven)
7层监听器 | LISTENER7 |[指标详情](https://docs.ksyun.com/documents/26#seven)
对等连接 | PEER |[指标详情](https://docs.ksyun.com/documents/26#eight)
共享带宽 | BWS |[指标详情](https://docs.ksyun.com/documents/26#nine)
专线网关 | DCGW|[指标详情](https://docs.ksyun.com/documents/26#seventeen)
关系型数据库 | KRDS |[指标详情](https://docs.ksyun.com/documents/26#five)
云数据库Redis | KCS |[指标详情](https://docs.ksyun.com/documents/26#four)

`后续会有更多的产品支持`

## 二、快速开始

### 1.构建

```shell
git clone https://github.com/KscSDK/kingsoftcloud-exporter.git
go build cmd/ksc-exporter/ksc_exporter.go
```

或从release列表获取预编译的二进制。

### 2. 配置

配置中主要分为如下两块：
- 云API的`credential`认证信息
- 产品`products`指标、实例导出信息

可在Git仓库中 `examples` 文件里, 找到含有支持的产品配置模版样例，用作参考。

配置参数说明如下：

```yaml
// 云监控拉取指标数据限制, 官方默认限制最大20qps
rate_limit: 15                                   

// 授权配置
credential:
  access_key: <YOUR_ACCESS_KEY>                  // 必须, 云API的SecretId
  secret_key: <YOUR_ACCESS_SECRET>               // 必须, 云API的SecretKey
  region: <REGION>                               // 必须, 实例所在区域信息
  
// 整个产品纬度配置, 每个产品一个item
// 特别需要注意：如果当前产品下有大量实例资源，不建议一次配置多个产品线，目前限制一个Expoter最多配置三个产品
products:
  - namespace: QCE/CMONGO                        // 必须, 产品命名空间
    all_metrics: true                            // 常用, 推荐开启, 导出支持的所有指标
    all_instances: true                          // 常用, 推荐开启, 导出该region下的所有实例
    only_include_metrics: [Inserts]              // 可选, 只导出这些指标, 配置时all_metrics失效
    statistics_types: [avg]                      // 可选, 拉取N个数据点, 再进行 max、min、avg计算, 默认max取最新值
    period_seconds: 60                           // 可选, 指标统计周期
    range_seconds: 300                           // 可选, 选取时间范围, 开始时间=now-range_seconds, 结束时间=now
    delay_seconds: 60                            // 可选, 时间偏移量, 结束时间=now-delay_seconds
    reload_interval_minutes: 60                  // 可选, 周期reload实例列表, 建议频率不要太频繁

// 单个指标纬度配置, 每个指标一个item
metrics:
  - tc_namespace: QCE/CMONGO                     // 产品命名空间, 同namespace
    tc_metric_name: Inserts                      // 云监控定义的指标名
    tc_metric_rename: Inserts                    // 导出指标的显示名
    tc_metric_name_type: 1                       // 可选，导出指标的名字格式化类型, 1=大写转小写加下划线, 2=转小写; 默认1
    tc_labels: [InstanceName]                    // 可选, 将实例的字段作为指标的lables导出
    tc_myself_dimensions:                        // 可选, 同custom_query_dimensions
    tc_statistics: [Avg]                         // 可选, 同statistics_types
    period_seconds: 60                           // 可选, 同period_seconds
    range_seconds: 300                           // 可选, 同range_seconds
    delay_seconds: 60                            // 可选, 同delay_seconds
```


**特殊说明:**

1. **custom_query_dimensions**  
   每个实例的纬度字段信息, 可从对应的云监控产品指标文档查询, 如mongo支持的纬度字段信息可由[云监控指标详情](https://cloud.tencent.com/document/product/248/45104#%E5%90%84%E7%BB%B4%E5%BA%A6%E5%AF%B9%E5%BA%94%E5%8F%82%E6%95%B0%E6%80%BB%E8%A7%88) 查询
2. **extra_labels**  
   每个导出metric的labels还额外上报实例对应的字段信息, 实例可选的字段列表可从对应产品文档查询, 如mongo实例支持的字段可从[实例查询api文档](https://cloud.tencent.com/document/product/240/38568) 获取, 目前只支持str、int类型的字段
3. **period_seconds**  
   每个指标支持的时间纬度统计, 一般支持60、300秒等, 具体可由对应产品的云监控产品指标文档查询, 如mongo可由[指标元数据查询](https://cloud.tencent.com/document/product/248/30351) , 假如不配置, 使用默认值(60), 假如该指标不支持60, 则自动使用该指标支持的最小值
4. **credential**  
   SecretId、SecretKey、Region可由环境变量获取

  ```bash
  export TENCENTCLOUD_SECRET_ID="YOUR_ACCESS_KEY"
  export TENCENTCLOUD_SECRET_KEY="YOUR_ACCESS_SECRET"
  export TENCENTCLOUD_REGION="REGION"
  ```

5. **region**  
   地域可选值参考[地域可选值](https://cloud.tencent.com/document/api/248/30346#.E5.9C.B0.E5.9F.9F.E5.88.97.E8.A1.A8)


### 3. 启动 Exporter

```bash
> ksc_exporter --config.file "exporter_config.yml"
```

访问 [http://127.0.0.1:9123/metrics](http://127.0.0.1:9123/metrics) 查看所有导出的指标


## 4、支持的命令行参数说明

命令行参数|说明|默认值
-------|----|-----
--web.listen-address|http服务的端口|9123
--web.telemetry-path|http访问的路径|/metrics
--web.enable-exporter-metrics|是否开启服务自身的指标导出, promhttp_\*, process_\*, go_*|false
--web.max-requests|最大同时抓取/metrics并发数, 0=disable|0
--config.file|产品实例指标配置文件位置|qcloud.yml
--log.level|日志级别|info