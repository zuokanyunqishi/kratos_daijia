### docker-compose 初始化环境
   [环境初始化](docker-compose-deploy/docker-compose.yaml)
   

### 环境初始化化成功后
    在 consul 中插入 高德路径规划接口的 apikey 
   [申请地址](https://console.amap.com/dev/key/app)


```Yaml
## path map/amap.yaml
 amap:
   direction:
      key: xxx
```