import { TranslateLoader, TranslateModule, TranslateService } from '@ngx-translate/core';
import { TranslateHttpLoader } from '@ngx-translate/http-loader';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { NgModule } from '@angular/core';

import { ArchwizardModule } from 'ng2-archwizard';

import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './app.component';
import { ApplicationsComponent } from './components/applications/applications.component';
import { AppEditComponent } from './components/appedit/appedit.component';
import { AppdetailComponent } from './components/appdetail/appdetail.component';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { NgxMdModule } from 'ngx-md';
import { OpenAPIUIComponent } from './components/openapi-ui/openapi-ui.component';
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
import { DataBadgeRatingsService } from './services/data-badgeratings.service';
import { StoreModule } from '@ngrx/store';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ConfigurationService } from './services/configuration.service';
import { SecurityService } from './services/security.service';
import { OuiProgressTrackerComponent } from './kit/oui-progress-tracker/oui-progress-tracker.component';
import { ApplicationSortPipe, DomainSortPipe } from './pipes/pipes-applications.component';
import { DataEnvironmentService } from './services/data-environment.service';
import { EnvironmentsStoreService } from './stores/environments-store.service';
import { GraphComponent } from './widget/graph/graph.component';
import { BadgeWidgetComponent } from './widget/badgewidget/badgewidget.component';
import { EnvChipComponent } from './components/env-chip/env-chip.component';
import { ApplicationResolver } from './resolver/resolve-app-detail';
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
import { environment } from '../environments/environment';
import { BadgesComponent } from './components/badges/badges.component';
import { BadgesResolver } from './resolver/resolve-badges';
import { DataBadgeService } from './services/data-badge.service';
import { DataBadgeStatsService } from './services/data-badgestats.service';
import { BadgesStoreService } from './stores/badges-store.service';
import { HelpsStoreService } from './stores/help-store.service';
import { ConfigStoreService } from './stores/config-store.service';
import { ErrorsStoreService } from './stores/errors-store.service';
import { HelpWidgetComponent } from './widget/help-widget/help-widget.component';

export function getBaseUrl() {
  return document.getElementsByTagName('base')[0].href;
}

// Cf. https://github.com/ngx-translate/core
export function createTranslateLoader(http: HttpClient) {
  return new TranslateHttpLoader(http, getBaseUrl() + '/assets/i18n/', '.json');
}

@NgModule({
  declarations: [
    AppComponent,
    ApplicationsComponent,
    GraphBrowseComponent,
    DomainsBrowseComponent,
    GraphComponent,
    BadgeWidgetComponent,
    AppEditComponent,
    AppdetailComponent,
    BadgesComponent,
    OpenAPIUIComponent,
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
    EnvChipComponent,
    HelpWidgetComponent,
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
    NgxMdModule.forRoot(),
    /**
     * store
     */
    StoreModule.forRoot({
      applications: ApplicationsStoreService.reducer,
      environments: EnvironmentsStoreService.reducer,
      graphs: GraphsStoreService.reducer,
      badges: BadgesStoreService.reducer,
      loaders: LoadersStoreService.reducer,
      helps: HelpsStoreService.reducer,
      config: ConfigStoreService.reducer,
      errors: ErrorsStoreService.reducer,
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
     * guards
     */
    ProfileGuard,
    RoutingGuard,
    /**
     * stores
     */
    ApplicationsStoreService,
    BadgesStoreService,
    EnvironmentsStoreService,
    LoadersStoreService,
    GraphsStoreService,
    HelpsStoreService,
    ConfigStoreService,
    ErrorsStoreService,
    /**
     * services
     */
    ConfigurationService,
    SecurityService,
    DataApplicationService,
    DataContentService,
    DataDeploymentService,
    DataBadgeRatingsService,
    DataEnvironmentService,
    TranslateService,
    DataDomainService,
    DataBadgeService,
    DataBadgeStatsService,
    DataGraphService,
    /**
     * resolvers
     */
    ApplicationResolver,
    ApplicationsResolver,
    DomainsResolver,
    BadgesResolver,
    GraphsResolver
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
