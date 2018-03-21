import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from './dashboard/dashboard.component'
import { ApplicationsComponent } from './applications/applications.component'
import { AppdetailComponent } from './appdetail/appdetail.component'

const routes: Routes = [
  {
    path: '',
    component: DashboardComponent
  },
  {
    path: 'applications',
    component: ApplicationsComponent
  },
  {
    path: 'applications/:domain/:name',
    component: AppdetailComponent
  }

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
