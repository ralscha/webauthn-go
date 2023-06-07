import {Component} from '@angular/core';
import {AuthService} from '../auth.service';
import {LoadingController, NavController} from '@ionic/angular';
import {MessagesService} from '../messages.service';
import {HttpClient} from '@angular/common/http';
import {CredentialRequestOptionsJSON, get, parseRequestOptionsFromJSON,} from "@github/webauthn-json/browser-ponyfill";
import {environment} from '../../environments/environment';

@Component({
  selector: 'app-sign-in',
  templateUrl: './sign-in.page.html'
})
export class SignInPage {

  constructor(private readonly authService: AuthService,
              private readonly loadingCtrl: LoadingController,
              private readonly navCtrl: NavController,
              private readonly httpClient: HttpClient,
              private readonly messagesService: MessagesService) {
  }

  async signIn(): Promise<void> {
    const loading = await this.messagesService.showLoading('Starting sign in ...');
    await loading.present();

    this.httpClient.post<CredentialRequestOptionsJSON>(`${environment.API_URL}/signin/start`, null)
      .subscribe({
        next: response => {
          loading.dismiss();
          this.handleSignInStartResponse(response);
        },
        error: () => {
          loading.dismiss();
          this.messagesService.showErrorToast('Sign in failed');
        }
      });
  }

  private async handleSignInStartResponse(response: CredentialRequestOptionsJSON): Promise<void> {
    let credential: object | null = null;
    try {
      credential = await get(parseRequestOptionsFromJSON(response));
    } catch (e) {
      await this.messagesService.showErrorToast('Sign in failed with error ' + e);
      return;
    }
    const loading = await this.messagesService.showLoading('Validating ...');
    await loading.present();

    this.httpClient.post<void>(`${environment.API_URL}/signin/finish`, credential).subscribe({
      next: () => {
        loading.dismiss();
        this.navCtrl.navigateRoot('/home', {replaceUrl: true});
      },
      error: () => {
        loading.dismiss();
        this.messagesService.showErrorToast('Sign in failed');
      }
    });
  }
}
