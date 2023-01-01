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

### Send data to device/id room

```
curl http://localhost:3000/device/publish 
```

Output

```json
{"id":"1e4832e7-1ffa-4cf4-b9d9-0b8eff286c52","name":"Temp"}
```