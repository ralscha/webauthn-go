import { HttpClient } from '@angular/common/http';
import { Component, inject, OnDestroy, OnInit, signal } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { firstValueFrom } from 'rxjs';
import {
  AuthenticationResponseJSON,
  browserSupportsWebAuthnAutofill,
  PublicKeyCredentialRequestOptionsJSON,
  startAuthentication,
} from '@simplewebauthn/browser';
import { MessagesService } from '../messages.service';
import { environment } from '../../environments/environment';

@Component({
  selector: 'app-authentication',
  templateUrl: './authentication.page.html',
  imports: [RouterLink],
})
export class AuthenticationPage implements OnInit, OnDestroy {
  readonly conditionalMediationAvailable = signal(false);

  readonly #router = inject(Router);
  readonly #httpClient = inject(HttpClient);
  readonly #messagesService = inject(MessagesService);
  #active = true;

  ngOnInit(): void {
    void this.startPasskeyAutofill();
  }

  ngOnDestroy(): void {
    this.#active = false;
  }

  async login(): Promise<void> {
    const loading = await this.#messagesService.showLoading('Starting login ...');
    try {
      const response = await firstValueFrom(
        this.#httpClient.post<PublicKeyCredentialRequestOptionsJSON>(
          `${environment.API_URL}/authentication/start`,
          null,
        ),
      );
      await loading.dismiss();
      await this.handleLoginStartResponse(response);
    } catch {
      await loading.dismiss();
      await this.#messagesService.showErrorToast('Login failed');
    }
  }

  private async handleLoginStartResponse(
    optionsJSON: PublicKeyCredentialRequestOptionsJSON,
  ): Promise<void> {
    try {
      const credential = await startAuthentication({ optionsJSON });
      await this.finishAuthentication(credential);
    } catch (error) {
      if (!this.isExpectedCredentialError(error)) {
        await this.#messagesService.showErrorToast('Login failed');
      }
    }
  }

  private async startPasskeyAutofill(): Promise<void> {
    try {
      this.conditionalMediationAvailable.set(await browserSupportsWebAuthnAutofill());

      if (!this.conditionalMediationAvailable()) {
        return;
      }

      const response = await firstValueFrom(
        this.#httpClient.post<PublicKeyCredentialRequestOptionsJSON>(
          `${environment.API_URL}/authentication/start`,
          null,
        ),
      );

      const credential = await startAuthentication({
        optionsJSON: response,
        useBrowserAutofill: true,
      });

      await this.finishAuthentication(credential);
    } catch (error) {
      if (!this.isExpectedCredentialError(error)) {
        await this.#messagesService.showErrorToast('Passkey autofill failed');
      }
    }
  }

  private async finishAuthentication(credential: AuthenticationResponseJSON): Promise<void> {
    if (!this.#active) {
      return;
    }

    const loading = await this.#messagesService.showLoading('Validating ...');

    try {
      await firstValueFrom(
        this.#httpClient.post<void>(`${environment.API_URL}/authentication/finish`, credential),
      );
      await loading.dismiss();
      await this.#router.navigate(['/home'], { replaceUrl: true });
    } catch {
      await loading.dismiss();
      await this.#messagesService.showErrorToast('Login failed');
    }
  }

  private isExpectedCredentialError(error: unknown): boolean {
    return (
      error instanceof DOMException &&
      (error.name === 'AbortError' || error.name === 'NotAllowedError')
    );
  }
}
