# mflight-api

[SONY のマルチファンクションライト](https://www.sony.co.jp/Products/multifunctional-light/) から、気温 (temperature)、湿度 (humidity)、照度 (illuminance) を取得するための api

# docker image

- https://hub.docker.com/repository/docker/otknoy/mflight-api

# dev

## for testing

- https://github.com/google/go-cmp
- https://github.com/golang/mock

### generate mock

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


# Multifunctional Light 解析


なんか POST できる
```
$ curl -v -X POST 'http://192.168.1.5:60001/'
```

設定取得
```
$ curl -s 'http://192.168.1.5:56000/ConfigurationData.xml?x-KEY_MOBILE_ID=foo&x-KEY_UPDATE_DATE=bar' | less
```


`/SensorMonitorV2.xml`
```
$ http://192.168.1.5:56002/SensorMonitorV2.xml?x-KEY_MOBILE_ID=foo&x-KEY_UPDATE_DATE=
```
temperature, humidity, illuminance が取得できる。
