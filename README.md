# http-mitigation


````
Usage:
   [flags]
   [command]

Available Commands:
  help        Help about any command
  version     Print the version.

Flags:
      --cookie-name string      Cookie Name (default "__mitigation")
  -h, --help                    help for this command
      --listen-port int         HTTP listen port (default 8000)
      --redis-addr string       Redis Server Address (default "127.0.0.1:6379")
      --redis-db int            Redis DB
      --redis-password string   Redis Password
      --threshold1 int          Threshold per domain per second (mitigation redirect 307) (default 10000)
      --threshold2 int          Threshold per domain per second (mitigation redirect javascript) (default 50000)
  -v, --verbose                 Enable verbose
````


Nginx Configuration

```
upstream mitigation_backend {
    server 127.0.0.1:8000;
}

server {
    listen       80;
    server_name  localhost;

    ######################## Mitigation #########################
    auth_request /__protection;
    auth_request_set $challenge $upstream_http_x_challenge;
    if ($challenge) {
        return 307 $challenge;
    }

    location /__protection {
        proxy_pass              http://mitigation_backend;
        proxy_pass_request_body off;
        proxy_set_header        Content-Length "";
        proxy_set_header        X-Original-URI $request_uri;
        proxy_set_header        X-Original-Query $query_string;
        proxy_set_header        X-Original-Host $host;
    }
    ############################################################


    location / {
        root   /usr/share/nginx/html;
        index  index.html index.htm;
    }
}

```