# go-cli-docker-image-tags

Docker イメージのタグを取得する CLI ツール

# 使い方

### バイナリの生成

```bash
$ make build
$ ls -l ./.bin/docker-image-tags
```

### イメージのタグを取得

```bash
$ ./.bin/docker-image-tags <namespace> <repository>
1.16.3
1.16.2
```
