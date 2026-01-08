# 在公网上部署应用

> 服务器系统：Ubuntu 24.04

## 数据库(pg)

ubuntu系统自带一个pg数据库的指定版本快照，如果想要安装其他版本，参考[官方文档](https://www.postgresql.org/download/linux/ubuntu/)

- `apt install postgresql`

### 创建用户/数据库

> 记得在云服务器厂商的控制台，编辑网络策略：允许访问5432端口（如果你使用了非默认的端口号，这里相应修改）

```cmd 
su - postgres // 切换linux用户
psql

create user [user name] with password '[password]'; // 注意分号、单引号
create database [db name] owner [user name];
grant all on database [db name] to [user name];

\du+ // 查看用户
\l   // 查看数据库
```

### 设置允许外部地址访问

config file: `/etc/postgresql/16/main/postgresql.conf`(注意版本号)
add: `listen_address: '*'`

config file: `/etc/postgresql/16/main/pg_hba.conf`
add: 
```txt
host cloud all 127.0.0.1/32 md5     // 允许本机连接cloud数据库
host cloud all 0.0.0.0/0    reject  // 不允许其他连接访问cloud
host all   all 0.0.0.0/0    md5     // 允许所有远程连接
```

重启/重载服务：`sudo systemctl restart/reload postgresql`

## nginx

### 下载与安装

可以使用包管理工具下载，或参考[nginx官网](https://nginx.org/en/linux_packages.html#Ubuntu)提供的下载与安装过程

- `sudo apt install nginx`
- `sudo systemctl start nginx`
- `curl http://127.0.0.1` // 验证安装结果

### 编辑反向代理规则

配置文件路径：`/etc/nginx/nginx.conf`
日志文件路径：`/var/log/nginx/access.log`（配置文件内可修改）

默认配置包含`/etc/nginx/conf.d/*.conf`

详细配置，见本文档同目录下，`*.conf`配置文件

- 如果想要使用https访问，可以使用自签证书（生成方式见下一节）

nginx常用命令：

- `nginx -t`：测试配置文件的语法是否正确
- `nginx -s reload`：重启nginx，更新并应用新的配置文件

### 自签证书

> 目的是允许使用https访问

生成http证书：

- `openssl genrsa -out server.key 2048`
- `openssl req -new -key server.key -out server.csr`
    - 会以问答的形式，要求我们输入一些信息，包括**域名**(Common Name)在内的一些关键信息就是在此处提供的
- `openssl x509 -req -in server.csr -out server.crt -signkey server.key -days 3650`

## 服务器设置

### 设置使用密钥登录

- `ssh-keygen`生成密钥对
    - 默认在`/root/.ssh/`目录生成`id_ed25519`和`id_ed25519.pub`文件
    - 将没有后缀名的私钥下载到本地
- 设置公钥：`cat ./id_ed25519.pub >> ./authorized_keys`(注意切换路径)
    - 将`id_ed25519.pub`写入`authorized_keys`
- 本地验证公钥的有效性：`ssh -i ./id_ed25519 localhost`
    - 修改ssh服务器安全策略：`sudo vim /etc/ssh/sshd_config`
        - 找到以下配置项并修改为对应值
          ```txt
          PermitRootLogin without-password
          PubkeyAuthentication yes
          PasswordAuthentication no // 禁用密码登录
          ```
        - 测试语法是否正确：`sudo ssh -t`
        - 重启并检查ssh服务状态：`sudo systemctl restart/status ssh`

### 设置开机自启动

在`/etc/rc.local`中，编辑想要在开机时执行的代码

```txt
#!/bin/sh -e
([.exe path] &)
exit 0
```
