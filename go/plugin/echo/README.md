## Echo Plugin

The Echo plugin allows you to echo a http request to the sender. 
All headers and the body will be printed as a json string:

```json
{
  "headers": {
    "header_key": [
      "header",
      "values"
    ]
  },
  "body": "Body of the request"
}
```

The endpoint is `/echo` by default, but can be configured:

```json
{
  "extra_config": {
    "plugin/http-server": {
      "name": [
        "echo"
      ],
      "echo_config": {
        "endpoint": "/my-echo"
      }
    }
  }
}
```