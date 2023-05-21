import {Component} from '@angular/core';
import {NavController} from '@ionic/angular';
import {HttpClient} from '@angular/common/http';
import {MessagesService} from '../messages.service';
import {
  create,
  CredentialCreationOptionsJSON,
  parseCreationOptionsFromJSON
} from "@github/webauthn-json/browser-ponyfill";
import {environment} from '../../environments/environment';
import {Errors, UsernameInput} from "../api/types";
import {NgForm, NgModel} from "@angular/forms";
import {displayFieldErrors} from "../util";
import {Router} from "@angular/router";

@Component({
  selector: 'app-sign-up',
  templateUrl: './sign-up.page.html'
})
export class SignUpPage {
  constructor(private readonly navCtrl: NavController,
              private readonly router: Router,
              private readonly messagesService: MessagesService,
              private readonly httpClient: HttpClient) {
  }

  async signUp(form: NgForm, username: string | null): Promise<void> {
    if (!username) {
      return;
    }

    const loading = await this.messagesService.showLoading('Starting sign up process...');
    await loading.present();

    const userNameInput: UsernameInput = {username};

    this.httpClient.post<CredentialCreationOptionsJSON>(`${environment.API_URL}/signup/start`, userNameInput)
      .subscribe({
        next: async (response) => {
          await loading.dismiss();
          await this.handleSignUpStartResponse(response);
        },
        error: (errorResponse) => {
          loading.dismiss();
          const response: Errors = errorResponse.error;
          if (response?.errors) {
            displayFieldErrors(form, response.errors)
          }
          this.messagesService.showErrorToast('Sign up failed');
        }
      });
  }

  errorMsg(username: NgModel): string | null {
    if (username.errors?.['required']) {
      return 'Username is required';
    }
    if (username.errors?.['exists']) {
      return 'Username already registered';
    }
    return null;
  }

  private async handleSignUpStartResponse(response: CredentialCreationOptionsJSON): Promise<void> {
    let credential: object | null = null;
    try {
      credential = await create(parseCreationOptionsFromJSON(response));
    } catch (e) {
      await this.messagesService.showErrorToast('Sign up failed with error ' + e);
      return;
    }
    const loading = await this.messagesService.showLoading('Finishing sign up process...');
    await loading.present();

    this.httpClient.post(`${environment.API_URL}/signup/finish`, credential)
      .subscribe({
        next: () => {
          loading.dismiss();
          this.messagesService.showSuccessToast('Sign up successful');
          this.router.navigate(['/signin']);
        },
        error: () => {
          loading.dismiss();
          this.messagesService.showErrorToast('Sign up failed');
        }
      });
  }
}


