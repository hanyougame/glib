# glib
k1 game pkg

### 安装
```shell
go get github.com/hanyougame/glib
```

### 目录结构
- [`gamepkg`](./gamepkg)
    - [`errors`](./errors) 错误处理
    - [`interceptor`](./interceptor) 拦截器
    - [`metadata`](./metadata) 上下文
    - [`stores`](./stores) 数据存储
        - [`cachex`](./stores/cachex) 缓存封装
        - [`gormx`](./stores/gormx) gorm封装
        - [`redisx`](./stores/redisx) redis封装
    - [`tracing`](./tracing) 链路追踪
    - [`utils`](./utils) 工具包
        - [`httpc`](./utils/httpc) http请求
        - [`uniqued`](./utils/uniqued) 生成唯一id
