
# Pushprom

Pushprom is a proxy (http/udp) to the prometheus go client. But why you need a proxy? Prometheus doesn't offers a php client and php clients are hard to implement because they would need to "hold state" (to count things for example) and php/apache setups don't encorage that.

Pushprom is different from [pushgateway](https://github.com/prometheus/pushgateway) because it's an aggregator. Pushgateway just "stores" the state.

We offer two flavors of php clients for it. [Vanilla](https://github.com/messagebird/pushprom-php-client) and [Yii2](https://github.com/messagebird/pushprom-yii2-client).

It accepts http and udp requests. The payloads are in json. Here is a full example:


```json
{
        "type":   "gauge",
        "name":   "trees",
        "help":   "the amount of trees in the forest.",
        "method": "add",
        "value":  3002,
        "labels": {
                "species": "araucaria angustifolia",
                "job": "tree-counter-bot"
        }
}
```

When pushprom receives this payload (from now on called Delta) it tries to register the metric with type **Gauge** named **trees** and then apply the operation **add** with value **3002** on it.

# Usage

```
$ pushprom -h
Usage of bin/pushprom:
  -debug
        Log debugging messages.
  -http-listen-address string
        The address to listen on for http stat and telemetry requests. (default ":9091")
  -udp-listen-address string
        The address to listen on for udp stats requests. (default ":9090")
```

run it:

```
$ pushprom
2016/01/15 16:35:25 main.go:26: exposing metrics on http://127.0.0.1::9092/metrics
2016/01/15 16:35:25 http.go:34: listening for stats on http://127.0.0.1:9091
2016/01/15 16:35:25 udp.go:10: listening for stats UDP on port :9090
```


# Protocol support

You can use http requests and udp packages to push deltas to pushprom.

## HTTP

When using http you should to a ```POST /```.

Watch out for error messages on the response body.

Example:

```bash
curl -H "Content-type: application/json" -X POST -d '{"type":"conter", "name":"gophers", "help":"little burrowing rodents", "method":"inc"}' http://localhost:9091/
```

## UDP

You move fast and break things.

Things may break: in the prometheus go client if you cannot register a metric with the same **name** and different **help** or **labels**. For example you
 register a metric with name "gophers" and with help "little rodents" a little later you think "but they are also burrowing animals!" and then you change the help string and push the same metric it wont work: you need to reboot pushprom.

Example:

```bash
echo "{\"type\":\"counter\", \"name\":\"gophers\", \"help\":\"little burrowing rodents\", \"method\":\"inc\"}" | nc -u -w1 127.0.0.1 9090
```

# Tests

```
go tests
```


