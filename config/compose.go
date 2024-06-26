package config

func GenDockerCompose(consul bool) string {
	if consul {
		return CONSUL_COMPOSE
	} else {
		return DIRECT_COMPOSE
	}
}

const DIRECT_COMPOSE = `version: '3.9'
services:
  dashboard:
    container_name: cossim_apisix_dashboard
    image: apache/apisix-dashboard:latest
    restart: always
    volumes:
      - ./config/common/apisix-dashboard.yaml:/usr/local/apisix-dashboard/conf/conf.yaml:ro
    depends_on:
      - etcd
    ports:
      - "8000:9000/tcp"
  livekit:
    image: livekit/livekit-server
    container_name: cossim_livekit
    ports:
      - "7880:7880"
      - "7881:7881"
      - "7882:7882/udp"
    volumes:
       - ./config/common/livekit.yaml:/livekit.yaml
    command:
      - "--config"
      - "/livekit.yaml"
      - "--node-ip=127.0.0.1"
  apisix:
    image: apache/apisix
    container_name: cossim_apisix
    restart: always
    volumes:
      - ./config/common/apisix.yaml:/usr/local/apisix/conf/config.yaml:ro
    depends_on:
      - etcd
    ports:
      - "9180:9180/tcp"
      - "8080:9080/tcp"
      - "9091:9091/tcp"
      - "443:9443/tcp"
      - "9092:9092/tcp"
  etcd:
    container_name: cossim_etcd
    image: bitnami/etcd:3.4.9
    user: root
    restart: always
    volumes:
      - ./data/var/lib/etcd_data:/etcd_data
    environment:
      ETCD_DATA_DIR: /etcd_data
      ETCD_ENABLE_V2: "true"
      ALLOW_NONE_AUTHENTICATION: "yes"
      ETCD_ADVERTISE_CLIENT_URLS: "http://etcd:2379"
      ETCD_LISTEN_CLIENT_URLS: "http://0.0.0.0:2379"
    ports:
      - "2379:2379/tcp"
  mysql:
    container_name: cossim_mysql
    image: mysql:5.7
    volumes:
      - ./data/var/lib/mysql:/var/lib/mysql
    command: [
      '--character-set-server=utf8mb4',
    ]
    expose:
      - "3306"
    ports:
      - "33306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "Hitosea@123.."
      MYSQL_DATABASE: coss
      MYSQL_USER: coss
      MYSQL_PASSWORD: "Hitosea@123.."
      MYSQL_TCP_PORT: '3306'
      MYSQL_ROOT_HOST: '%'
    #      MARIADB_AUTO_UPGRADE: 'true'
    #      MARIADB_DISABLE_UPGRADE_BACKUP: 'true'
    healthcheck:
      test: mysqladmin ping -h mysql -P 3306 -p$$MYSQL_ROOT_PASSWORD
      interval: 5s
      timeout: 10s
      retries: 10
      start_period: 30s
  minio:
    image: hub.hitosea.com/cossim/minio
    container_name: cossim_minio
    ports:
      - "9000:9000"
      - "9001:9001"
    expose:
      - "9000"
      - "9001"
    environment:
      MINIO_ROOT_USER: root
      MINIO_ROOT_PASSWORD: Hitosea@123..
    volumes:
      - ./data/var/lib/minio:/data
    command: server /data --console-address ":9001"
  rabbitmq:
    image: "rabbitmq:management"
    container_name: cossim_rabbitmq
    hostname: rabbitmq3-management-master
    logging:
      driver: json-file
      options:
        max-size: "100m"
        max-file: "1"
    volumes:
      - "./data/var/lib/rabbitmq:/var/lib/rabbitmq"
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: root
      RABBITMQ_DEFAULT_PASS: Hitosea@123..
  redis:
    image: redis:latest
    container_name: cossim_redis
    ports:
      - 6379:6379
    command: redis-server --requirepass Hitosea@123..
  dtm:
    container_name: cossim_dtm
    image: hub.hitosea.com/cossim/dtm
    ports:
      - '36789:36789'
      - '36790:36790'
  consul:
    image: hub.hitosea.com/cossim/consul:latest
    container_name: cossim_consul
    volumes:
      - ./data/var/lib/consul:/consul/data
      - ./config/common/consul.json:/etc/consul/consul-config.json
    command: consul agent -server -bootstrap-expect=1 -client=0.0.0.0 -ui -data-dir=/consul/data -config-dir=/etc/consul
    ports:
      - '8300:8300'
      - '8301:8301'
      - '8301:8301/udp'
      - '8500:8500'
      - '8600:8600'
      - '8600:8600/udp'
  admin:
    container_name: cossim_admin
    image: hub.hitosea.com/cossim/admin
    command:
      - "--config"
      - "/config/config.yaml"
      - "--discover"
    volumes:
      - ./config/service/admin.yaml:/config/config.yaml
      - ./config/pgp:/.cache
    depends_on:
      - dtm
      - redis
      - rabbitmq
      - group
      - relation
      - user
      - msg
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  live:
    container_name: cossim_live
    image: hub.hitosea.com/cossim/live
    command:
      - "--config"
      - "/config/config.yaml"
      - "--discover"
    volumes:
      - ./config/service/live.yaml:/config/config.yaml
      - ./config/pgp:/.cache
    depends_on:
      - consul
      - redis
      - user
      - relation
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  msg:
    container_name: cossim_msg
    image: hub.hitosea.com/cossim/msg
    command:
      - "--config"
      - "/config/config.yaml"
      - "--discover"
    depends_on:
      mysql:
        condition: service_healthy
    volumes:
      - ./config/service/msg.yaml:/config/config.yaml
      - ./config/pgp:/.cache
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  relation:
    container_name: cossim_relation
    image: hub.hitosea.com/cossim/relation
    command:
      - "--config"
      - "/config/config.yaml"
      - "--discover"
    volumes:
      - ./config/service/relation.yaml:/config/config.yaml
      - ./config/pgp:/.cache
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  user:
    container_name: cossim_user
    image: hub.hitosea.com/cossim/user
    command:
      - "--config"
      - "/config/config.yaml"
      - "--discover"
    volumes:
      - ./config/service/user.yaml:/config/config.yaml
      - ./config/pgp:/.cache
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  group:
    container_name: cossim_group
    image: hub.hitosea.com/cossim/group
    command:
      - "--config"
      - "/config/config.yaml"
      - "--discover"
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
    volumes:
      - ./config/service/group.yaml:/config/config.yaml
      - ./config/pgp:/.cache
  storage:
    container_name: cossim_storage
    image: hub.hitosea.com/cossim/storage
    command:
      - "--config"
      - "/config/config.yaml"
      - "--discover"
    volumes:
      - ./config/service/storage.yaml:/config/config.yaml
      - ./config/pgp:/.cache
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  push:
    container_name: cossim_push
    image: hub.hitosea.com/cossim/push
    command:
      - "--config"
      - "/config/config.yaml"
      - "--discover"
    volumes:
      - ./config/service/push.yaml:/config/config.yaml
      - ./config/pgp:/.cache
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
`

const CONSUL_COMPOSE = `version: '3.9'
services:
  dashboard:
    container_name: cossim_apisix_dashboard
    image: apache/apisix-dashboard:latest
    restart: always
    volumes:
      - ./config/common/apisix-dashboard.yaml:/usr/local/apisix-dashboard/conf/conf.yaml:ro
    depends_on:
      - etcd
    ports:
      - "8000:9000/tcp"
  livekit:
    image: livekit/livekit-server
    container_name: cossim_livekit
    ports:
      - "7880:7880"
      - "7881:7881"
      - "7882:7882/udp"
    volumes:
       - ./config/common/livekit.yaml:/livekit.yaml
    command:
      - "--config"
      - "/livekit.yaml"
      - "--node-ip=127.0.0.1"
  apisix:
    image: apache/apisix
    container_name: cossim_apisix
    restart: always
    volumes:
      - ./config/common/apisix.yaml:/usr/local/apisix/conf/config.yaml:ro
    depends_on:
      - etcd
    ports:
      - "9180:9180/tcp"
      - "8080:9080/tcp"
      - "9091:9091/tcp"
      - "443:9443/tcp"
      - "9092:9092/tcp"
  etcd:
    container_name: cossim_etcd
    image: bitnami/etcd:3.4.9
    user: root
    restart: always
    volumes:
      - ./data/var/lib/etcd_data:/etcd_data
    environment:
      ETCD_DATA_DIR: /etcd_data
      ETCD_ENABLE_V2: "true"
      ALLOW_NONE_AUTHENTICATION: "yes"
      ETCD_ADVERTISE_CLIENT_URLS: "http://etcd:2379"
      ETCD_LISTEN_CLIENT_URLS: "http://0.0.0.0:2379"
    ports:
      - "2379:2379/tcp"
  mysql:
    container_name: cossim_mysql
    image: mysql:5.7
    volumes:
      - ./data/var/lib/mysql:/var/lib/mysql
    command: [
      '--character-set-server=utf8mb4',
    ]
    expose:
      - "3306"
    ports:
      - "33306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "Hitosea@123.."
      MYSQL_DATABASE: coss
      MYSQL_USER: coss
      MYSQL_PASSWORD: "Hitosea@123.."
      MYSQL_TCP_PORT: '3306'
      MYSQL_ROOT_HOST: '%'
    #      MARIADB_AUTO_UPGRADE: 'true'
    #      MARIADB_DISABLE_UPGRADE_BACKUP: 'true'
    healthcheck:
      test: mysqladmin ping -h mysql -P 3306 -p$$MYSQL_ROOT_PASSWORD
      interval: 5s
      timeout: 10s
      retries: 10
      start_period: 30s
  minio:
    image: hub.hitosea.com/cossim/minio
    container_name: cossim_minio
    ports:
      - "9000:9000"
      - "9001:9001"
    expose:
      - "9000"
      - "9001"
    environment:
      MINIO_ROOT_USER: root
      MINIO_ROOT_PASSWORD: Hitosea@123..
    volumes:
      - ./data/var/lib/minio:/data
    command: server /data --console-address ":9001"
  rabbitmq:
    image: "rabbitmq:management"
    container_name: cossim_rabbitmq
    hostname: rabbitmq3-management-master
    logging:
      driver: json-file
      options:
        max-size: "100m"
        max-file: "1"
    volumes:
      - "./data/var/lib/rabbitmq:/var/lib/rabbitmq"
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: root
      RABBITMQ_DEFAULT_PASS: Hitosea@123..
  redis:
    image: redis:latest
    container_name: cossim_redis
    ports:
      - 6379:6379
    command: redis-server --requirepass Hitosea@123..
  dtm:
    container_name: cossim_dtm
    image: hub.hitosea.com/cossim/dtm
    ports:
      - '36789:36789'
      - '36790:36790'
  consul:
    image: hub.hitosea.com/cossim/consul:latest
    container_name: cossim_consul
    volumes:
      - ./data/var/lib/consul:/consul/data
      - ./config/common/consul.json:/etc/consul/consul-config.json
    command: consul agent -server -bootstrap-expect=1 -client=0.0.0.0 -ui -data-dir=/consul/data -config-dir=/etc/consul
    ports:
      - '8300:8300'
      - '8301:8301'
      - '8301:8301/udp'
      - '8500:8500'
      - '8600:8600'
      - '8600:8600/udp'
  admin:
    container_name: cossim_admin
    image: hub.hitosea.com/cossim/admin
    command:
      - "--discover"
      - "--register"
      - "--remote-config"
      - "--hot-reload"
      - "--config-center-addr"
      - "consul:8500"
      - "--pprof-bind-address"
      - ":6060"
      - "--metrics-bind-address"
      - ":9090"
      - "--http-health-probe-bind-address"
      - ":9091"
      - "--grpc-health-probe-bind-address"
      - ":9092"
    volumes:
      - ./config/pgp:/.cache
    depends_on:
      - dtm
      - redis
      - rabbitmq
      - group
      - relation
      - user
      - msg
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  live:
    container_name: cossim_live
    image: hub.hitosea.com/cossim/live
    command:
      - "--discover"
      - "--register"
      - "--remote-config"
      - "--hot-reload"
      - "--config-center-addr"
      - "consul:8500"
      - "--pprof-bind-address"
      - ":6060"
      - "--metrics-bind-address"
      - ":9090"
      - "--http-health-probe-bind-address"
      - ":9091"
      - "--grpc-health-probe-bind-address"
      - ":9092"
    volumes:
      - ./config/pgp:/.cache
    depends_on:
      - consul
      - redis
      - user
      - relation
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  msg:
    container_name: cossim_msg
    image: hub.hitosea.com/cossim/msg
    command:
      - "--discover"
      - "--register"
      - "--remote-config"
      - "--hot-reload"
      - "--config-center-addr"
      - "consul:8500"
      - "--pprof-bind-address"
      - ":6060"
      - "--metrics-bind-address"
      - ":9090"
      - "--http-health-probe-bind-address"
      - ":9091"
      - "--grpc-health-probe-bind-address"
      - ":9092"
    depends_on:
      mysql:
        condition: service_healthy
    volumes:
      - ./config/pgp:/.cache
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  relation:
    container_name: cossim_relation
    image: hub.hitosea.com/cossim/relation
    command:
      - "--discover"
      - "--register"
      - "--remote-config"
      - "--hot-reload"
      - "--config-center-addr"
      - "consul:8500"
      - "--pprof-bind-address"
      - ":6060"
      - "--metrics-bind-address"
      - ":9090"
      - "--http-health-probe-bind-address"
      - ":9091"
      - "--grpc-health-probe-bind-address"
      - ":9092"
    volumes:
      - ./config/pgp:/.cache
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  user:
    container_name: cossim_user
    image: hub.hitosea.com/cossim/user
    command:
      - "--discover"
      - "--register"
      - "--remote-config"
      - "--hot-reload"
      - "--config-center-addr"
      - "consul:8500"
      - "--pprof-bind-address"
      - ":6060"
      - "--metrics-bind-address"
      - ":9090"
      - "--http-health-probe-bind-address"
      - ":9091"
      - "--grpc-health-probe-bind-address"
      - ":9092"
    volumes:
      - ./config/pgp:/.cache
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  group:
    container_name: cossim_group
    image: hub.hitosea.com/cossim/group
    command:
      - "--discover"
      - "--register"
      - "--remote-config"
      - "--hot-reload"
      - "--config-center-addr"
      - "consul:8500"
      - "--pprof-bind-address"
      - ":6060"
      - "--metrics-bind-address"
      - ":9090"
      - "--http-health-probe-bind-address"
      - ":9091"
      - "--grpc-health-probe-bind-address"
      - ":9092"
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
    volumes:
      - ./config/pgp:/.cache
  storage:
    container_name: cossim_storage
    image: hub.hitosea.com/cossim/storage
    command:
      - "--discover"
      - "--register"
      - "--remote-config"
      - "--hot-reload"
      - "--config-center-addr"
      - "consul:8500"
      - "--pprof-bind-address"
      - ":6060"
      - "--metrics-bind-address"
      - ":9090"
      - "--http-health-probe-bind-address"
      - ":9091"
      - "--grpc-health-probe-bind-address"
      - ":9092"
    volumes:
      - ./config/pgp:/.cache
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
  push:
    container_name: cossim_push
    image: hub.hitosea.com/cossim/push
    command:
      - "--discover"
      - "--register"
      - "--remote-config"
      - "--hot-reload"
      - "--config-center-addr"
      - "consul:8500"
      - "--pprof-bind-address"
      - ":6060"
      - "--metrics-bind-address"
      - ":9090"
      - "--http-health-probe-bind-address"
      - ":9091"
      - "--grpc-health-probe-bind-address"
      - ":9092"
    volumes:
      - ./config/pgp:/.cache
    environment:
      CONSUL_HTTP_TOKEN:
    restart: on-failure
`
