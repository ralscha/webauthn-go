import {Component, inject} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {MessagesService} from '../messages.service';
import {environment} from '../../environments/environment';
import {Errors, UsernameInput} from "../api/types";
import {FormsModule, NgForm, NgModel} from "@angular/forms";
import {displayFieldErrors} from "../util";
import {Router} from "@angular/router";
import {
  PublicKeyCredentialCreationOptionsJSON,
  RegistrationResponseJSON,
  startRegistration,
} from '@simplewebauthn/browser';
import {
  IonBackButton,
  IonButton,
  IonButtons,
  IonCol,
  IonContent,
  IonGrid,
  IonHeader,
  IonInput,
  IonItem,
  IonRow,
  IonTitle,
  IonToolbar
} from "@ionic/angular/standalone";

@Component({
  selector: 'app-registration',
  templateUrl: './registration.page.html',
  imports: [FormsModule, IonHeader, IonToolbar, IonTitle, IonContent, IonGrid, IonRow, IonCol, IonButton, IonButtons, IonBackButton, IonInput, IonItem]
})
export class RegistrationPage {
  readonly #router = inject(Router);
  readonly #httpClient = inject(HttpClient);
  readonly #messagesService = inject(MessagesService);

  async register(form: NgForm, username: string | null): Promise<void> {
    if (!username) {
      return;
    }

    const loading = await this.#messagesService.showLoading('Starting registration process...');
    await loading.present();

    const userNameInput: UsernameInput = {username};

    this.#httpClient.post<PublicKeyCredentialCreationOptionsJSON>(`${environment.API_URL}/registration/start`, userNameInput)
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
          this.#messagesService.showErrorToast('Registration failed');
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

  private async handleSignUpStartResponse(optionsJSON: PublicKeyCredentialCreationOptionsJSON): Promise<void> {
    let registrationResponse: RegistrationResponseJSON | null = null;
    try {
      registrationResponse = await startRegistration({optionsJSON})
    } catch (e) {
      await this.#messagesService.showErrorToast('Registration failed with error ' + e);
      return;
    }
    const loading = await this.#messagesService.showLoading('Finishing registration process...');
    await loading.present();

    this.#httpClient.post(`${environment.API_URL}/registration/finish`, registrationResponse)
      .subscribe({
        next: () => {
          loading.dismiss();
          this.#messagesService.showSuccessToast('Registration successful');
          this.#router.navigate(['/login']);
        },
        error: () => {
          loading.dismiss();
          this.#messagesService.showErrorToast('Registration failed');
        }
      });
  }
}


