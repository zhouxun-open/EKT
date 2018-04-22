# EKT文档[](https://github.com/EducationEKT/EKT/tree/master/docs)
   * [EKT白皮书](docs/whitepaper.md)
   * [EKT路线图](docs/roadmap.md)

# 百万代币空投技术社区!
官方技术QQ群:699726921


# 部署

1. 安装Golang环境和gcc环境
    1.1 安装gcc环境
        * centos
		 `yum -y install gcc`
        * ubuntu
		 `sudo apt-get install gcc`

        大家也可以自己下载源码进行安装。安装完成之后可以使用`gcc -v`查看是否安装成功。

    1.2 安装Golang环境
	从golang官网下载安装包，解压到/usr/local目录下，解压后的可以加上版本号，软连接成/usr/local/go目录，方便以后更新golang版本。
	修改/etc/profile，设置Go语言的环境，在/etc/profile最后增加一下的代码

```
export GOROOT=/usr/local/go
export GOPATH=/opt/gopath
export PATH=$GOROOT/bin:$PATH:$GOPATH/bin
```
	如果/opt目录下没有gopath文件夹的话，可以先新建gopath文件夹 `mkdir /opt/gopath`
	最后让设置生效，执行`source /etc/profile`
	大家可以使用go version查看go语言是否已经安装成功，也可以通过go env判断go的一些其他配置。


2. 从GitHub上下载最新源码
```
    cd $GOPATH
    mkdir -p src/github.com/EducationEKT
    cd src/github.com/EducationEKT
    git clone https://github.com/EducationEKT/EKT
    cd github.com/EducationEKT/EKT
```

3. 修改本地的配置文件genesis.json
```
    vim genesis.json
```
把dbPath、logPath、node和blockchainManagePwd修改成自己的。

4. 启动节点,在测试阶段可以不用打包，直接命令行运行就可以了
```
    mkdir -p /var/log/EKT
    go run io/ekt8/main.go genesis.json 1>/var/log/EKT/stdout 2>/var/log/EKT/stderr &
```

5. 查看stdout或者stderr可以使用 `tail -f /var/log/EKT/stdout`, 如果需要看其他日志，可以cd到genesis.json中配置的日志的目录中进行查看
