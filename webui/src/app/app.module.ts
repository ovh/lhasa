import { BrowserModule } from '@angular/platform-browser';

import { NgModule } from '@angular/core';

import { ArchwizardModule } from 'ng2-archwizard';

import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './app.component';
import { ApplicationsComponent } from './applications/applications.component';
import { ApplicationsService } from './app.service';
import { DashboardComponent } from './dashboard/dashboard.component';
import { AppdetailComponent } from './appdetail/appdetail.component';
import { HttpClientModule } from '@angular/common/http';
import { MarkdownModule } from 'angular2-markdown';



@NgModule({
  declarations: [
    AppComponent,
    ApplicationsComponent,
    DashboardComponent,
    AppdetailComponent,

  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    ArchwizardModule,
    MarkdownModule.forRoot()
  ],
  providers: [ApplicationsService],
  bootstrap: [AppComponent]
})
export class AppModule { }
