# local-ehentai

Local E-Hentai Gallery Index (with ~830k galleries' metadata)

[中文简介](README-zh.md)

# Features

- Search multiple tags, in like mode / exact mode
- Local thumbnails cache
- Export search results as csv / json
- Expunged / min-rating filter
- (Filename / URL based) local files filter, with auto discover

# Build

## Requirement

- Go >= 1.12

## Build

Build package `cmd/local-ehentai`:

```bash
go build ./cmd/local-ehentai
```

Set your workdir to `web` when running

# Usage

## Quick Start

1. Download gdata.json from [Mega](https://mega.nz/#F!oh1U0SIA!WBUcf3PaOvrfIF238fnbTg), shared by Sachia Lanlus.
2. Download & unzip local-ehentai.win.zip ([Releases](https://github.com/firefoxchan/local-ehentai/releases))
3. Put gdata.json in the same directory with local-ehentai.exe
4. Run local-ehentai.exe
5. Open [http://127.0.0.1:8080/](http://127.0.0.1:8080/)

## Search

### Syntax 

Comma-separated tag:value pair  

`tag1:exact value1$, tag2:like value2 [, tag3:value3, value4, ...]`

If the tag is omitted, all tags will be searched.  
If the value ends with $, an exact match will be performed.

### Valid Tags

- `category`
- `uploader`
- `parody`
- `character`
- `artist`
- `group`
- `female`
- `male`
- `language`
- `misc`
- `gid` (exact match)

### Example

`artist:toyo$, female: swim`  
Will match  
(`artist` is `toyo`) AND (`female` contains `swim`)

## Local Files Filter

### Filename Based

#### Auto Discover 

1. Put your archive files under `files` directory (or sub directory)
2. Restart local-ehentai.exe

#### Manually Link File with Gallery

1. Modify file `filesMap.txt` in the same directory with local-ehentai.exe
2. Syntax: `filename/gallery_id`, eg: `abcd.zip/1234`
3. Restart local-ehentai.exe

### URL Based

1. Modify file `existUrls.txt` in the same directory with local-ehentai.exe
2. Restart local-ehentai.exe

## Enable Local Thumbnails Cache

If you have downloaded thumbnails form [this torrent](https://sukebei.nyaa.si/view/2770267), you can use local thumbs cache without connecting to eght.org

1. Unzip thumbs_raw.7z, rename the directory to thumbs
2. Put thumbs directory in the same directory with local-ehentai.exe, then restart it

# Demo

![Galleries](/assets/galleries_demo_v0.0.6.png)
