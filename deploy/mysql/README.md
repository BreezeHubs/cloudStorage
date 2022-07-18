主从模式  
  
docker创建主从模式的数据库实例  
需要一个主节点、一个从节点
https://www.freesion.com/article/3931408189/
```
docker run --name mysql-master -p 3306:3306 -e MYSQL_ROOT_PASSWORD=12345678 -d mysql:latest
docker run --name mysql-slave -p 3307:3306 -e MYSQL_ROOT_PASSWORD=12345678 -d mysql:latest
```
![图片(./res/20220715235445.png)


主节点查看状态
```
show master status; #可以看到binlog文件的位置
```

从节点需要指定主节点的IP地址和端口号
```
CHANGE MASTER TO MASTER_HOST = ' 主节点的IP地址', MASTER_PORT = 3306, MASTER_USER = 'root', MASTER_PASSWORD = '密码', MASTER_LOG_FILE = 'mysql-bin.000001', MASTER_LOG_POS = 0;

start slave; # 启动从节点
show slave status\G; # 查看从节点的状态
```