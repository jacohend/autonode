# Autonode

Peer-to-peer event bus and queue in golang, for creating horizontally-scalable services.

## Installation

See the example for an autonode implementation.

To test the example, run

```
make build

./autonode --autonode.listen="<YourIP>:8082" --addr="0.0.0.0:8081"
```

and on another computer, run 

```
make build

./autonode --autonode.listen="<YourIP>:8082" --addr="0.0.0.0:8081" --autonode.seeds="<The FirstIP>:8082"
```

`curl -v http://127.0.0.1:8081/` will result in 
an event being acked and processed by the other node