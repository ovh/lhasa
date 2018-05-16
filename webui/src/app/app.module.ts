import { TranslateLoader, TranslateModule, TranslateService } from '@ngx-translate/core';
import { TranslateHttpLoader } from '@ngx-translate/http-loader';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { NgModule } from '@angular/core';

import { ArchwizardModule } from 'ng2-archwizard';

import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './app.component';
import { ApplicationsComponent } from './components/applications/applications.component';
import { EnrollmentComponent } from './components/enrollment/enrollment.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { AppdetailComponent } from './components/appdetail/appdetail.component';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { MarkdownModule } from 'ngx-md';
/**
 * material2
 */
import {
  MatButtonModule,
  MatCardModule,
  MatCheckboxModule,
  MatFormFieldModule,
  MatGridListModule,
  MatInputModule,
  MatOptionModule,
  MatSelectModule,
  MatSidenavModule,
  MatSnackBarModule,
  MatTableModule,
  MatTabsModule
} from '@angular/material';
import { MatDividerModule } from '@angular/material/divider';
import { MatStepperModule } from '@angular/material/stepper';
/**
 * primeng
 */
import {
  AccordionModule,
  ButtonModule,
  CalendarModule,
  CarouselModule,
  ChartModule,
  CheckboxModule,
  CodeHighlighterModule,
  ConfirmDialogModule,
  DataGridModule,
  DataListModule,
  DataTableModule,
  DialogModule,
  DropdownModule,
  EditorModule,
  FieldsetModule,
  GrowlModule,
  InputTextareaModule,
  InputTextModule,
  MenubarModule,
  MenuModule,
  MessagesModule,
  OrderListModule,
  OverlayPanelModule,
  PanelMenuModule,
  PanelModule,
  SharedModule,
  SidebarModule,
  SliderModule,
  SpinnerModule,
  SplitButtonModule,
  StepsModule,
  TabMenuModule,
  TabViewModule,
  ToggleButtonModule,
  ToolbarModule,
  TooltipModule,
  TreeTableModule
} from 'primeng/primeng';
import { CardModule } from 'primeng/card';
import { MatSliderModule } from '@angular/material/slider';
import { TableModule } from 'primeng/table';
/**
 * guard
 */
import { ProfileGuard } from './guards/profile.guard';
import { RoutingGuard } from './guards/routing.guard';
/**
 * stores
 */
import { ApplicationsStoreService } from './stores/applications-store.service';
import { DataApplicationService } from './services/data-application-version.service';
import { DataDeploymentService } from './services/data-deployment.service';
import { StoreModule } from '@ngrx/store';
import { APP_BASE_HREF, CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ConfigurationService } from './services/configuration.service';
import { SecurityService } from './services/security.service';
import { BitbucketService } from './services/data-bitbucket.service';
import { OuiProgressTrackerComponent } from './kit/oui-progress-tracker/oui-progress-tracker.component';
import { ApplicationSortPipe, DomainSortPipe } from './pipes/pipes-applications.component';
import { environment } from '../environments/environment';
import { AppdetailsActiveDeploymentsPipe } from './pipes/pipes-appdetails.component';
import { DataEnvironmentService } from './services/data-environment.service';
import { EnvironmentsStoreService } from './stores/environments-store.service';
import { DomainsComponent } from './components/domains/domains.component';
import { GraphComponent } from './widget/graph/graph.component';
import { EnvChipComponent } from './components/env-chip/env-chip.component';
import { ApplicationResolver } from './resolver/resolve-app-detail';
import { ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { DomainsBrowseComponent } from './components/domains-browse/domains-browse.component';
import { DomainsResolver } from './resolver/resolve-domains';
import { DataDomainService } from './services/data-domain.service';
import { OuiPaginationComponent } from './kit/oui-pagination/oui-pagination.component';
import { ApplicationsResolver } from './resolver/resolve-applications';
import { OuiMessageComponent } from './kit/oui-message/oui-message.component';
import { OuiNavBarComponent } from './kit/oui-nav-bar/oui-nav-bar.component';

// Cf. https://github.com/ngx-translate/core
export function createTranslateLoader(http: HttpClient) {
  return new TranslateHttpLoader(http, './assets/i18n/', '.json');
}

@NgModule({
  declarations: [
    AppComponent,
    ApplicationsComponent,
    DomainsComponent,
    DomainsBrowseComponent,
    GraphComponent,
    EnrollmentComponent,
    DashboardComponent,
    AppdetailComponent,
    /**
     * UI Kit component
     */
    OuiProgressTrackerComponent,
    OuiNavBarComponent,
    OuiPaginationComponent,
    OuiMessageComponent,
    ApplicationSortPipe,
    DomainSortPipe,
    AppdetailsActiveDeploymentsPipe,
    EnvChipComponent,
  ],
  imports: [
    CommonModule,
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpClientModule,
    ReactiveFormsModule,
    /**
     * material2
     */
    MatSidenavModule,
    MatButtonModule,
    MatGridListModule,
    MatInputModule,
    MatCheckboxModule,
    MatTableModule,
    MatTabsModule,
    MatOptionModule,
    MatSelectModule,
    MatCardModule,
    MatSnackBarModule,
    MatFormFieldModule,
    MatDividerModule,
    MatStepperModule,
    /**
     * primeface
     */
    CardModule,
    OverlayPanelModule,
    SidebarModule,
    MatSliderModule,
    CarouselModule,
    DataTableModule,
    SharedModule,
    MenuModule,
    MenubarModule,
    CheckboxModule,
    InputTextModule,
    AccordionModule,
    CodeHighlighterModule,
    InputTextareaModule,
    DataListModule,
    TabViewModule,
    OrderListModule,
    EditorModule,
    DataGridModule,
    PanelModule,
    GrowlModule,
    MessagesModule,
    StepsModule,
    ButtonModule,
    PanelMenuModule,
    DialogModule,
    FieldsetModule,
    DropdownModule,
    ConfirmDialogModule,
    SplitButtonModule,
    ToolbarModule,
    TooltipModule,
    TreeTableModule,
    ChartModule,
    CalendarModule,
    SpinnerModule,
    SliderModule,
    ToggleButtonModule,
    TabMenuModule,
    CarouselModule,
    TableModule,
    // Local modules
    AppRoutingModule,
    ArchwizardModule,
    MarkdownModule.forRoot(),
    /**
     * store
     */
    StoreModule.forRoot({
      applications: ApplicationsStoreService.reducer,
      environments: EnvironmentsStoreService.reducer,
    }),
    /**
     * i18n
     */
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useFactory: createTranslateLoader,
        deps: [HttpClient]
      }
    })
  ],
  providers: [
    /**
     * Base Href
     */
    {
      provide: APP_BASE_HREF, useValue: environment.href
    },
    /**
     * guards
     */
    ProfileGuard,
    RoutingGuard,
    /**
     * stores
     */
    ApplicationsStoreService,
    EnvironmentsStoreService,
    /**
     * services
     */
    ConfigurationService,
    SecurityService,
    DataApplicationService,
    DataDeploymentService,
    DataEnvironmentService,
    BitbucketService,
    TranslateService,
    DataDomainService,
    /**
     * resolvers
     */
    ApplicationResolver,
    ApplicationsResolver,
    DomainsResolver
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
