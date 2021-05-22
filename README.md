# web-screenshot

网页截屏工具，基于chromedep实现对访问网址并实现截屏功能。

chromedep需要依赖于chrome，在服务器可直接使用`chromedp/headless-shell`来部署，执行脚本如下：

```bash
docker run -d -p 9222:9222 --restart=always \
  --name headless-shell \
  chromedp/headless-shell
```