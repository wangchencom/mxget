# mxget

`mxget` 是一款用Go语言编写的命令行程序，是 [music-get](https://github.com/winterssy/music-get) 的升级版，开发的初衷是为程序员提供更优雅的音乐下载体验。

[![Actions Status](https://github.com/winterssy/mxget/workflows/Build/badge.svg)](https://github.com/winterssy/mxget/actions)
[![Actions Status](https://github.com/winterssy/mxget/workflows/Publish%20Docker/badge.svg)](https://github.com/winterssy/mxget/actions)

## 功能特性

- [网易云音乐](https://music.163.com) / [QQ音乐](https://y.qq.com) / [咪咕音乐](http://music.migu.cn/v3) / [酷狗音乐](http://www.kugou.com) / [酷我音乐](http://www.kuwo.cn/) 一站式音乐搜索和下载。
- 单曲、专辑、歌单以及歌手热门歌曲，只需一步，就能搞定！
- 支持自动嵌入音乐标签/下载歌词。
- 利用goroutines的先天优势实现快速批量下载。
- 支持库调用和RESTful API。

## 下载安装

```
go get -u github.com/winterssy/mxget
```

## 使用说明

> `mxget` 并不是为破解音乐平台的数字版权限制而生的，仅提供试听版音质下载，如果你喜欢高音质/无损资源，请支持正版。

**本项目不提供可执行程序下载，也不提供任何线上demo演示，如须使用请自行编译。**

### 作为CLI使用

这是 `mxget` 的基础功能，你可以通过终端调用 `mxget` 实现音乐搜索、下载功能。以网易云音乐为例，

- 搜索歌曲

```sh
$ mxget search --from nc -k Faded
```

- 下载歌曲

```sh
$ mxget song --from nc --id 36990266
```

- 下载专辑

```sh
$ mxget album --from nc --id 3406843
```

- 下载歌单

```sh
$ mxget playlist --from nc --id 156934569
```

- 下载歌手热门歌曲

```sh
$ mxget artist --from nc --id 1045123
```

如果你希望 `mxget` 为你自动更新音乐标签，可使用 `--tag` 指令，如：

```sh
$ mxget song --from nc --id 36990266 --tag
```

当使用 `--tag` 指令时，`mxget` 会同时将歌词内嵌到音乐文件中，一般而言你无须再额外下载歌词。如果你确实需要 `.lrc` 格式的歌词文件，可使用 `--lyric` 指令，如：

```sh
$ mxget song --from nc --id 36990266 --lyric
```

- 设置默认下载目录

默认情况下，`mxget` 会下载音乐到当前目录下的 `downloads` 文件夹，如果你想要更改此行为，可以这样做：

```sh
$ mxget config --cwd [directory]
```

>  **注：** `directory` 必须为绝对路径。

- 设置默认音乐平台

`mxget` 允许你设置默认使用的音乐平台，如：

```sh
$ mxget config --from qq
```

这样，如果你不通过 `--from` 指令指定音乐平台，`mxget` 便会使用默认值。

在上述命令中，你会经常用到 `--from` 以及 `--id` 这两个指令，它们分别表示音乐平台标识和音乐id。`mxget` 使用的平台标识如下：

|  音乐平台  |     音乐标识     | 识别码 |
| :--------: | :--------------: | :----: |
| 网易云音乐 | `netease` / `nc` |  1000  |
|   QQ音乐   |       `qq`       |  1001  |
|  咪咕音乐  |  `migu` / `mg`   |  1002  |
|  酷狗音乐  |  `kugou` / `kg`  |  1003  |
|  酷我音乐  |  `kuwo` / `kw`   |  1004  |

音乐id为各平台为对应资源分配的唯一id，当使用 `mxget` 进行搜索时，歌曲id会显示在每条结果的后面，你也可以通过网页搜索相关资源，然后从其URL中获取音乐id。值得注意的是，酷狗音乐对应的歌曲id即为文件哈希 `hash` 。

- 多任务并发下载

`mxget` 支持快速批量下载，你可以通过 `--limit` 参数指定同时下载的任务数，最大32。如不指定默认为CPU核心数。

```sh
$ mxget playlist --from nc --id 156934569 --limit 16
```

### 作为库调用

`mxget` 封装了一些便捷的API，Go开发者可以直接调用，举个例子：

```go
package main

import (
	"fmt"
	"github.com/winterssy/mxget/pkg/provider/netease"
)

func main() {
	resp, err := netease.GetSong("36990266")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
```

Tip：使用前你须对作者开发的另一款网络请求库 [sreq](https://github.com/winterssy/sreq) 有所了解，更多细节请阅读 `mxget` 的源码。

> 网易云音乐API的加解密算法参考 [Binaryify/NeteaseCloudMusicApi](https://github.com/Binaryify/NeteaseCloudMusicApi) 并用Golang实现，但 `mxget` 并未移植原项目的所有API，如开发者需要，可fork本项目实现，很简单。

### 作为API服务部署

`mxget` 提供了简易的RESTful API，允许你基于 `mxget` 开发web应用。启动服务：

```sh
$ mxget serve
```

Docker版：

```sh
$ docker pull winterssy/mxget
$ docker run -d --name mxget -p 8080:8080 winterssy/mxget
```

请求方法均为 `GET` ，示例：

- 从QQ音乐获取 `周杰伦` 的搜索结果

```sh
$ curl -X GET "http://127.0.0.1:8080/api/qq/search/周杰伦" -H "accept: application/json"
```

- 从网易云音乐获取id为 `36990266` 的歌曲资源

```sh
$ curl -X GET "http://127.0.0.1:8080/api/netease/song/36990266" -H "accept: application/json"
```

- 从咪咕音乐获取id为 `1121438701` 的专辑资源

```sh
$ curl -X GET "http://127.0.0.1:8080/api/migu/album/1121438701" -H "accept: application/json"
```

- 从酷狗音乐获取id为 `547134` 的歌单资源

```sh
$ curl -X GET "http://127.0.0.1:8080/api/kugou/playlist/547134" -H "accept: application/json"
```

- 从酷我音乐获取id为 `336` 的歌手资源

```sh
$ curl -X GET "http://127.0.0.1:8080/api/kuwo/artist/336" -H "accept: application/json"
```

**注：** 由于音乐平台的限制，`mxget` 的API服务仅在本地测试通过。如果你将 `mxget` 部署到公网，特别是海外VPS上，开发者不保证能工作，遇到的问题需要你自行解决。

## 免责声明

- 本项目仅供学习研究使用，禁止商业用途。
- 本项目使用的接口如无特别说明均为官方接口，音乐版权归源音乐平台所有，侵删。

## License

GPLv3。