import { environment } from '../environments/environment';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from './components/dashboard/dashboard.component';
import { ApplicationsComponent } from './components/applications/applications.component';
import { AppdetailComponent } from './components/appdetail/appdetail.component';
import { ProfileGuard } from './guards/profile.guard';
import { EnrollmentComponent } from './components/enrollment/enrollment.component';
import { RoutingGuard } from './guards/routing.guard';
import { DomainsComponent } from './components/domains/domains.component';
import { ApplicationResolver } from './resolver/resolve-app-detail';

const routes: Routes = [
  {
    path: '',
    component: DashboardComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'domains',
    component: DomainsComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'applications',
    component: ApplicationsComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'register/:domain/:name',
    resolve: {
      application: ApplicationResolver
    },
    component: EnrollmentComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'applications/:domain/:name',
    resolve: {
      application: ApplicationResolver
    },
    component: AppdetailComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { enableTracing: environment.tracing })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
