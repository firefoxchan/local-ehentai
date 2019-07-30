# local-ehentai

Local E-Hentai Gallery Searcher

### 下载 

[Releases](https://github.com/firefoxchan/local-ehentai/releases)

### 数据文件

gdata.json: [Mega](https://mega.nz/#F!oh1U0SIA!WBUcf3PaOvrfIF238fnbTg)  
thumbs: [Nyaa](https://sukebei.nyaa.si/view/2770267)

# 使用方法

## 启动服务

1. 下载 Sachia Lanlus 提供的 gdata.json : [Mega](https://mega.nz/#F!oh1U0SIA!WBUcf3PaOvrfIF238fnbTg) 
2. 下载并解压 [Releases](https://github.com/firefoxchan/local-ehentai/releases) 中的 local-ehentai.win.zip
3. 将 gdata.json 放到 local-ehentai.exe 同级目录
4. 运行 local-ehentai.exe
5. 访问 [http://127.0.0.1:8080](http://127.0.0.1:8080/)

## 搜索

### 语法  

一组逗号分隔的tag:value

`tag1:exact value1$, tag2:like value2 [, tag3:value3, value4, ...]`

tag留空时会搜索所有tag, value以$结尾时会精确匹配

### Tag列表

- category
- uploader
- parody
- character
- artist
- group
- female
- male
- language
- misc

### 样例 

`artist:toyo$, female: swim`  
会匹配  
`artist`是`toyo` 并且 `female`中包含`swim`的数据

## 启用本地缩略图缓存

如果你从 [这里](https://sukebei.nyaa.si/view/2770267) 下载了缩略图包, 可以直接开启本地缩略图缓存, 不用访问eght.org

1. 解压 thumbs_raw.7z, 重命名文件夹为 thumbs
2. 把 thumbs 文件夹放到与 local-ehentai.exe 同级目录下, 然后重启 local-ehentai.exe

# 预览

![Galleries](/assets/galleries_demo.png)
