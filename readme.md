```go
.
├── configs    //配置文件目录
│   ├── certs
│   │   ├── moss.crt
│   │   └── moss.key
│   ├── moss.sql
│   └── moss.yaml
├── go.mod
├── go.sum
├── internal 
│   ├── pkg
│   │   ├── code  // 错误码
│   │   │   ├── base.go
│   │   │   ├── code.go
│   │   │   ├── code_generated.go
│   │   │   └── userservice.go
│   │   ├── logger  // 日志包
│   │   │   ├── doc.go
│   │   │   ├── logger.go
│   │   │   ├── sql.go
│   │   │   └── sql_test.go
│   │   ├── middleware //通用中间件
│   │   │   ├── auth //认证相关中间件
│   │   │   │   ├── auto.go
│   │   │   │   ├── basic.go
│   │   │   │   └── jwt.go
│   │   │   ├── auth.go
│   │   │   ├── context.go
│   │   │   ├── cors.go // 跨域
│   │   │   ├── limit.go // 令牌桶
│   │   │   ├── logger.go 
│   │   │   ├── middleware.go 
│   │   │   ├── requestsid.go //请求ID
│   │   │   └── user_validation.go // 用户校验
│   │   ├── options // 应用通用配置选项
│   │   │   ├── feature.go
│   │   │   ├── grpc.go
│   │   │   ├── insecure_serving.go
│   │   │   ├── jwt_options.go
│   │   │   ├── mysql_options.go
│   │   │   ├── redis_options.go
│   │   │   ├── secure_serving.go
│   │   │   └── server_run_options.go
│   │   └── server // 原生http服务
│   │       ├── config.go
│   │       └── genericapiserver.go
│   └── userservice //微服务案例,用户中心服务
│       ├── app.go // 创建App实例
│       ├── auth.go // jwt认证
│       ├── cmd // 程序主入口
│       │   └── userservice.go
│       ├── config
│       │   ├── config.go
│       │   └── doc.go
│       ├── model // 模型层
│       │   └── v1
│       │       ├── users.go
│       │       └── validate.go
│       ├── controller // 控制层
│       │   └── v1
│       │       ├── cache
│       │       │   └── cache.go
│       │       └── user
│       │           ├── create.go
│       │           ├── get.go
│       │           └── user.go
│       ├── grpc.go
│       ├── options
│       │   ├── options.go
│       │   └── validation.go
│       ├── proto
│       │   └── v1
│       │       ├── cache.pb.go
│       │       └── cache.proto
│       ├── route.go 
│       ├── run.go
│       ├── server.go
│       ├── service //业务层
│       │   └── v1
│       │       ├── services.go
│       │       └── users.go
│       └── store //仓库层
│           ├── mysql
│           │   ├── mysql.go
│           │   └── user.go
│           ├── store.go
│           └── user.go
├── pkg
│   ├── app // 通用APP模板
│   │   ├── app.go
│   │   ├── cmd.go
│   │   ├── config.go
│   │   ├── doc.go
│   │   ├── help.go
│   │   └── options.go
│   ├── db // 数据库连接
│   │   └── mysql.go
│   ├── log // 日志包
│   │   ├── LICENSE
│   │   ├── README.md
│   │   ├── context.go
│   │   ├── cronlog
│   │   │   ├── doc.go
│   │   │   └── log.go
│   │   ├── distribution
│   │   │   ├── doc.go
│   │   │   └── logger.go
│   │   ├── doc.go
│   │   ├── encoder.go
│   │   ├── example
│   │   │   ├── context
│   │   │   │   ├── doc.go
│   │   │   │   └── main.go
│   │   │   ├── doc.go
│   │   │   ├── example.go
│   │   │   ├── simple
│   │   │   │   ├── doc.go
│   │   │   │   └── simple.go
│   │   │   └── vlevel
│   │   │       ├── doc.go
│   │   │       └── v_level.go
│   │   ├── go.sum
│   │   ├── klog
│   │   │   ├── doc.go
│   │   │   └── logger.go
│   │   ├── log.go
│   │   ├── log_test.go
│   │   ├── logrus
│   │   │   ├── doc.go
│   │   │   ├── hook.go
│   │   │   └── logger.go
│   │   ├── options.go
│   │   ├── options_test.go
│   │   └── types.go
│   ├── shutdown //优雅关闭
│   │   ├── LICENCE
│   │   ├── README.md
│   │   ├── doc.go
│   │   ├── shutdown.go
│   │   ├── shutdown_test.go
│   │   └── shutdownmanagers
│   │       └── posixsignal
│   │           ├── doc.go
│   │           ├── posixsignal.go
│   │           └── posixsignal_test.go
│   ├── storage // redis
│   │   ├── redis_cluster.go
│   │   └── storage.go
│   ├── util //通用工具包
│   │   └── gormutil
│   │       └── gorm.go
│   └── validator
│       ├── README.md
│       ├── doc.go
│       ├── error.go
│       ├── example
│       │   └── example.go
│       ├── options.go
│       ├── types.go
│       ├── validation.go
│       └── validator.go
├── readme.md
└── third_party // 第三方包
    └── forked
        └── murmur3
            ├── LICENSE
            ├── README.md
            ├── murmur.go
            ├── murmur128.go
            ├── murmur32.go
            ├── murmur32_legacy.go
            ├── murmur64.go
            └── murmur_test.go


```