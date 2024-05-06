# Docker обёртка Scorecard

```shell
docker build -t scorecard .
docker run -e GITHUB_AUTH_TOKEN=аuth_token -e REPO_URL=url -e Checks=Checks -it scorecard
```
Checks перечисляются через запятую. Например, CI-Tests,Binary-Artifacts
