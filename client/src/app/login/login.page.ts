import {Component} from '@angular/core';
import {AuthService} from '../auth.service';
import {LoadingController, NavController} from '@ionic/angular';
import {MessagesService} from '../messages.service';
import {HttpClient} from '@angular/common/http';
import {CredentialRequestOptionsJSON, get, parseRequestOptionsFromJSON,} from "@github/webauthn-json/browser-ponyfill";
import {environment} from '../../environments/environment';

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html'
})
export class LoginPage {

  constructor(private readonly authService: AuthService,
              private readonly loadingCtrl: LoadingController,
              private readonly navCtrl: NavController,
              private readonly httpClient: HttpClient,
              private readonly messagesService: MessagesService) {
  }

  async login(): Promise<void> {
    const loading = await this.messagesService.showLoading('Starting login ...');
    await loading.present();

    this.httpClient.post<CredentialRequestOptionsJSON>(`${environment.API_URL}/login/start`, null)
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

  private async handleLoginStartResponse(response: CredentialRequestOptionsJSON): Promise<void> {
    let credential: object | null = null;
    try {
      credential = await get(parseRequestOptionsFromJSON(response));
    } catch (e) {
      await this.messagesService.showErrorToast('Login failed with error ' + e);
      return;
    }
    const loading = await this.messagesService.showLoading('Validating ...');
    await loading.present();

    this.httpClient.post<void>(`${environment.API_URL}/login/finish`, credential).subscribe({
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
