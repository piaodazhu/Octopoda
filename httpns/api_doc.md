## 1 ping
- protocol: http, https
- url: /ping
- method: GET
- argument: none
- response: none. statuscode 200


## 2 register a name entry
- protocol: https
- url: /register
- method: POST
- argument: in post json body
- request body:
```json
[
	{
		"key": "a string",
		"type": "a string, one of addr, num, str, array",
		"value": "a string",
		"description": "description",
		"ttl": "a int, means ms",
	},
	{
		"key": "a string",
		"type": "a string, one of addr, num, str, array",
		"value": "a string",
		"description": "description",
		"ttl": "a int, means ms",
	}
]
```
- response:
```json
{
	"msg":"OK"
}
```


## 3 query a name entry
- protocol: http, https
- url: /query
- method: GET
- argument: in query
- fields:
	- name (string, required): it is the key of this record.
- response:
```json
{
	"msg":"OK",
	"entry":{
		"key": "brain",
		"type": "string",
		"value": "1.1.1.1",
		"description": "ip of brain",
		"ttl": 3000,
	}
}
```

## 4 delete a name record
- protocol: https
- url: /delete
- method: POST
- argument: in post form
- fields:
	- name (string, required): it is the key of this record.
- response:
```json
{
	"msg":"OK"
}
```

## 5 list stored keys
- protocol: https
- url: /list
- method: GET
- argument: in query
- fields:
	- match (string, required): string for pattern matching.
	- method (string, required): must be one of `prefix` `suffix` `contain` `contain` `all`, for making pattern.
- response:
```json
{
	"msg":"OK",
	"list":[
		"net01.brain01",
		"net01.tentacle01",
		"net01.tentacle02",
		"net01.tentacle03",
	]
}
```

## 6 summery of service
- protocol: https
- url: /summery
- method: GET
- argument: no
- response:
```json
{
	"total_request": 666,
	"since": 12345124,
	"api_stats":{
		"register[GET]": {
			"requests": 100,
			"success": 100,
		},
		"register[POST]": {
			"requests": 80,
			"success": 78,
		},
		"sshinfo[GET]": {
			"requests": 30,
			"success": 30,
		},
	}
}
```