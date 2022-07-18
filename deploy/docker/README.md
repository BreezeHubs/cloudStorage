#### 1、更新系统里的所有的能更新的软件
```
sudo apt-get update
```

#### 2、安装几个工具软件 
```
sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
```

#### 3、增加一个docker的官方GPG key  
gpgkey：是用来验证软件的真伪 ——防伪的
```
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
```

#### 4、下载仓库文件
```
echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

#### 5、再次更新系统
```
sudo apt-get update
```

#### 6、安装docker-ce软件
```
sudo apt-get install docker-ce docker-ce-cli containerd.io -y
```