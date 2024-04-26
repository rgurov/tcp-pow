# tcp-pow

tcp-pow requires docker compose installed

How to run:

```bash
make docker
```

Example output

```
server  | 2024/04/19 11:50:39 INFO server started on 0.0.0.0:7771
client  | 2024/04/19 11:50:39 INFO connected
server  | 2024/04/19 11:50:39 INFO new connection ip=172.21.0.3:41816
server  | 2024/04/19 11:50:39 INFO puzzle sent, hash=fbc2b99dd1b40411eddcf9224c79215e1629e0aec4f4be4eed9a9aa30f7f1af1 ip=172.21.0.3:41816
client  | 2024/04/19 11:50:39 INFO puzzle hash received
client  | 2024/04/19 11:50:39 INFO solving the puzzle...
client  | 34242381
client  | 2024/04/19 11:50:41 INFO puzzle solved, hash=00000073a2e84535c423cc8655fe9f01007a831a313338dabc0f1a73fcff1176
server  | 2024/04/19 11:50:41 INFO valid solution received ip=172.21.0.3:41816
server  | 2024/04/19 11:50:41 INFO word of wisdom sent ip=172.21.0.3:41816
server  | 2024/04/19 11:50:41 INFO connection closed ip=172.21.0.3:41816
client  | 2024/04/19 11:50:41 INFO got message from server: Admit your mistakes and don’t repeat them. If you can’t admit your mistakes, you are destined to repeat them.
client  | 2024/04/19 11:50:41 INFO done! waiting for signal to close
```


Or you can run *client* and *server* with:

```bash
make client
make server
```

Flags for server and client:

```
-h Host       (default 127.0.0.1)
-p Port       (default 7771)
-c Complexity (default 6)
```

## How it works

1. Server starts listening for a incoming connections
2. Client connects to server
3. Client received a random hash
4. Client compute a solution (hash) with N leading zeros. N - complexity
5. Client sends solution to server
6. Server verifies that solution is correct
7. Server sends random "Word of Wisdom"
8. Client receive message and wait for term signal to close

This proof of work alghorithm was choosen for it's simplicity and widespread  

To run tests use this:
```bash
make test
```
