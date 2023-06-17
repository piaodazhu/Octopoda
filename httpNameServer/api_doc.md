## 0 ping
- url: /ping
- method: GET
- argument: none
- response: none. statuscode 200
```

## 1 register a name entry
- url: /register
- method: POST
- argument: in post form
- fields:
	- type (string, required): must be one of `brain` `tentacle` `octl` `other`.
	- name (string, required): it is the key of this record.
	- ip (string, required): should be a ip address.
	- port (int, required): should be a integer in 0-65535.
	- ttl (int, optional)
	- description (string, optional)
- response:
```json
{
	"msg":"OK"
}
```


## 2 query a name entry
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
		"type":"brain",
		"name":"net01.brain01",
		"ip":"1.1.1.1",
		"port":10101,
		"description":"hello world",
		"ts":12345124,
	}
}
```

## 3 upload a config file
- url: /conf
- method: POST
- argument: in post form
- fields:
	- type (string, required): must be one of `brain` `tentacle` `octl` `other`.
	- name (string, required): it is the key of this config list record.
	- method (string, required): must be one of `reset` `append` `clear`.
		- reset: delete all historical configs of this name and save current conf.
		- append: append current conf to configs of this name.
		- clear: delete all historical configs of this name.
	- conf (int, optional): the serailized configuration file.
- response:
```json
{
	"msg":"OK"
}
```

## 4 download a config file
- url: /conf
- method: GET
- argument: in query
- fields:
	- name (string, required): it is the key of this config list record.
	- index (int, required): starting index of configs of this name. **0 means the latest.**
	- amount (int, requried): number of configs to get from given index.
- response:
```json
{
	"msg":"OK",
	"conflist":[{
		"type":"tentacle",
		"name":"net01.tentacle03",
		"conf":"{'a':'b'}",
		"ts":12345161,
	},{
		"type":"tentacle",
		"name":"net01.tentacle03",
		"conf":"{'a':'b'}",
		"ts":12345120,
	}]
}
```

## 5 upload a ssh info
- url: /sshinfo
- method: POST
- argument: in post form
- fields:
	- type (string, required): must be one of `brain` `tentacle` `octl` `other`.
	- name (string, required): it is the key of this config list record.
	- username (string, required): the ssh login username.
	- ip (string, required): the ssh login ip address.
	- port (int, required): the ssh login port number.
	- password (string, required): the ssh login port password.
- response:
```json
{
	"msg":"OK"
}
```

## 6 download a ssh info
- url: /sshinfo
- method: GET
- argument: in query
- fields:
	- name (string, required): it is the key of this record.
- response:
```json
{
	"msg":"OK",
	"sshinfo":{
		"type":"brain",
		"name":"net01.tentacle02",
		"username":"pi",
		"ip":"2.2.2.2",
		"port":10101,
		"password":"123456",
		"ts":12345124,
	}
}
```

## 7 list stored keys
- url: /list
- method: GET
- argument: in query
- fields:
	- scope	(string, required): must be one of `name` `config` `ssh`, means the 3 scopes of records mentioned above.
	- match (string, required): for pattern matching.
	- method (string, required): must be one of `prefix` `suffix` `contain` `contain` `all`, for making pattern.
- response:
```json
{
	"msg":"OK",
	"conflist":[
		"net01.brain01",
		"net01.tentacle01",
		"net01.tentacle02",
		"net01.tentacle03",
	]
}
```

## 8 summery of service
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