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
import { ErrorBean, ErrorsStoreService, DropErrorAction, NewErrorAction } from './stores/errors-store.service';
import { BehaviorSubject } from 'rxjs';
import { ConfigStoreService, ConfigBean } from './stores/config-store.service';

import { each, find } from 'lodash';
import { Observable } from 'rxjs/internal/Observable';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {


  /**
   * internal streams and store
   */
  protected errorsStream: Observable<ErrorBean[]>;
  public errors: ErrorBean[] = [];
  protected loadersStream: Observable<LoaderBean[]>;
  public loaders: LoaderBean[] = [];
  protected helpStream: Observable<HelpBean>;
  protected configStream: Observable<ConfigBean>;

  title = 'app';

  public items: UiKitMenuItem[];
  public sideItems: UiKitMenuItem[];
  public displaySidebar = false;
  public displayHelp = false;
  public helpContent: ContentBean;

  private TOOLBAR = [
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

  private SIDEBAR = [
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

  constructor(
    private router: Router,
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService,
    private environmentsStoreService: EnvironmentsStoreService,
    private loaderstoreService: LoadersStoreService,
    private errorsStoreService: ErrorsStoreService,
    private environmentService: DataEnvironmentService,
    private translate: TranslateService,
    private helpsStoreService: HelpsStoreService,
    private configStoreService: ConfigStoreService,
  ) {
    // this language will be used as a fallback when a translation isn't found in the current language
    this.translate.setDefaultLang('en');

    // the lang to use, if the lang isn't available, it will use the current loader to get them
    this.translate.use('en');

    // Simple menu model
    this.items = [];

    // Loaders
    this.loadersStream = this.loaderstoreService.loaders();
    // Errors
    this.errorsStream = this.errorsStoreService.errors();

    // help
    this.helpStream = this.helpsStoreService.help();
    // config
    this.configStream = this.configStoreService.help();
    // Dispath configuration loading
    this.configStoreService.request('configuration');

    // Simple menu model
    this.sideItems = [];
  }

  /**
   * handle error
   */
  public onMessageEvent(event: any) {
    if (event.data.type === 'close') {
      event.data.sender.hide();
      this.errorsStoreService.dispatch(new DropErrorAction(
        {
          code: event.data.message,
        }, null
      ));
    }
  }

  private filterByID(keys: any, values: any) {
    if (!keys) {
      return values;
    }

    const filtered = [];
    each(values, (item) => {
      const found = find(keys, (f) => {
        return f.id === item.id;
      });
      if (found) {
        filtered.push(item);
      }
    });
    return filtered;
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
    // config
    this.configStream.subscribe(
      (data: ConfigBean) => {
        if (data.config) {
          this.items = this.filterByID(data.config.toolbar, this.TOOLBAR);
          this.sideItems = this.filterByID(data.config.toolbar, this.SIDEBAR);
        }
      }
    );
    // errors
    this.errorsStream.subscribe(
      (errors: ErrorBean[]) => {
        this.errors = errors;
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
      },
      (error) => {
        this.errorsStoreService.dispatch(new NewErrorAction(
          <ErrorBean>{
            code: 'ERROR-ENVS',
            stack: JSON.stringify(error, null, 2),
          }, new BehaviorSubject<any>('select environements')
        ));
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
