import {Component, inject, OnInit} from '@angular/core';
import {AuthService} from '../auth.service';
import {IonButton, IonContent, IonHeader, IonTitle, IonToolbar, NavController} from '@ionic/angular/standalone';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {SecretOutput} from "../api/types";

@Component({
  selector: 'app-home',
  templateUrl: './home.page.html',
  imports: [
    IonHeader,
    IonToolbar,
    IonTitle,
    IonContent,
    IonButton
  ]
})
export class HomePage implements OnInit {
  secretMessage: string | null = null;
  private readonly authService = inject(AuthService);
  private readonly navCtrl = inject(NavController);
  private readonly httpClient = inject(HttpClient);

  async logout(): Promise<void> {
    this.authService.logout().subscribe(() => this.navCtrl.navigateRoot('/login'));
  }

  ngOnInit(): void {
    this.httpClient.get<SecretOutput>(`${environment.API_URL}/secret`)
      .subscribe(response => this.secretMessage = response.message);
  }

}
