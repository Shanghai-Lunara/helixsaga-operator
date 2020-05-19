# hexlisaga-operator
HelixSaga Kubernetes Custom Operator


## todo
1. how to generate or edit business's configuration files

2. configMap imports all business's configuration files


3. structure:
   (1) api/game/gmt/version/pay_notify: nginx/php-fpm php:5.6  replicas
   (2) friend/rank/queue: php:5.6 swoole-ext:1.9.23,
   (3) chat/heart: php:5.6 workerman-framework
   (6) campaign/guildwar/app-notification: go binary executable files