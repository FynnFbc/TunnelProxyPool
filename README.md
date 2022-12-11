# TunnelProxyPool

- 一款无环境依赖开箱即用的免费代理IP池   
- 内置18个免费代理源，均使用内置的简单正则获取  
- 支持webApi获取、删除等代理池内的IP 
- 支持 http，socket5 隧道代理模式，无需手动更换IP，每一次请求IP都不同
- 遇到bug或有好的建议，欢迎提issue
- 欢迎访问我的 **[博客](https://blog.fynn.vip/ "博客")**

## 隧道代理
- 隧道代理是代理IP存在的一种方式。  
- 相对于传统固定代理IP，它的特点是自动地在代理服务器上改变IP，这样每个请求都使用一个不同的IP。

# 代理IP特征
- 这里提供一些代理IP的特征，师傅们可通过特征自己写代理源，api获取的话内置的正则方式就能写  
- 360网络空间测绘_socket5：
```text
protocol:"socks5" AND "Accepted Auth Method: 0x0" AND "connection: close" AND country: "China"  
```
## fofa_http:
```text  
"HTTP/1.1 403 Forbidden Server: nginx/1.12.1" && port="9091"   

port="3128" && title="ERROR: The requested URL could not be retrieved"  

"X-Cache: 'MISS from VideoCacheBox/CE8265A63696DECD7F0D17858B1BDADC37771805'" && "X-Squid-Error: ERR_ACCESS_DENIED 0"  
```
## hunter_http：
```text
header.server="nginx/2.2.200603d"&&web.title="502 Bad Gateway" && ip.port="8085"
```

# 截图
[![zuz6TU.png](https://s1.ax1x.com/2022/11/19/zuz6TU.png)](https://s1.ax1x.com/2022/11/19/zuz6TU.png)
# 使用说明
### 下载 
```
git clone https://github.com/pingc0y/go_proxy_pool.git
```
- 编译（直接使用成品，就无需编译）  
- 以下是在Windows环境下，编译出各平台可执行文件的命令  
```
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o ../ProxyPool-win-64.exe

set CGO_ENABLED=0
set GOOS=windows
set GOARCH=386
go build -ldflags "-s -w"  -o ../ProxyPool-win-86.exe

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o ../ProxyPool-linux-64

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm64
go build -ldflags "-s -w" -o ../ProxyPool-linux-arm64

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=386
go build -ldflags "-s -w" -o ../ProxyPool-linux-86

set CGO_ENABLED=0
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "-s -w" -o ../ProxyPool-macos-64

set CGO_ENABLED=0
set GOOS=darwin
set GOARCH=arm64
go build -ldflags "-s -w" -o ../ProxyPool-macos-arm64

```
### 运行  
- 需要与config.yml在同一目录  
- 注意：抓取代理会进行类型地区等验证会比较缓慢，存活验证会快很多
```
.\goProxyPool.exe
```

### 代理源中有部分需要翻墙才能访问，有条件就设置下config.yml的代理配置
```yml
proxy:
  host: 127.0.0.1
  port: 10809
```
## webAPi说明
- 查看代理池情况
```
http://127.0.0.1:8080/
```

### 获取代理
```
http://127.0.0.1:8080/get?type=HTTP&count=10&anonymity=all
可选参数：
type        代理类型
anonymity   匿名度
country     国家
source      代理源
count       代理数量
获取所有：all
```

### 删除代理 (默认没有开启，需自行修改源代码编译)
```
http://127.0.0.1:8080/delete?ip=127.0.0.1&port=8888
必须传参：
ip      代理ip
port    代理端口
```

### 验证代理
```
http://127.0.0.1:8080/verify
```

### 抓取代理
```
http://127.0.0.1:8080/spider
```
## 代理字段解读
```go
type ProxyIp struct {
    Ip         string //IP地址
    Port       string //代理端口
    Country    string //代理国家
    Province   string //代理省份
    City       string //代理城市
    Isp        string //IP提供商
    Type       string //代理类型
    Anonymity  string //代理匿名度, 透明：显示真实IP, 普匿：显示假的IP, 高匿：无代理IP特征
    Time       string //代理验证
    Speed      string //代理响应速度
    SuccessNum int    //验证请求成功的次数
    RequestNum int    //验证请求的次数
    Source     string //代理源
}
```
## 配置文件
```yaml
# 使用代理去获取代理IP
proxy:
  host: 127.0.0.1
  port: 10809
  
# 代理身份验证
auth:
  # 用户名
  user: abcde
  # 密码
  pass: qwert

# 配置信息
config:
  #监听IP
  ip: 0.0.0.0
  #web监听端口
  port: 8080
  #http隧道代理端口
  httpTunnelPort: 8111
  #socket隧道代理端口
  socketTunnelPort: 8112
  #隧道代理更换时间秒
  tunnelTime: 60
  #可用IP数量小于‘proxyNum’时就去抓取
  proxyNum: 50
  #代理IP验证间隔秒
  verifyTime: 1800
  #抓取/检测状态线程数
  threadNum: 200

```

## 更新说明
2022/11/22  
修复 ip归属地接口更换  
优化 验证代理   

2022/11/19  
新增 socket5代理  
新增 文件导入代理  
新增 显示验证进度  
新增 验证webApi  
修改 扩展导入格式  
优化 代理验证方式  
优化 匿名度改为自动识别  
修复 若干bug  


