# local-ehentai

Local E-Hentai Gallery Searcher

### 下载 

[Releases](https://github.com/firefoxchan/local-ehentai/releases)

### 数据文件

gdata.json: [Mega](https://mega.nz/#F!oh1U0SIA!WBUcf3PaOvrfIF238fnbTg)  
thumbs: [Nyaa](https://sukebei.nyaa.si/view/2770267)

### 使用方法

下载解压release中local-ehentai.win.zip后, 将上面的gdata.json放在与exe文件同一目录下, 执行local-ehentai.exe即可  
待载入完成后访问浏览器http://127.0.0.1:8080/使用

### 搜索规则

搜索语法  
`tag1:exact value1$, tag2:like value2 [, tag3:value3 ...]`

例  
`artist:toyo$, female: swim`  
会匹配  
`artist`是`toyo` 并且 `female`中包含`swim`的数据

### Demo

![Galleries](/assets/galleries_demo.png)
