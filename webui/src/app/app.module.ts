import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { NgModule } from '@angular/core';

import { ArchwizardModule } from 'ng2-archwizard';

import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './app.component';
import { ApplicationsComponent } from './components/applications/applications.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { AppdetailComponent } from './components/appdetail/appdetail.component';
import { HttpClientModule, HttpClient } from '@angular/common/http';
import { MarkdownModule } from 'angular2-markdown';

/**
 * material2
 */
import { MatSidenavModule } from '@angular/material';
import { MatCheckboxModule } from '@angular/material';
import { MatButtonModule } from '@angular/material';
import { MatGridListModule } from '@angular/material';
import { MatInputModule } from '@angular/material';
import { MatTableModule } from '@angular/material';
import { MatTabsModule } from '@angular/material';
import { MatSelectModule } from '@angular/material';
import { MatOptionModule } from '@angular/material';
import { MatCardModule } from '@angular/material';
import { MatSnackBarModule } from '@angular/material';
import { MatFormFieldModule } from '@angular/material';
import { MatDividerModule } from '@angular/material/divider'
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

/**
 * guard
 */
import { ProfileGuard } from './guards/profile.guard';
import { RoutingGuard } from './guards/routing.guard';

/**
 * stores
 */
import { ApplicationsStoreService } from './stores/applications-store.service';
import { DataApplicationServiceService } from './services/data-application-version.service';
import { StoreModule } from '@ngrx/store';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { ConfigurationService } from './services/configuration.service';
import { SecurityService } from './services/security.service';

@NgModule({
  declarations: [
    AppComponent,
    ApplicationsComponent,
    DashboardComponent,
    AppdetailComponent
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
    // Local modules
    AppRoutingModule,
    ArchwizardModule,
    MarkdownModule.forRoot(),
    /**
     * store
     */
    StoreModule.forRoot({
      applications: ApplicationsStoreService.reducer
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
    /**
     * services
     */
    ConfigurationService,
    SecurityService,
    DataApplicationServiceService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
