FROM rabbitmq:alpine
RUN apk update && apk add wget
RUN wget -P /plugins https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/3.10.2/rabbitmq_delayed_message_exchange-3.10.2.ez
RUN rabbitmq-plugins enable --offline rabbitmq_mqtt rabbitmq_federation_management rabbitmq_stomp rabbitmq_delayed_message_exchange
EXPOSE 5672
EXPOSE 15672