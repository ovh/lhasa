import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from './components/dashboard/dashboard.component'
import { ApplicationsComponent } from './components/applications/applications.component'
import { AppdetailComponent } from './components/appdetail/appdetail.component'
import { ProfileGuard } from './guards/profile.guard';
import { RoutingGuard } from './guards/routing.guard';

const routes: Routes = [
  {
    path: '',
    component: DashboardComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'applications',
    component: ApplicationsComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'applications/:domain/:name',
    component: AppdetailComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
