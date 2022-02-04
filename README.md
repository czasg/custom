# custom

自定义后端响应服务

## docker 安装并运行
```shell script
docker pull czasg/custom
docker run --rm -it -e PORT=8080 -p 8080:8080 czasg/custom
```

## 自定义参数
1、自定义code
```shell script
curl -v localhost:8080/custom?code=204
```

2、自定义header
```shell script
curl -v localhost:8080/custom?header=cookie:custom&header=content-type:application/json
```

3、自定义body
```shell script
curl -v localhost:8080/custom?body=你好
```

4、常用type
```shell script
curl -v localhost:8080/custom?body={"route":"/custom","method":"GET"}&type=json
```

```shell script
curl -o custom.zip localhost:8080/custom?body=test&type=zip
```
