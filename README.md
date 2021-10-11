# rabc_dom 
数据库使用的是pg
可以自己在本地建一个，别忘了把/config.yaml 中数据库的配置改为自己的

直接使用
go run main.go
搞起

依赖可以自行拉取
go mod tidy
go mod vendor

编译启动后添加用户，添加项目，登录用户，使用获取项目详情验证用户是否权限正确
也可以自己添加一些模块和接口
仿照/config/policy_source.yml 添加接口权限控制
进行接口权限验证
