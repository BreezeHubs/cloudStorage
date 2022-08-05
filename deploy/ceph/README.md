### 修改docker源为国内源
```json
{
    "registry-mirrors": [
        "https://registry.docker-cn.com",
        "https://docker.mirrors.ustc.edu.cn",
        "http://hub-mirror.c.163.com",
        "https://cr.console.aliyun.com/"
  ]
}
```

### 拉取镜像
```shell
docker pull ceph/mon
docker pull  ceph/osd
docker pull  ceph/radosgw
```

### 创建ceph网桥
```shell
docker network create --driver bridge --subnet 172.20.0.0/16 ceph-network
```

### 创建相关目录及修改权限，用于挂载volume
```shell
mkdir -p /www/ceph /var/lib/ceph/osd /www/osd/
 
chown -R 64045:64045 /var/lib/ceph/osd
 
chown -R 64045:64045 /www/osd/
```

### 创建monitor节点
```shell
docker run -itd --name monnode --network ceph-network --ip 172.20.0.10 -e NON_NAME=monnode -e MON_IP=172.20.0.10 -v E:/Desktop/项目/cloudStorage/deploy/ceph/www/ceph:/etc/ceph ceph/mon
```

### 在monitor节点上标识3个osd节点
```shell
docker exec monnode ceph osd create
 
docker exec monnode ceph osd create
 
docker exec monnode ceph osd create
```

### 创建osd节点
```shell
docker run -itd --name osdnode0 --network ceph-network -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v E:/Desktop/项目/cloudStorage/deploy/ceph/www/ceph:/etc/ceph -v E:/Desktop/项目/cloudStorage/deploy/ceph/www/osd0:/var/lib/ceph/osd/ceph-0 ceph/osd
 
docker run -itd --name osdnode1 --network ceph-network -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v E:/Desktop/项目/cloudStorage/deploy/ceph/www/ceph:/etc/ceph -v E:/Desktop/项目/cloudStorage/deploy/ceph/www/osd1:/var/lib/ceph/osd/ceph-1 ceph/osd
 
docker run -itd --name osdnode2 --network ceph-network -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v E:/Desktop/项目/cloudStorage/deploy/ceph/www/ceph:/etc/ceph -v E:/Desktop/项目/cloudStorage/deploy/ceph/www/osd2:/var/lib/ceph/osd/ceph-2 ceph/osd
```

### 增加monitor节点，组件成机器
```shell
docker run -itd --name monnode_1 --network ceph-network --ip 172.20.0.11 -e NON_NAME=monnode_1 -e MON_IP=172.20.0.11 -v E:/Desktop/项目/cloudStorage/deploy/ceph/www/ceph:/etc/ceph ceph/mon
 
docker run -itd --name monnode_2 --network ceph-network --ip 172.20.0.12 -e NON_NAME=monnode_2 -e MON_IP=172.20.0.12 -v E:/Desktop/项目/cloudStorage/deploy/ceph/www/ceph:/etc/ceph ceph/mon
```

### 创建gateway节点
```shell
docker run -itd --name gwnode --network ceph-network --ip 172.20.0.9 -p 9080:80 -e RGW_NAME=gwnode -v E:/Desktop/项目/cloudStorage/deploy/ceph/www/ceph:/etc/ceph ceph/radosgw
```

### 查看ceph集群状态
```shell
docker exec monnode ceph -s
```
运行结果，可以看到ceph集群状态，health不为HEALTH_OK、osdmap e19: 3 osds: 3 up, 3 in不是全部up，说明集群未成功启动，重启集群后再查看
```
cluster f559585b-8ff7-426d-8f62-acf71484f376
    health HEALTH_OK
    monmap e3: 3 mons at {650e8a0f441a=172.20.0.10:6789/0,65774d73ee8d=172.20.0.12:6789/0,af003949dcad=172.20.0.11:6789/0}
        election epoch 8, quorum 0,1,2 650e8a0f441a,af003949dcad,65774d73ee8d
    osdmap e19: 3 osds: 3 up, 3 in
        flags sortbitwise
    pgmap v39: 104 pgs, 6 pools, 848 bytes data, 81 objects
        119 GB used, 172 GB / 292 GB avail
                104 active+clean
client io 10874 B/s rd, 0 B/s wr, 25 op/s
```

### 创建accessKey、secretKey
```shell
docker exec -it gwnode radosgw-admin user create --uid=user1 --display-name=user1
```
运行结果，可以看到accessKey、secretKey
```
{
    "user_id": "user1",
    "display_name": "user1",
    "email": "",
    "suspended": 0,
    "max_buckets": 1000,
    "auid": 0,
    "subusers": [],
    "keys": [
        {
            "user": "user1",
            "access_key": "4R6VUG9BLUVTDAZIQ231",
            "secret_key": "yRHsmkfKaedtUxfGKJdULNVAfNG9q5mO9lxmb4ZP"
        }
    ],
    "swift_keys": [],
    "caps": [],
    "op_mask": "read, write, delete",
    "default_placement": "",
    "placement_tags": [],
    "bucket_quota": {
        "enabled": false,
        "max_size_kb": -1,
        "max_objects": -1
    },
    "user_quota": {
        "enabled": false,
        "max_size_kb": -1,
        "max_objects": -1
    },
    "temp_url_keys": []
}
```