import {Component, OnInit} from '@angular/core';
import {ApplicationsStoreService, LoadApplicationsAction} from './stores/applications-store.service';
import {DataApplicationService} from './services/data-application-version.service';
import {ApplicationBean, EnvironmentBean} from './models/commons/applications-bean';
import {ContentListResponse} from './models/commons/entity-bean';
import {Router} from '@angular/router';
import {EnvironmentsStoreService, LoadEnvironmentsAction} from './stores/environments-store.service';
import {DataEnvironmentService} from './services/data-environment.service';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'app';

  constructor(
    private router: Router,
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService,
    private environmentsStoreService: EnvironmentsStoreService,
    private environmentService: DataEnvironmentService,
    private translate: TranslateService
  ) {
      // this language will be used as a fallback when a translation isn't found in the current language
      this.translate.setDefaultLang('en');

      // the lang to use, if the lang isn't available, it will use the current loader to get them
      this.translate.use('en');
  }

  /**
   * dispatch load applications
   * @param event
   */
  public loadApplications() {
    // load all applications from a content return
    this.applicationsService.GetAllFromContent('', <Map<string, string>> {size: 1000}).subscribe(
      (data: ContentListResponse<ApplicationBean>) => {
        this.applicationsStoreService.dispatch(
          new LoadApplicationsAction(
            data.content
          )
        );
      }
    );
  }

  ngOnInit(): void {
    this.environmentService.GetAllFromContent('/', new Map<string, string>()).subscribe(
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
}
