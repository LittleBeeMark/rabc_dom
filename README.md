# rabc_dom 
## 这是一个简单的，使用 casbin 的 rabc-dom 就是多租户模型或叫角色域模型实现的一个权限控制示例
具体描述可以看一下我博客的这片[文章](https://marksuper.xyz/2021/06/18/casbin_rabc_dom/)

简单描述一下 rabc-dom 模型，简单的权限控制一般是用户，角色，访问资源这三者进行排列组合形成权限。
举个简单的例子：用户A为管理员角色可以访问订单模块的增删改查方法，用户B为审计员只能访问订单的查方法。
使用 Casbin 的 rabc 模型就可以把权限控制抽离出来不用耦合在不同的 server 层代码中

rabc-dom 加入一个角色域的概念。相当于项目。用户，角色，项目，访问资源，四者进行排列组合形成权限。
同样举个简单的例子，用户A在X项目中是管理员，在Y项目中是审计员，那么他能访问X中的订单增删改查方法，但只能访问Y项目中的订单的查方法。

这个例子就实现了使用 Casbin 的 rabc-dom 模型进行用户，角色，项目，访问资源这四者的权限抽离形成中间键

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
