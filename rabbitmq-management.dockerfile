# Sử dụng một base image có sẵn với RabbitMQ đã cài đặt và kích hoạt Management Plugin
FROM rabbitmq:3-management

# Thiết lập các biến môi trường cho RabbitMQ Management Plugin
ENV RABBITMQ_DEFAULT_USER=admin
ENV RABBITMQ_DEFAULT_PASS=admin

# EXPOSE cổng 15672 (cho giao diện web quản lý)
EXPOSE 15672

# Bật RabbitMQ Management Plugin
RUN rabbitmq-plugins enable rabbitmq_management

# Đảm bảo RabbitMQ Management Plugin đã được kích hoạt và đang chạy khi container được khởi chạy
CMD ["rabbitmq-server"]