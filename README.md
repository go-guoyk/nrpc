# nrpc
low overhead rpc library

## Wire Format

`nrpc` uses text-based wire format, thus it's `nc` friendly.

**Request (no body)**

```text
flake,create\n
instance_id: localhost\n
track_id: acde38fb6\n
\n
\n
```

**Request (JSON body)**

```text
user,create\n
instance_id: localhost\n
track_id: acde38fb6\n
\n
{"name":"guoyk"}\n
```

**Response (no body)**

```text
ok,custom success message 自定义成功信息\n
instance_id: localhost\n
track_id: acde38fb6\n
\n
\n
```

**Response (no body, with error)**

```text
err_internal,custom error message 自定义错误信息\n
instance_id: localhost\n
track_id: acde38fb6\n
\n
\n
```


**Response (JSON body)**

```text
ok,custom success message\n
instance_id: localhost\n
track_id: acde38fb6\n
\n
{"user_id": 1}\n
```
