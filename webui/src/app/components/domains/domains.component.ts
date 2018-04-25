import { Component, OnInit, Pipe } from '@angular/core';
import { ApplicationsStoreService, SelectApplicationAction, LoadApplicationsAction } from '../../stores/applications-store.service';
import { Store } from '@ngrx/store';
import { ApplicationBean, DomainBean } from '../../models/commons/applications-bean';
import { DataApplicationService } from '../../services/data-application-version.service';
import { ContentListResponse } from '../../models/commons/entity-bean';

import * as _ from 'lodash';

@Component({
  selector: 'app-domains',
  templateUrl: './domains.component.html',
  styleUrls: ['./domains.component.css'],
})
export class DomainsComponent implements OnInit {

  /**
   * internal streams and store
   */
  protected domainsStream: Store<DomainBean[]>;
  protected domains: DomainBean[] = [];

  constructor(
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService
  ) {
    /**
     * subscribe
     */
    this.domainsStream = this.applicationsStoreService.domains();

    this.domainsStream.subscribe(
      (domains: DomainBean[]) => {
        this.domains = domains
        _.each(domains, (domain) => {
          let applications = domain.applications;
          _.each(applications, (app) => {
            if (app.manifest && app.manifest.description) {
              app.description = (app.manifest.description.length > 200) ? (app.manifest.description.substr(0, 200) + "...") : (app.manifest.description)
            } else {
              app.description = "No description ..."
              if (!app.manifest) {
                app.manifest = {}
              }
            }
          });
        });
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }

  ngOnInit() {
    if (!this.domains || this.domains.length == 0) {
      this.loadApplications(null)
    }
  }

  /**
 * dispatch load applications
 * @param event 
 */
  protected loadApplications(event: any) {
    // load all applications from a content return
    this.applicationsService.GetAllFromContent("", <Map<string, string>>{ size: 1000 }).subscribe(
      (data: ContentListResponse<ApplicationBean>) => {
        this.applicationsStoreService.dispatch(
          new LoadApplicationsAction(
            data.content
          )
        )
      }
    );
  }

  /**
   * dispatch load applications
   * @param event 
   */
  protected selectApplication(application: ApplicationBean) {
    // load all applications from a content return
    this.applicationsService.GetAllFromContent('/' + application.domain + '/' + application.name, <Map<string, string>>{ size: 1 }).subscribe(
      (data: ContentListResponse<ApplicationBean>) => {
        this.applicationsStoreService.dispatch(
          new SelectApplicationAction(
            data.content[0]
          )
        )
      }
    );
  }
}
