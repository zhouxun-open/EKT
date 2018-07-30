# playground

基于docker 快速搭建N 个测试节点

## 运行环境
需要提前安装ecli，当前脚本依赖ecli 生成配置信息

python 依赖安装如下:
```
pip install -r requirement.txt
```

## deploy.py 帮助脚本

目前支持可以自定义配置如下
- `EKTCLI`: 指定ecli 执行路径
- `HOST_ADDR`: 指定一组用于部署节点的机器（需要保证已安装并启动docker）
- `NODE_NUM`: 每台机器上安装节点的个数（节点会安装在容器中做隔离）
- `PORT_RANGE`: 每台机器上节点可以使用的端口号，个数需要与`NODE_NUM` 一致
配置完成后即可执行以下几部进行部署


## 1.生成配置

```
python deploy.py gen_conf
gen conf done.
```

## 2.部署

```
python deploy.py deploy
publish conf done.
latest: Pulling from registry.cloudhua.com/ekt8/ekt8
......
```

## 3.查看

```
python deploy.py run_cmd 'docker ps'
CONTAINER ID        IMAGE                                    COMMAND                CREATED             STATUS              PORTS                      NAMES
88c4e349d2eb        registry.cloudhua.com/ekt8/ekt8:latest   "/bin/sh entrypoint.   6 minutes ago       Up 6 minutes        0.0.0.0:19990->19990/tcp   192.168.6.54_19990   
6ef5a1fe5175        registry.cloudhua.com/ekt8/ekt8:latest   "/bin/sh entrypoint.   6 minutes ago       Up 6 minutes        0.0.0.0:19991->19991/tcp   192.168.6.54_19991   
8919c5c70883        registry.cloudhua.com/ekt8/ekt8:latest   "/bin/sh entrypoint.   6 minutes ago       Up 6 minutes        0.0.0.0:19992->19992/tcp   192.168.6.54_19992   
CONTAINER ID        IMAGE                                    COMMAND                CREATED             STATUS              PORTS                      NAMES
b76d9dea5fde        registry.cloudhua.com/ekt8/ekt8:latest   "/bin/sh entrypoint.   6 minutes ago       Up 6 minutes        0.0.0.0:19991->19991/tcp   192.168.6.55_19991   
45c23bedcf08        registry.cloudhua.com/ekt8/ekt8:latest   "/bin/sh entrypoint.   6 minutes ago       Up 6 minutes        0.0.0.0:19992->19992/tcp   192.168.6.55_19992   
9f0969434193        registry.cloudhua.com/ekt8/ekt8:latest   "/bin/sh entrypoint.   6 minutes ago       Up 6 minutes        0.0.0.0:19990->19990/tcp   192.168.6.55_19990
......
```

## 清理容器
- 注意数据无法恢复，一般用于重建测试网络

```
python deploy.py clean
192.168.6.54_19992
192.168.6.54_19991
192.168.6.54_19990
192.168.6.55_19990
192.168.6.55_19992
...
```

## 注意
- 宿主机的iptables 可能会阻止容器访问网络，可以关闭`iptables` 或增加以下规则，`docker0` 是docker 的默认网卡(按需配置)

```
iptables -A INPUT -i docker0 -j ACCEPT
```