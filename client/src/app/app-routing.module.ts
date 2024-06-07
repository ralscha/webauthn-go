import {inject, NgModule} from '@angular/core';
import {PreloadAllModules, Router, RouterModule, Routes} from '@angular/router';
import {AuthenticationPage} from './authentication/authentication.page';
import {AuthService} from "./auth.service";
import {map} from "rxjs/operators";

export const authGuard = (authService = inject(AuthService), router = inject(Router)) => {
  if (authService.isLoggedIn()) {
    return true;
  }

  return authService.isAuthenticated().pipe(
    map(success => {
      if (success) {
        return true;
      }
      return router.createUrlTree(['/login']);
    })
  );
}

const routes: Routes = [
  {path: '', redirectTo: 'home', pathMatch: 'full'},
  {
    path: 'home',
    canActivate: [() => authGuard()],
    loadChildren: () => import('./home/home.module').then(m => m.HomePageModule)
  },
  {path: 'login', component: AuthenticationPage},
  {
    path: 'registration',
    loadChildren: () => import('./registration/registration.module').then(m => m.RegistrationModule)
  }
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, {preloadingStrategy: PreloadAllModules, useHash: true})
  ],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
