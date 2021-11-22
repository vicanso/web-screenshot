# web-screenshot

网页截屏工具，基于chromedep实现对访问网址并实现截屏功能。

chromedep需要依赖于chrome，在服务器可直接使用`chromedp/headless-shell`来部署，命令如下：

```bash
docker run -d -p 9222:9222 --restart=always \
  --name headless-shell \
  --shm-size=256m \
  chromedp/headless-shell
```

上面版本的镜像中并没有提供中文字库的支持，如果需要支持中文，可以参考`Dockerfile_headless_chinese`的生成自定义镜像，或者可直接使用`vicanso/headless-shell`。docker指定更大的shm空间`--shm-size=256m`

```bash
docker build -f Dockerfile_headless_chinese -t vicanso/headless-shell .
```


启动web-screenshot，命令如下：

```bash
docker run -d -p 7000:7000 \
  -e DEV_TOOLS_WS_URL=ws://172.16.23.175:9222/ \
  --name web-screenshot \
  vicanso/web-screenshot
```

其中`DEV_TOOLS_WS_URL`从`headless-shell`的启动日志中获取，之后则可通过访问地址：`http://127.0.0.1:7000/capture-screenshot?url=https%3A%2F%2Fwww.baidu.com%2F&width=800&height=600&overrideHeader=1`来生成对应地址的截屏。

- `url`: 必填参数，需要截屏的地址
- `width`: 可选参数，浏览器的窗口宽度
- `height`: 可选参数，浏览器的窗口高度
- `overrideHeader`: 可选参数，是否覆盖请求的header，如果设置此参数，则以访问`/capture-screenshot`的请求头替换headless的默认请求头
- `dalay`: 可选参数，延时等待后截屏，如delay=5s
- `visible`: 可选参数，等待元素可见，如visible=#load，使用css选择器
