import { httpResource } from '@angular/common/http';
import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../auth.service';
import { environment } from '../../environments/environment';
import { SecretOutput } from '../api/types';

@Component({
  selector: 'app-home',
  templateUrl: './home.page.html',
})
export class HomePage {
  readonly secret = httpResource<SecretOutput>(() => `${environment.API_URL}/secret`);

  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  async logout(): Promise<void> {
    this.authService
      .logout()
      .subscribe(() => this.router.navigate(['/login'], { replaceUrl: true }));
  }
}
