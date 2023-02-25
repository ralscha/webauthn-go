import {inject, NgModule} from '@angular/core';
import {PreloadAllModules, Router, RouterModule, Routes} from '@angular/router';
import {SignInPage} from './signin/sign-in.page';
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
      return router.createUrlTree(['/signin']);
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
  {path: 'signin', component: SignInPage},
  {
    path: 'signup',
    loadChildren: () => import('./signup/sign-up.module').then(m => m.SignUpModule)
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
