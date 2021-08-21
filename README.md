# previewQiitaGo

## What's this ??

Qiita の記事の Markdown を取得し、HTML に変換するサーバ



## How to use

1. サーバ起動

   ```go
   go run main.go
   ```

2. 起動したらブラウザで、`http://localhost:8080?target={{ Qiita の記事の URL }}` にアクセスすると変換したものが見えます



## For what ??

これを GAE などにデプロイしておけば、プロキシの関係などで PC から直接 Qiita にアクセスできなくても Qiita の記事を読むことができます