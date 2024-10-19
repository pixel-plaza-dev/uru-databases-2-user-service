import { Injectable } from '@nestjs/common';
import { UserChangePasswordDTO } from '../common/dto/users-service/user-change-password.dto';
import { UserChangeUsernameDTO } from '../common/dto/users-service/user-change-username.dto';
import { UserChangeEmailDTO } from '../common/dto/users-service/user-change-email.dto';
import { UserSendEmailVerificationTokenDTO } from '../common/dto/users-service/user-send-email-verification-token.dto';
import { UserVerifyEmailDTO } from '../common/dto/users-service/user-verify-email.dto';
import { UserForgotPasswordDTO } from '../common/dto/users-service/user-forgot-password.dto';
import { UserResetPasswordDTO } from '../common/dto/users-service/user-reset-password.dto';
import { UserUpdateDTO } from '../common/dto/users-service/user-update.dto';
import { RmqContext } from '@nestjs/microservices';

@Injectable()
export class AppService {
  // Acknowledge the message after successful processing
  acknowledge(context: RmqContext) {
    const channel = context.getChannelRef();
    const originalMessage = context.getMessage();
    channel.ack(originalMessage);
  }

  async updateUser(data: UserUpdateDTO, context: RmqContext) {
    console.log(data)

    // Acknowledge the message after successful processing
    this.acknowledge(context);

    return 100
  }

  async changeUsername(data: UserChangeUsernameDTO, context: RmqContext) {
    // Acknowledge the message after successful processing
    this.acknowledge(context);
  }

  async changePassword(data: UserChangePasswordDTO, context: RmqContext) {
    // Acknowledge the message after successful processing
    this.acknowledge(context);
  }

  async changeEmail(data: UserChangeEmailDTO, context: RmqContext) {
    // Acknowledge the message after successful processing
    this.acknowledge(context);
  }

  async sendEmailVerificationToken(data: UserSendEmailVerificationTokenDTO, context: RmqContext) {
    // Acknowledge the message after successful processing
    this.acknowledge(context);
  }

  async verifyEmail(data: UserVerifyEmailDTO, context: RmqContext) {
    // Acknowledge the message after successful processing
    this.acknowledge(context);
  }

  async forgotPassword(data: UserForgotPasswordDTO, context: RmqContext) {
    // Acknowledge the message after successful processing
    this.acknowledge(context);
  }

  async resetPassword(data: UserResetPasswordDTO, context: RmqContext) {
    // Acknowledge the message after successful processing
    this.acknowledge(context);
  }
}
