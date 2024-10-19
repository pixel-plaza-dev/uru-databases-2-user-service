import { Controller } from '@nestjs/common';
import { AppService } from './app.service';
import { Ctx, MessagePattern, Payload, RmqContext } from '@nestjs/microservices';
import USERS_PATTERN from '../common/message-pattern/users-service';
import { UserUpdateDTO } from '../common/dto/users-service/user-update.dto';
import { UserChangePasswordDTO } from '../common/dto/users-service/user-change-password.dto';
import { UserChangeUsernameDTO } from '../common/dto/users-service/user-change-username.dto';
import { UserChangeEmailDTO } from '../common/dto/users-service/user-change-email.dto';
import { UserSendEmailVerificationTokenDTO } from '../common/dto/users-service/user-send-email-verification-token.dto';
import { UserVerifyEmailDTO } from '../common/dto/users-service/user-verify-email.dto';
import { UserForgotPasswordDTO } from '../common/dto/users-service/user-forgot-password.dto';
import { UserResetPasswordDTO } from '../common/dto/users-service/user-reset-password.dto';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  @MessagePattern(USERS_PATTERN.UPDATE_USER)
  async updateUser(@Payload() data: UserUpdateDTO, @Ctx() context: RmqContext) {
    return this.appService.updateUser(data, context);
  }

  @MessagePattern(USERS_PATTERN.CHANGE_USERNAME)
  async changeUsername(@Payload() data: UserChangeUsernameDTO, @Ctx() context: RmqContext) {
    return this.appService.changeUsername(data, context);
  }

  @MessagePattern(USERS_PATTERN.CHANGE_PASSWORD)
  async changePassword(@Payload() data: UserChangePasswordDTO, @Ctx() context: RmqContext) {
    return this.appService.changePassword(data, context);
  }

  @MessagePattern(USERS_PATTERN.CHANGE_EMAIL)
  async changeEmail(@Payload() data: UserChangeEmailDTO, @Ctx() context: RmqContext) {
    return this.appService.changeEmail(data, context);
  }

  @MessagePattern(USERS_PATTERN.SEND_EMAIL_VERIFICATION_TOKEN)
  async sendEmailVerificationToken(@Payload() data: UserSendEmailVerificationTokenDTO, @Ctx() context: RmqContext) {
    return this.appService.sendEmailVerificationToken(data, context);
  }

  @MessagePattern(USERS_PATTERN.VERIFY_EMAIL)
  async verifyEmail(@Payload() data: UserVerifyEmailDTO, @Ctx() context: RmqContext) {
    return this.appService.verifyEmail(data, context);
  }

  @MessagePattern(USERS_PATTERN.FORGOT_PASSWORD)
  async forgotPassword(@Payload() data: UserForgotPasswordDTO, @Ctx() context: RmqContext) {
    return this.appService.forgotPassword(data, context);
  }

  @MessagePattern(USERS_PATTERN.RESET_PASSWORD)
  async resetPassword(@Payload() data: UserResetPasswordDTO, @Ctx() context: RmqContext) {
    return this.appService.resetPassword(data, context);
  }
}
