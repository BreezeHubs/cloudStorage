```shell
docker run -d --hostname rabbit-svr --name rabbit -p 5672:5672 -p 15672:15672 -p 25672:25672 -v E:/Desktop/项目/cloudStorage/deploy/rabbitmq/data/rabbitmq:/var/lib/rabbitmq rabbitmq:management
```