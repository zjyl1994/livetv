# LiveTV
将 Youtube 直播作为 IPTV 电视源

## 安装方法

### 从源码安装

本程序依赖 `youtube-dl` ，请自行通过包管理器等方式进行安装。

```bash
git clone https://github.com/zjyl1994/livetv.git
cd livetv
go build
sudo cp livetv /usr/local/bin/
sudo mkdir -p /etc/livetv
sudo cp config.ini playlist.m3u channel.txt /etc/livetv/
sudo cp livetv.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start livetv
```
## 配置

### config.ini

**如果你不清楚你在做什么请不要修改本文件。**

BaseURL 外部访问的URL地址

ListenOn  程序监听的本地URL地址

ProxyStream 是否代理直播流，默认true

M3UTemplate m3u文件模板位置

ChannelFile 频道定义文件位置

YtdlCmd youtube-dl调用命令

YtdlArgs youtube-dl调用参数

### playlist.m3u

如果你需要修改生成的m3u节目清单，可以在尾部追加任意符合格式的条目，但请不要删除已有内容。

### channel.txt

频道定义文件，两行一个频道。第一行一般为频道名，以`#`号开头。第二行为Youtube直播地址，无需做任何转换。

## 使用

假设当前 `baseURL` 设置为：`http://192.168.31.100:10088`

### 直接播放某个油管直播

访问 `http://192.168.31.100:10088/live.m3u8?url=` + Youtube地址。注意，此处Youtube地址需要做Urlencode转码。

比如想观看 三立新闻台的直播 [https://www.youtube.com/watch?v=4ZVUmEUFwaY](https://www.youtube.com/watch?v=4ZVUmEUFwaY)。
那么，可以使用 `http://192.168.31.100:10088/live.m3u8?url=https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3D4ZVUmEUFwaY` 直接播放对应的直播。

由于后台使用了`youtube-dl`进行视频解析，所以画面加载速度不会特别快，但是后续的播放会很流畅，不用担心。

### 在电视APP或者Kodi上播放

修改 `channel.txt` ，将你的直播链接直接写在里面。具体格式为：

```text
#直播的显示名称
直播地址（无需任何转换直接贴进来）
```

可以参考代码中带的`channel.txt`文件，设置多个直播源。#号标注的直播名称会作为IPTV的频道显示。

在Kodi或电视上操作IPTV客户端，设置播放清单为 `http://192.168.31.100:10088/lives.m3u` ，重启Kodi就可以看到这些频道了。

## 代理模式
本程序支持代理模式，当`config.ini`中的`ProxyStream`设置为true时， 并且在`livetv.service`文件中设定`HTTP_PROXY` 环境变量后。所有直播流量会被本程序代理。
这种使用场景可以解决电视上没法科学上网的困扰，只需要在路由器或者家庭服务器里运行本程序，使用给定的M3U文件，电视就可以以IPTV形式直接观看油管直播了。