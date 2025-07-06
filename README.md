##使用方法
1. 编译前端
```
npm run build
```
2. 复制文件
```
rm -rf dev_tool/dist/
cp -r web_source/dist dev_tool/
```
3. 编译服务
```
cd dev_tool
go build
```

## 开发
```
cd dev_tool
go run app.go
cd web_source
npm run dev
```

## Windows编译
```
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
```
windows运行需要 config目录和data目录