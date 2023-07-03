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
  selector: 'app-registration',
  templateUrl: './registration.page.html'
})
export class RegistrationPage {
  constructor(private readonly navCtrl: NavController,
              private readonly router: Router,
              private readonly messagesService: MessagesService,
              private readonly httpClient: HttpClient) {
  }

  async register(form: NgForm, username: string | null): Promise<void> {
    if (!username) {
      return;
    }

    const loading = await this.messagesService.showLoading('Starting registration process...');
    await loading.present();

    const userNameInput: UsernameInput = {username};

    this.httpClient.post<CredentialCreationOptionsJSON>(`${environment.API_URL}/registration/start`, userNameInput)
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
          this.messagesService.showErrorToast('Registration failed');
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
      await this.messagesService.showErrorToast('Registration failed with error ' + e);
      return;
    }
    const loading = await this.messagesService.showLoading('Finishing registration process...');
    await loading.present();

    this.httpClient.post(`${environment.API_URL}/registration/finish`, credential)
      .subscribe({
        next: () => {
          loading.dismiss();
          this.messagesService.showSuccessToast('Registration successful');
          this.router.navigate(['/login']);
        },
        error: () => {
          loading.dismiss();
          this.messagesService.showErrorToast('Registration failed');
        }
      });
  }
}


