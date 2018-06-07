#### Should be started with 'docker-compose up' command.
* By default main service works on localhost:8000
* Main service uses in-memory cache system that can be configured through config.json file.

* In case of all services started not by docker-compose but from the host system, main service should be run with '-l' flag, that will set outer services urls to localhost.

#### Existing resources:
    localhost:8000/info/user1
    localhost:8000/info/user2
    localhost:8000/info/user3