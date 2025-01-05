import {PreloadAllModules, provideRouter, RouteReuseStrategy, withHashLocation, withPreloading} from '@angular/router';
import {IonicRouteStrategy} from '@ionic/angular';
import {provideHttpClient, withInterceptorsFromDi} from '@angular/common/http';
import {bootstrapApplication} from '@angular/platform-browser';
import {AppComponent} from './app/app.component';
import {provideIonicAngular} from "@ionic/angular/standalone";
import {routes} from "./app/app.routes";


bootstrapApplication(AppComponent, {
  providers: [
    provideIonicAngular(),
    {provide: RouteReuseStrategy, useClass: IonicRouteStrategy},
    provideHttpClient(withInterceptorsFromDi()),
    provideRouter(routes, withHashLocation(), withPreloading(PreloadAllModules))
  ]
})
  .catch(err => console.error(err));

