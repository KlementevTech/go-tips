# go-tips

## Pprof сервер

[Pprof](http://localhost:6060/debug/pprof/)

## Примеры запросов

```shell
grpcurl -plaintext \
    -d '{"id": "019d5da1-46dd-7b0d-82e5-49345ac87e78", "name": "Turbo21"}' \
    127.0.0.1:50051 \
    catalog.v1.PcPartService/CreatePcPart
```

```shell
grpcurl -plaintext \
    -d '{"id": "019d5da1-46dd-7b0d-82e5-49345ac87e78"}' \
    127.0.0.1:50051 \
    catalog.v1.PcPartService/GetPcPart
```

```shell
grpcurl -plaintext \
    -d '{"pc_part": {"id": "019d5da1-46dd-7b0d-82e5-49345ac87e78", "version": "1", "name": "turbo2"}}' \
    127.0.0.1:50051 \
    catalog.v1.PcPartService/UpdatePcPart
```

```shell
grpcurl -plaintext \
    -d '{"id": "019d5da1-46dd-7b0d-82e5-49345ac87e78", "version": "2"}' \
    127.0.0.1:50051 \
    catalog.v1.PcPartService/DeletePcPart
```