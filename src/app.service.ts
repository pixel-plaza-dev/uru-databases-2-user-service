import { Injectable } from '@nestjs/common';
import { UserChangePasswordDTO } from '../common/dto/users-service/user-change-password.dto';
import { UserChangeUsernameDTO } from '../common/dto/users-service/user-change-username.dto';
import { UserChangeEmailDTO } from '../common/dto/users-service/user-change-email.dto';
import { UserSendEmailVerificationTokenDTO } from '../common/dto/users-service/user-send-email-verification-token.dto';
import { UserVerifyEmailDTO } from '../common/dto/users-service/user-verify-email.dto';
import { UserForgotPasswordDTO } from '../common/dto/users-service/user-forgot-password.dto';
import { UserResetPasswordDTO } from '../common/dto/users-service/user-reset-password.dto';
import { UserUpdateDTO } from '../common/dto/users-service/user-update.dto';

@Injectable()
export class AppService {
  async updateUser(data: UserUpdateDTO) {}

  async changeUsername(data: UserChangeUsernameDTO) {}

  async changePassword(data: UserChangePasswordDTO) {}

  async changeEmail(data: UserChangeEmailDTO) {}

  async sendEmailVerificationToken(data: UserSendEmailVerificationTokenDTO) {}

  async verifyEmail(data: UserVerifyEmailDTO) {}

  async forgotPassword(data: UserForgotPasswordDTO) {}

  async resetPassword(data: UserResetPasswordDTO) {}
}
