# quickapi-rpc
RPC take on the quickapi stuff.

As with quickapi, there's 2 ways to get the thing going.

1. `Run` which you feed your entities. It will start a rpc-server and expose it over nats according to nuts.
2.  `For` which feed one of your entities binds it to topic `<prefix>.<entity.Name()>.<create|read|update|delete...>`.

## Differences

Obviously there are differencenes. Quickapi-rpc have it's own sub-api. Ie `Create` expects a `api.Create`-struct and so on. How that struct is created and sent, is your problem 🙃.