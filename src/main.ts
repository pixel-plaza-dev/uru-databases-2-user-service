import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { RABBITMQ } from './config/microservice';

async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    AppModule,
    {
      transport: Transport.RMQ,
      options: {
        urls: [RABBITMQ.URL],
        queue: RABBITMQ.USERS_QUEUE,
        queueOptions: {
          durable: true,
        },
      },
    },
  );
  await app.listen();
}

bootstrap();
