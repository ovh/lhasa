import { Component } from '@angular/core';
import { ApplicationsStoreService, LoadApplicationsAction } from './stores/applications-store.service';
import { DataApplicationService } from './services/data-application-version.service';
import { ApplicationBean } from './models/commons/applications-bean';
import { ContentListResponse } from './models/commons/entity-bean';
import { Router } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'app';

  constructor(
    private router: Router,
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService
  ) {

  }

  /**
   * dispatch load applications
   * @param event 
   */
  protected loadApplications(event: any) {
    // load all applications from a content return
    this.applicationsService.GetAllFromContent("", new Map<string,string>([ [ "size", "1" ]])).subscribe(
      (data: ContentListResponse<ApplicationBean>) => {
        this.applicationsStoreService.dispatch(
          new LoadApplicationsAction(
            data.content
          )
        )
      }
    );
  }
}
