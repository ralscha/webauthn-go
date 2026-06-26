import { Component, inject, signal } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { MessagesService } from '../messages.service';
import { environment } from '../../environments/environment';
import { Errors, UsernameInput } from '../api/types';
import { FormField, FormRoot, form, required } from '@angular/forms/signals';
import type { FieldTree, TreeValidationResult } from '@angular/forms/signals';
import { Router, RouterLink } from '@angular/router';
import { firstValueFrom } from 'rxjs';
import {
  PublicKeyCredentialCreationOptionsJSON,
  RegistrationResponseJSON,
  startRegistration,
} from '@simplewebauthn/browser';

interface RegistrationFormModel {
  username: string;
}

@Component({
  selector: 'app-registration',
  templateUrl: './registration.page.html',
  imports: [FormRoot, FormField, RouterLink],
})
export class RegistrationPage {
  readonly #router = inject(Router);
  readonly #httpClient = inject(HttpClient);
  readonly #messagesService = inject(MessagesService);

  readonly registrationModel = signal<RegistrationFormModel>({ username: '' });
  readonly registrationForm = form(
    this.registrationModel,
    (path) => {
      required(path.username, { message: 'Username is required' });
    },
    {
      submission: {
        action: (field) => this.register(field().value().username, field.username),
      },
    },
  );

  usernameError(): string | null {
    const error = this.registrationForm.username().errors()[0];
    return error?.message ?? null;
  }

  private async register(
    username: string,
    usernameField: FieldTree<string>,
  ): Promise<TreeValidationResult> {
    const loading = await this.#messagesService.showLoading('Starting registration process...');
    const userNameInput: UsernameInput = { username };

    let response: PublicKeyCredentialCreationOptionsJSON;
    try {
      response = await firstValueFrom(
        this.#httpClient.post<PublicKeyCredentialCreationOptionsJSON>(
          `${environment.API_URL}/registration/start`,
          userNameInput,
        ),
      );
    } catch (error) {
      const errors = this.extractFieldErrors(error);
      if (errors.length > 0) {
        return errors.map((kind) => ({
          kind,
          message: this.serverErrorMessage(kind),
          fieldTree: usernameField,
        }));
      }

      await this.#messagesService.showErrorToast('Registration failed');
      return { kind: 'registrationStartFailed', message: 'Registration failed' };
    } finally {
      await loading.dismiss();
    }

    await this.handleSignUpStartResponse(response);
  }

  private extractFieldErrors(error: unknown): string[] {
    if (!(error instanceof HttpErrorResponse)) {
      return [];
    }

    const response: Errors | undefined = error.error;
    return response?.errors?.['username'] ?? [];
  }

  private serverErrorMessage(kind: string): string {
    if (kind === 'exists') {
      return 'Username already registered';
    }
    return 'Registration failed';
  }

  private async handleSignUpStartResponse(
    optionsJSON: PublicKeyCredentialCreationOptionsJSON,
  ): Promise<void> {
    let registrationResponse: RegistrationResponseJSON;
    try {
      registrationResponse = await startRegistration({ optionsJSON });
    } catch (e) {
      await this.#messagesService.showErrorToast('Registration failed with error ' + e);
      return;
    }
    const loading = await this.#messagesService.showLoading('Finishing registration process...');

    try {
      await firstValueFrom(
        this.#httpClient.post(`${environment.API_URL}/registration/finish`, registrationResponse),
      );
      await this.#messagesService.showSuccessToast('Registration successful');
      await this.#router.navigate(['/login']);
    } catch {
      await this.#messagesService.showErrorToast('Registration failed');
    } finally {
      await loading.dismiss();
    }
  }
}
