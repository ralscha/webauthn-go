import {Component, inject} from '@angular/core';
import {
  IonButton,
  IonCol,
  IonContent,
  IonGrid,
  IonHeader,
  IonRouterLink,
  IonRow,
  IonTitle,
  IonToolbar,
  NavController
} from '@ionic/angular/standalone';
import {MessagesService} from '../messages.service';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {
  AuthenticationResponseJSON,
  PublicKeyCredentialRequestOptionsJSON,
  startAuthentication
} from '@simplewebauthn/browser';
import {RouterLink} from '@angular/router';

@Component({
  selector: 'app-authentication',
  templateUrl: './authentication.page.html',
  imports: [RouterLink, IonRouterLink, IonHeader, IonToolbar, IonTitle, IonContent, IonGrid, IonRow, IonCol, IonButton]
})
export class AuthenticationPage {
  readonly #navCtrl = inject(NavController);
  readonly #httpClient = inject(HttpClient);
  readonly #messagesService = inject(MessagesService);

  async login(): Promise<void> {
    const loading = await this.#messagesService.showLoading('Starting login ...');
    await loading.present();

    this.#httpClient.post<PublicKeyCredentialRequestOptionsJSON>(`${environment.API_URL}/authentication/start`, null)
      .subscribe({
        next: response => {
          loading.dismiss();
          this.handleLoginStartResponse(response);
        },
        error: () => {
          loading.dismiss();
          this.#messagesService.showErrorToast('Login failed');
        }
      });
  }

  private async handleLoginStartResponse(optionsJSON: PublicKeyCredentialRequestOptionsJSON): Promise<void> {
    let authenticationResponse: AuthenticationResponseJSON | null = null;
    try {
      authenticationResponse = await startAuthentication({optionsJSON});
    } catch (e) {
      await this.#messagesService.showErrorToast('Login failed with error ' + e);
      return;
    }
    const loading = await this.#messagesService.showLoading('Validating ...');
    await loading.present();

    this.#httpClient.post<void>(`${environment.API_URL}/authentication/finish`, authenticationResponse).subscribe({
      next: () => {
        loading.dismiss();
        this.#navCtrl.navigateRoot('/home', {replaceUrl: true});
      },
      error: () => {
        loading.dismiss();
        this.#messagesService.showErrorToast('Login failed');
      }
    });
  }
}
