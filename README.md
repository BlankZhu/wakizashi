# wakizashi

Wakizashi is a traffic probe for cloud native environment. It allows you to sniff the bytes passing through your machine (or in K8S, your pod). 

This repository provides wakizashi that only sinffs the amount of bytes running through, if you want a traffic content analyzer, you may fork then build your own version. It won't be hard to make your own traffic filter logic.

## Consist

Wakizashi consists two parts: `center` & `probe`.

### Probe

`probe` is a traffic probe that captures every byte going through given network devices. It can be deployed on any machine, or alongside with other continer in a K8S pod. Once starts working, `probe` will read the amount of bytes going through given network devices like `eth0`, `tunl0` and so on. Later, the traffic status will be posted to `center` for aggregation.

### Center
`center` works as a aggregator of traffic data. The data from `probe` will be interpreted into structural data record, and later written in backend databse like InfluxDB, redis or MongoDB. It's OK to deploy `center` in multiple replicas.

## Get Started

For now, `wakizashi` only supports InfluxDB1 as data backend. So we will get started like below:

```txt
+-------+-------+               +--------+              +------------+
| Nginx | Probe |               |        |+             |            |
+-------+-------+  =Raw Data=>  | Center ||  =Record=>  |  Backend   |
|  Pod (or VM)  |               |        ||             | (InfluxDB) |
+---------------+               +--------+|             +------------+
                                 +--------+
```

`Nginx` here represents a actual working application, which is deployed together with a `probe` to collect network traffic data. `probe` will ignore the raw datas between `center` and itself to avoid confusing final records.

`Center` is deployed in other network space (other vm, or pod) to avoid confusing raw data. It communicate with `probe` in `GRPC`, so it's easy to make load-balance between them.

`Backend` stores the aggregated records reported from `center`. Users can analyse the traffic data from it easily.

For `Nginx` and `InfluxDB`, check their official documents to deploy them.

### Config

First deploy `center` in following command:
```shell
./center -c ./center-config.yaml
```
For configuration example check `config/center-config.yaml`.

Then deploy `probe`:
```shell
./probe -c ./probe-config.yaml
```
For configuration example check `config/probe-config.yaml`.

### Validate

Once `Nginx` (as user application) alongside `probe`, `center` and `backend` are all up, make a request to Nginx by CURL, after a period of time (defined the configuration of `center` & `probe`), you will see the record in `backend`.

## Build

### Binary

If you want a binary version that runs directory on your machine, you need to install Go 1.16 (or later), and just run:

```shell
./build-binary.sh
```

It gives:

```shell
generating wakizashi center
wakizashi center generated
generating wakizashi probe
wakizashi probe generate
```

Then you can find the binary in `build/`

### Container Image

If you want a container version (to work in K8S), just run:

```shell
./build-container.sh [TAG]
```

Specially for CN network environment, use `build-container-cn.sh`

## More

For those who want to make a real *sniffer*, you will need to focus on these core codes:
- pkg/dump/dump.go for data dumping behaviour
- pkg/entity/rawtrafficrecord.go for raw traffic between `probe` & `center`
- pkg/entity/trafficrecord.go which represents the data in `backend`

## And More

`wakizashi` is a sidetime project and is tested on K8S. Feel free to make any issue or pr.