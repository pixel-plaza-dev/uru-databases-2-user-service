// User Microservice configuration
export const IS_PRODUCTION = process.env.NODE_ENV === 'production';
export const USERS_MICROSERVICE_PORT = parseInt(
  process.env.USERS_MICROSERVICE_PORT,
);

// RabbitMQ configuration
export const RABBITMQ = {
  URL: process.env.RABBITMQ_URL,
  USERS_QUEUE: process.env.RABBITMQ_USERS_QUEUE,
};
