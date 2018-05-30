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
import { ButtonModule, OrderListModule, EditorModule, OverlayPanelModule } from 'primeng/primeng';
import { CardModule } from 'primeng/card';
import { SidebarModule } from 'primeng/primeng';
import { CarouselModule } from 'primeng/primeng';
import { ChartModule } from 'primeng/primeng';
import { DataTableModule, SharedModule } from 'primeng/primeng';
import { MenubarModule, MenuModule } from 'primeng/primeng';
import { CheckboxModule } from 'primeng/primeng';
import { InputTextModule } from 'primeng/primeng';
import { AccordionModule } from 'primeng/primeng';
import { CodeHighlighterModule } from 'primeng/primeng';
import { InputTextareaModule } from 'primeng/primeng';
import { DataListModule } from 'primeng/primeng';
import { TabViewModule } from 'primeng/primeng';
import { DataGridModule } from 'primeng/primeng';
import { PanelModule } from 'primeng/primeng';
import { GrowlModule } from 'primeng/primeng';
import { MessagesModule } from 'primeng/primeng';
import { StepsModule } from 'primeng/primeng';
import { PanelMenuModule } from 'primeng/primeng';
import { DialogModule } from 'primeng/primeng';
import { FieldsetModule } from 'primeng/primeng';
import { DropdownModule } from 'primeng/primeng';
import { ConfirmDialogModule, ConfirmationService } from 'primeng/primeng';
import { SplitButtonModule } from 'primeng/primeng';
import { ToolbarModule } from 'primeng/primeng';
import { TooltipModule } from 'primeng/primeng';
import { TreeTableModule } from 'primeng/primeng';
import { CalendarModule } from 'primeng/primeng';
import { SpinnerModule } from 'primeng/primeng';
import { SliderModule } from 'primeng/primeng';
import { MatSliderModule } from '@angular/material/slider';
import { ToggleButtonModule } from 'primeng/primeng';
import { TabMenuModule } from 'primeng/primeng';
import { TableModule } from 'primeng/table';
import { BlockUIModule } from 'primeng/blockui';

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
import { LoadersStoreService } from './stores/loader-store.service';
import { OuiNavBarComponent } from './kit/oui-nav-bar/oui-nav-bar.component';
import { OuiActionMenuComponent } from './kit/oui-action-menu/oui-action-menu.component';
import { DataContentService } from './services/data-content.service';
import { GraphsStoreService } from './stores/graphs-store.service';
import { GraphBrowseComponent } from './components/graphs/graph-browse.component';
import { GraphsResolver } from './resolver/resolve-graph';
import { DataGraphService } from './services/data-graph.service';

// Cf. https://github.com/ngx-translate/core
export function createTranslateLoader(http: HttpClient) {
  return new TranslateHttpLoader(http, './assets/i18n/', '.json');
}

@NgModule({
  declarations: [
    AppComponent,
    ApplicationsComponent,
    GraphBrowseComponent,
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
    OuiActionMenuComponent,
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
    BlockUIModule,
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
      graphs: GraphsStoreService.reducer,
      loaders: LoadersStoreService.reducer,
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
    LoadersStoreService,
    GraphsStoreService,
    /**
     * services
     */
    ConfigurationService,
    SecurityService,
    DataApplicationService,
    DataContentService,
    DataDeploymentService,
    DataEnvironmentService,
    BitbucketService,
    TranslateService,
    DataDomainService,
    DataGraphService,
    /**
     * resolvers
     */
    ApplicationResolver,
    ApplicationsResolver,
    DomainsResolver,
    GraphsResolver
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
