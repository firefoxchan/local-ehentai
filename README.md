# local-ehentai

Local E-Hentai Gallery Searcher

### Download

[Releases](https://github.com/firefoxchan/local-ehentai/releases)

### Usage

```bash
local-ehentai -j gdata.json -h 127.0.0.1:8080 -t thumbs
```

### Gallery Data Files

gdata.json: [Mega](https://mega.nz/#F!oh1U0SIA!WBUcf3PaOvrfIF238fnbTg)  
thumbs: [Nyaa](https://sukebei.nyaa.si/view/2770267)

### Search Rule

Syntax:  
`tag1:exact value1$, tag2:like value2 [, tag3:value3 ...]`

Example:  
`artist:toyo$, female: swim`  
Will match  
`artist is toyo AND female contains swim`

### Demo

![Galleries](/assets/galleries_demo.png)
