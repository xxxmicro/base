


### Microbase
Microbase 是 lucfish 开源的一套 GO 微服务框架，集成了目前微服务架构的大量最佳实践策略.

Goals
我们致力于提供完整的微服务研发体验，整合相关框架及工具后，微服务治理相关部分可对整体业务开发周期无感，从而更加聚焦于业务交付。对每位开发者而言，整套 Microbase 框架也是不错的学习仓库，可以了解和参考到 lucfish 在微服务方面的技术积累和经验。

Features

* http: 核心基于 gin 进行模块化设计
* Ioc: 应用的模块依赖组装通过 uber的依赖注入框架fx完成，目前另一个对标是 facebook 的 facebookgo/inject, 大相径庭，只是做了二选一策略
* config: 基于 gomicro 实现，source 部分实现了 远程配置中心 apollo 的扩展,
* broker: 基于 gomicro 实现，实现了 rocketmq 的扩展
* cache: 优雅接口设计，充分参考了 kratos(github.com/go-kratos/kratos), gomicro 等
* database: 集成 gorm, 添加 熔断保护 和 统计支持，可快速发现数据压力
* log: 依赖 zap 实现高性能日志库，并结合 log-agent实现远程日志管理
* trace: 基于 opentracing, 集成了全链路trace支持 (gPRC/HTTP/MySQL/Redis), 可切换zipkin/jaeger
* generator: 工具链，可快速生成标准项目

