# LiveTV
Use Youtube live as IPTV feeds

## Install from source code

Install `youtube-dl` by your package manager,we need it.

Install Go 1.13 and above compilation environment.

Enter the following command to install.

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

## Usage

Change `ProxyStream` in `config.ini` to `false` if you not live in China.Will more faster to watch live.

### Direct play Youtube live

`baseURL + "/live.m3u8?url=" + Youtube link`

Default baseURL is `http://127.0.0.1:10088`, and you want to watch `https://www.youtube.com/watch?v=dp8PhLsUcFE`

You can use `http://127.0.0.1:10088/live.m3u8?url=https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3Ddp8PhLsUcFE` to watch this live in Kodi or IPTV.

Because of `youtube-dl` parsing is slow,first time to watch will slower.If live started,it speed only depend on you bandwidth.

### Use in Kodi or TV

Modify `channel.txt` , add Youtube live link in it. First line is channel name,must start with `#`.Second line is Youtube live link.
If you have more, repeate it.

Add `http://127.0.0.1:10088/lives.m3u`(default setting) or `{baseURL}/lives.m3u`(replace baseURL by your setting) in Kodi,and enjoy watch live in Kodi.