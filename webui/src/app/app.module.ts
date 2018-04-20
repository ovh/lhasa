import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';

import {NgModule} from '@angular/core';

import {ArchwizardModule} from 'ng2-archwizard';

import {AppRoutingModule} from './app-routing.module';

import {AppComponent} from './app.component';
import {ApplicationsComponent} from './components/applications/applications.component';
import {EnrollmentComponent} from './components/enrollment/enrollment.component';
import {DashboardComponent} from './components/dashboard/dashboard.component';
import {AppdetailComponent} from './components/appdetail/appdetail.component';
import {HttpClientModule, HttpClient} from '@angular/common/http';
import {MarkdownModule} from 'ngx-md';

/**
 * material2
 */
import {MatSidenavModule} from '@angular/material';
import {MatCheckboxModule} from '@angular/material';
import {MatButtonModule} from '@angular/material';
import {MatGridListModule} from '@angular/material';
import {MatInputModule} from '@angular/material';
import {MatTableModule} from '@angular/material';
import {MatTabsModule} from '@angular/material';
import {MatSelectModule} from '@angular/material';
import {MatOptionModule} from '@angular/material';
import {MatCardModule} from '@angular/material';
import {MatSnackBarModule} from '@angular/material';
import {MatFormFieldModule} from '@angular/material';
import {MatDividerModule} from '@angular/material/divider';
import {MatStepperModule} from '@angular/material/stepper';

/**
 * primeng
 */
import {ButtonModule, OrderListModule, EditorModule, OverlayPanelModule} from 'primeng/primeng';
import {CardModule} from 'primeng/card';
import {SidebarModule} from 'primeng/primeng';
import {CarouselModule} from 'primeng/primeng';
import {ChartModule} from 'primeng/primeng';
import {DataTableModule, SharedModule} from 'primeng/primeng';
import {MenubarModule, MenuModule} from 'primeng/primeng';
import {CheckboxModule} from 'primeng/primeng';
import {InputTextModule} from 'primeng/primeng';
import {AccordionModule} from 'primeng/primeng';
import {CodeHighlighterModule} from 'primeng/primeng';
import {InputTextareaModule} from 'primeng/primeng';
import {DataListModule} from 'primeng/primeng';
import {TabViewModule} from 'primeng/primeng';
import {DataGridModule} from 'primeng/primeng';
import {PanelModule} from 'primeng/primeng';
import {GrowlModule} from 'primeng/primeng';
import {MessagesModule} from 'primeng/primeng';
import {StepsModule} from 'primeng/primeng';
import {PanelMenuModule} from 'primeng/primeng';
import {DialogModule} from 'primeng/primeng';
import {FieldsetModule} from 'primeng/primeng';
import {DropdownModule} from 'primeng/primeng';
import {ConfirmDialogModule, ConfirmationService} from 'primeng/primeng';
import {SplitButtonModule} from 'primeng/primeng';
import {ToolbarModule} from 'primeng/primeng';
import {TooltipModule} from 'primeng/primeng';
import {TreeTableModule} from 'primeng/primeng';
import {CalendarModule} from 'primeng/primeng';
import {SpinnerModule} from 'primeng/primeng';
import {SliderModule} from 'primeng/primeng';
import {MatSliderModule} from '@angular/material/slider';
import {ToggleButtonModule} from 'primeng/primeng';
import {TabMenuModule} from 'primeng/primeng';
import {TableModule} from 'primeng/table';

/**
 * guard
 */
import {ProfileGuard} from './guards/profile.guard';
import {RoutingGuard} from './guards/routing.guard';

/**
 * stores
 */
import {ApplicationsStoreService} from './stores/applications-store.service';
import {DataApplicationService} from './services/data-application-version.service';
import {DataDeploymentService} from './services/data-deployment.service';
import {StoreModule} from '@ngrx/store';
import {CommonModule, APP_BASE_HREF} from '@angular/common';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {HttpModule} from '@angular/http';
import {ConfigurationService} from './services/configuration.service';
import {SecurityService} from './services/security.service';
import {BitbucketService} from './services/data-bitbucket.service';
import {OuiProgressTrackerComponent} from './kit/oui-progress-tracker/oui-progress-tracker.component';
import {ApplicationSortPipe, DomainSortPipe} from './pipes/pipes-applications.component';
import {environment} from '../environments/environment';
import {AppdetailsActiveDeploymentsPipe} from './pipes/pipes-appdetails.component';
import {DataEnvironmentService} from './services/data-environment.service';
import {EnvironmentsStoreService} from './stores/environments-store.service';
import {DomainsComponent} from './components/domains/domains.component';
import {GraphComponent} from './widget/graph/graph.component';
import {EnvChipComponent} from './components/env-chip/env-chip.component';

@NgModule({
  declarations: [
    AppComponent,
    ApplicationsComponent,
    DomainsComponent,
    GraphComponent,
    EnrollmentComponent,
    DashboardComponent,
    AppdetailComponent,
    OuiProgressTrackerComponent,
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
    HttpModule,
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
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
