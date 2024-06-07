import {Component} from '@angular/core';
import {NavController} from '@ionic/angular';
import {MessagesService} from '../messages.service';
import { HttpClient } from '@angular/common/http';
import {environment} from '../../environments/environment';
import { startAuthentication } from '@simplewebauthn/browser';
import {AuthenticationResponseJSON, PublicKeyCredentialRequestOptionsJSON} from "@simplewebauthn/types";
@Component({
  selector: 'app-authentication',
  templateUrl: './authentication.page.html'
})
export class AuthenticationPage {

  constructor(private readonly navCtrl: NavController,
              private readonly httpClient: HttpClient,
              private readonly messagesService: MessagesService) {
  }

  async login(): Promise<void> {
    const loading = await this.messagesService.showLoading('Starting login ...');
    await loading.present();

    this.httpClient.post<PublicKeyCredentialRequestOptionsJSON>(`${environment.API_URL}/authentication/start`, null)
      .subscribe({
        next: response => {
          loading.dismiss();
          this.handleLoginStartResponse(response);
        },
        error: () => {
          loading.dismiss();
          this.messagesService.showErrorToast('Login failed');
        }
      });
  }

  private async handleLoginStartResponse(response: PublicKeyCredentialRequestOptionsJSON): Promise<void> {
    let authenticationResponse: AuthenticationResponseJSON | null = null;
    try {
      authenticationResponse = await startAuthentication(response);
    } catch (e) {
      await this.messagesService.showErrorToast('Login failed with error ' + e);
      return;
    }
    const loading = await this.messagesService.showLoading('Validating ...');
    await loading.present();

    this.httpClient.post<void>(`${environment.API_URL}/authentication/finish`, authenticationResponse).subscribe({
      next: () => {
        loading.dismiss();
        this.navCtrl.navigateRoot('/home', {replaceUrl: true});
      },
      error: () => {
        loading.dismiss();
        this.messagesService.showErrorToast('Login failed');
      }
    });
  }
}
