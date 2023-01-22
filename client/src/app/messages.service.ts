import {Injectable} from '@angular/core';
import {LoadingController, ToastController} from '@ionic/angular';

@Injectable({
  providedIn: 'root'
})
export class MessagesService {

  constructor(private readonly toastCtrl: ToastController,
              private readonly loadingCtrl: LoadingController) {
  }

  async showLoading(message = 'Working'): Promise<HTMLIonLoadingElement> {
    const loading = await this.loadingCtrl.create({
      spinner: 'bubbles',
      message
    });
    await loading.present();
    return loading;
  }

  async showErrorToast(message = 'Unexpected error occurred'): Promise<void> {
    const toast = await this.toastCtrl.create({
      message,
      duration: 4000,
      position: 'bottom',
      color: 'danger'
    });
    await toast.present();
  }

  async showSuccessToast(message = 'Success'): Promise<void> {
    const toast = await this.toastCtrl.create({
      message,
      duration: 4000,
      position: 'bottom',
      color: 'success'
    });
    await toast.present();
  }

}
