# mflight-api

[SONY のマルチファンクションライト](https://www.sony.co.jp/Products/multifunctional-light/) から、気温 (temperature)、湿度 (humidity)、照度 (illuminance) を取得するための api

# for testing

- https://github.com/google/go-cmp
- https://github.com/golang/mock

## generate mock

```
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE`
```

```sh
$ go generate ./...
```

```sh
$ find . -name 'mock_*.go'
./infrastructure/mflight/mock_mflight/mock_client.go
./infrastructure/cache/mock_cache/mock_cache.go
```

