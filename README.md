# GO Socket.IO Example

Authorization Socket.IO with JWT Example

### Run server

```shell
go run .
```

### Open browser

```shell
http://localhost:3000/
```

### Test with API

```
curl http://localhost:3000/device/publish 
```

Output

```json
{"id":"1e4832e7-1ffa-4cf4-b9d9-0b8eff286c52","name":"Temp"}
```

### Test with MQTTX

```json
url: wss://broker.emqx.io:8084/mqtt
username: emqx_test
password: emqx_test
topic: device
payload: {"id": "1",  "msg": "hello"}
```