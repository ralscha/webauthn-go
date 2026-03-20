import {HttpClient} from '@angular/common/http';
import {Component, inject, OnDestroy, OnInit} from '@angular/core';
import {RouterLink} from '@angular/router';
import {firstValueFrom} from 'rxjs';
import {
  IonButton,
  IonCol,
  IonContent,
  IonGrid,
  IonHeader,
  IonInput,
  IonItem,
  IonRouterLink,
  IonRow,
  IonTitle,
  IonToolbar,
  NavController
} from '@ionic/angular/standalone';
import {MessagesService} from '../messages.service';
import {environment} from '../../environments/environment';

@Component({
  selector: 'app-authentication',
  templateUrl: './authentication.page.html',
  imports: [RouterLink, IonRouterLink, IonHeader, IonToolbar, IonTitle, IonContent, IonGrid, IonRow, IonCol, IonItem, IonInput, IonButton]
})
export class AuthenticationPage implements OnInit, OnDestroy {
  conditionalMediationAvailable = false;

  readonly #navCtrl = inject(NavController);
  readonly #httpClient = inject(HttpClient);
  readonly #messagesService = inject(MessagesService);
  #conditionalMediationAbortController: AbortController | null = null;

  ngOnInit(): void {
    void this.startPasskeyAutofill();
  }

  ngOnDestroy(): void {
    this.abortConditionalMediation();
  }

  async login(): Promise<void> {
    this.abortConditionalMediation();

    const loading = await this.#messagesService.showLoading('Starting login ...');
    try {
      const response = await firstValueFrom(
        this.#httpClient.post<PublicKeyCredentialRequestOptionsJSON>(`${environment.API_URL}/authentication/start`, null)
      );
      await loading.dismiss();
      await this.handleLoginStartResponse(response);
    }
    catch {
      await loading.dismiss();
      await this.#messagesService.showErrorToast('Login failed');
    }
  }

  private async handleLoginStartResponse(optionsJSON: PublicKeyCredentialRequestOptionsJSON): Promise<void> {
    try {
      const publicKey = PublicKeyCredential.parseRequestOptionsFromJSON(optionsJSON);
      const credential = await navigator.credentials.get({publicKey}) as PublicKeyCredential | null;

      if (!credential) {
        return;
      }

      await this.finishAuthentication(credential.toJSON());
    }
    catch (error) {
      if (!this.isExpectedCredentialError(error)) {
        await this.#messagesService.showErrorToast('Login failed');
      }
    }
  }

  private async startPasskeyAutofill(): Promise<void> {
    if (!window.PublicKeyCredential
      || typeof PublicKeyCredential.isConditionalMediationAvailable !== 'function') {
      return;
    }

    try {
      this.conditionalMediationAvailable =
        await PublicKeyCredential.isConditionalMediationAvailable();

      if (!this.conditionalMediationAvailable) {
        return;
      }

      const response = await firstValueFrom(
        this.#httpClient.post<PublicKeyCredentialRequestOptionsJSON>(`${environment.API_URL}/authentication/start`, null)
      );
      const publicKey = PublicKeyCredential.parseRequestOptionsFromJSON(response);

      this.#conditionalMediationAbortController = new AbortController();
      const credential = await navigator.credentials.get({
        publicKey,
        mediation: 'conditional',
        signal: this.#conditionalMediationAbortController.signal
      }) as PublicKeyCredential | null;
      this.#conditionalMediationAbortController = null;

      if (!credential) {
        return;
      }

      await this.finishAuthentication(credential.toJSON());
    }
    catch (error) {
      this.#conditionalMediationAbortController = null;
      if (!this.isExpectedCredentialError(error)) {
        await this.#messagesService.showErrorToast('Passkey autofill failed');
      }
    }
  }

  private async finishAuthentication(credential: PublicKeyCredentialJSON): Promise<void> {
    const loading = await this.#messagesService.showLoading('Validating ...');

    try {
      await firstValueFrom(this.#httpClient.post<void>(`${environment.API_URL}/authentication/finish`, credential));
      await loading.dismiss();
      await this.#navCtrl.navigateRoot('/home', {replaceUrl: true});
    }
    catch {
      await loading.dismiss();
      await this.#messagesService.showErrorToast('Login failed');
    }
  }

  private abortConditionalMediation(): void {
    this.#conditionalMediationAbortController?.abort();
    this.#conditionalMediationAbortController = null;
  }

  private isExpectedCredentialError(error: unknown): boolean {
    return error instanceof DOMException
      && (error.name === 'AbortError' || error.name === 'NotAllowedError');
  }
}
