import { GraphBrowseComponent } from './components/graphs/graph-browse.component';
import { environment } from '../environments/environment';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from './components/dashboard/dashboard.component';
import { ApplicationsComponent } from './components/applications/applications.component';
import { AppdetailComponent } from './components/appdetail/appdetail.component';
import { ProfileGuard } from './guards/profile.guard';
import { AppEditComponent } from './components/appedit/appedit.component';
import { RoutingGuard } from './guards/routing.guard';
import { ApplicationResolver } from './resolver/resolve-app-detail';
import { DomainsBrowseComponent } from './components/domains-browse/domains-browse.component';
import { DomainsResolver } from './resolver/resolve-domains';
import { ApplicationsResolver } from './resolver/resolve-applications';
import { GraphsResolver } from './resolver/resolve-graph';
import { BadgesComponent } from './components/badges/badges.component';
import { BadgesResolver } from './resolver/resolve-badges';


const routes: Routes = [
  {
    path: '',
    redirectTo: 'applications', 
    pathMatch: 'full',
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'graph/deployments',
    resolve: {
      graph: GraphsResolver
    },
    component: GraphBrowseComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'applications',
    resolve: {
      applications: ApplicationsResolver
    },
    component: ApplicationsComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'register/:domain/:name',
    resolve: {
      application: ApplicationResolver
    },
    component: AppEditComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'applications/:domain/:name/versions/:version',
    resolve: {
      application: ApplicationResolver
    },
    component: AppdetailComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'applications/:domain/:name/versions/:version/edit',
    resolve: {
      application: ApplicationResolver
    },
    component: AppEditComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'domains',
    resolve: {
      domains: DomainsResolver
    },
    component: DomainsBrowseComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
  {
    path: 'badges',
    resolve: {
      badges: BadgesResolver
    },
    component: BadgesComponent,
    canActivate: [ProfileGuard, RoutingGuard]
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { enableTracing: environment.tracing })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
