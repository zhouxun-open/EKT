# EKT文档[](https://github.com/EducationEKT/EKT/tree/master/docs)
   * [EKT白皮书](docs/whitepaper.md)
   * [EKT路线图](docs/roadmap.md)

# 百万代币空投技术社区!
官方技术QQ群:699726921


# 部署

1. 从GitHub上下载最新源码
```
    cd $GOPATH
    mkdir -p src/github.com/EducationEKT
    cd src/github.com/EducationEKT
    git clone https://github.com/EducationEKT/EKT
    cd github.com/EducationEKT/EKT
```

2. 修改本地的配置文件genesis.json
```
    vim genesis.json
```
把dbPath、logPath、node和blockchainManagePwd修改成自己的。

3. 启动节点,在测试阶段可以不用打包，直接命令行运行就可以了
```
    mkdir -p /var/log/EKT
    go run io/ekt8/main.go genesis.json 1>/var/log/EKT/stdout 2>/var/log/EKT/stderr &
```

4. 查看stdout或者stderr可以使用 `tail -f /var/log/EKT/stdout`, 如果需要看其他日志，可以cd到genesis.json中配置的日志的目录中进行查看
