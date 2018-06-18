import { HelpsStoreService, HelpBean } from './stores/help-store.service';
import { Component, OnInit } from '@angular/core';
import { ApplicationsStoreService, LoadApplicationsAction } from './stores/applications-store.service';
import { DataApplicationService } from './services/data-application-version.service';
import { ApplicationBean, EnvironmentBean, ApplicationPagesBean } from './models/commons/applications-bean';
import { ContentListResponse } from './models/commons/entity-bean';
import { Router } from '@angular/router';
import { EnvironmentsStoreService, LoadEnvironmentsAction } from './stores/environments-store.service';
import { DataEnvironmentService } from './services/data-environment.service';
import { TranslateService } from '@ngx-translate/core';
import { UiKitMenuItem } from './models/kit/navbar';
import { Store } from '@ngrx/store';
import { LoaderBean, LoadersStoreService } from './stores/loader-store.service';
import { SidebarModule } from 'primeng/sidebar';
import { ContentBean } from './models/commons/content-bean';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {


  /**
   * internal streams and store
   */
  protected loadersStream: Store<LoaderBean[]>;
  public loaders: LoaderBean[] = [];
  protected helpStream: Store<HelpBean>;

  title = 'app';

  public items: UiKitMenuItem[];
  public sideItems: UiKitMenuItem[];
  public displaySidebar = false;
  public displayHelp = false;
  public helpContent: ContentBean;

  constructor(
    private router: Router,
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService,
    private environmentsStoreService: EnvironmentsStoreService,
    private loaderstoreService: LoadersStoreService,
    private environmentService: DataEnvironmentService,
    private translate: TranslateService,
    private helpsStoreService: HelpsStoreService,
  ) {
    // this language will be used as a fallback when a translation isn't found in the current language
    this.translate.setDefaultLang('en');

    // the lang to use, if the lang isn't available, it will use the current loader to get them
    this.translate.use('en');

    // Simple menu model
    this.items = [
      {
        id: 'domains',
        label: 'DOMAINS',
        routerLink: '/domains'
      },
      {
        id: 'applications',
        label: 'APPLICATIONS',
        routerLink: '/applications'
      },
      {
        id: 'badges',
        label: 'BADGES',
        routerLink: '/badges'
      },
      {
        id: 'maps',
        label: 'MAPS',
        routerLink: '/graph/deployments'
      }
    ];

    // Loaders
    this.loadersStream = this.loaderstoreService.loaders();

    // help
    this.helpStream = this.helpsStoreService.help();

    // Simple menu model
    this.sideItems = [
      {
        id: 'domains',
        label: 'DOMAINS',
        items: [
          {
            id: 'domains-all',
            label: 'DOMAINS-ALL',
            routerLink: '/domains'
          }
        ]
      },
      {
        id: 'applications',
        label: 'APPLICATIONS',
        items: [
          {
            id: 'applications-all',
            label: 'APPLICATIONS-ALL',
            routerLink: '/applications'
          }
        ]
      },
      {
        id: 'badges',
        label: 'BADGES',
        items: [
          {
            id: 'badges-all',
            label: 'BADGES-ALL',
            routerLink: '/badges'
          }
        ]
      },
      {
        id: 'maps',
        label: 'MAPS',
        items: [
          {
            id: 'deployments-graph',
            label: 'MAPS-ALL',
            routerLink: '/graph/deployments'
          }
        ]
      },
    ];
  }

  ngOnInit(): void {
    // loaders
    this.loadersStream.subscribe(
      (loaders: LoaderBean[]) => {
        this.loaders = loaders;
      }
    );
    // help
    this.helpStream.subscribe(
      (help: HelpBean) => {
        if (help.token) {
          this.displayHelp = true;
          this.helpContent = help.content;
        }
      }
    );
    // read domains
    this.environmentService.GetAllFromContent('/', null).subscribe(
      (value: ContentListResponse<EnvironmentBean>) => {
        const environmentMap = new Map<string, EnvironmentBean>();
        value.content.forEach(env => {
          environmentMap[env.slug] = env;
        });
        this.environmentsStoreService.dispatch(
          new LoadEnvironmentsAction(
            environmentMap
          )
        );
      }
    );
  }

  public handler(event) {
    switch (event.data) {
      case '/':
        this.displaySidebar = true;
        break;
      default:
        break;
    }
  }
}
