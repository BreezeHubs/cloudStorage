```
docker run --restart=always -p 6379:6379 --name myredis -v E:/Desktop/项目/cloudStorage/deploy/redis/myredis/myredis.conf:/etc/redis/redis.conf -v /home/redis/myredis/data:/data -d redis redis-server /etc/redis/redis.conf  --appendonly yes  --requirepass 12345678
```