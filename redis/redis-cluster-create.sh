# wait for the docker-compose depends_on to spin up the redis nodes usually takes this long
sleep 10



redis-cli --cluster create 127.0.0.1:6372 127.0.0.1:6373 127.0.0.1:6374 127.0.0.1:6375 127.0.0.1:6376 127.0.0.1:6377 --cluster-replicas 1 --cluster-yes