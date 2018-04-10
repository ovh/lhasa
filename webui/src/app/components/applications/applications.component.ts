import { Component, OnInit, Pipe } from '@angular/core';
import { ApplicationsStoreService, SelectApplicationAction, LoadApplicationsAction } from '../../stores/applications-store.service';
import { Store } from '@ngrx/store';
import { ApplicationBean, DomainBean } from '../../models/commons/applications-bean';
import { DataApplicationService } from '../../services/data-application-version.service';
import { ContentListResponse } from '../../models/commons/entity-bean';

import * as _ from 'lodash';

@Component({
  selector: 'app-applications',
  templateUrl: './applications.component.html',
  styleUrls: ['./applications.component.css'],
})
export class ApplicationsComponent implements OnInit {

  /**
   * internal streams and store
   */
  protected applicationsStream: Store<ApplicationBean[]>;
  protected applications: ApplicationBean[];
  protected orderedDomains = new Map<string, ApplicationBean[]>()
  protected domains: DomainBean[] = [];

  constructor(
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService
  ) {
    /**
     * subscribe
     */
    this.applicationsStream = this.applicationsStoreService.applications();

    this.applicationsStream.subscribe(
      (element: ApplicationBean[]) => {
        this.applications = element;
        _.each(this.applications, (app) => {
          if(app.manifest) {
            app.description = (app.manifest.description.length>200)? (app.manifest.description.substr(0,200)+"...") : (app.manifest.description)
          } else {
            app.description = "No description ..."
          }
        });
        this.orderedDomains = new Map<string, ApplicationBean[]>()
        _.each(this.applications, (app) => {
          if(!this.orderedDomains.has(app.domain)) {
            this.orderedDomains.set(app.domain, [])
          }
          this.orderedDomains.get(app.domain).push(app)
        });
        this.domains = []
        this.orderedDomains.forEach((v,k) => {
          this.domains.push({name: k, applications: v})
        })
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }

  ngOnInit() {
    if(!this.applications || this.applications.length == 0) {
      this.loadApplications(null)
    }
  }

    /**
   * dispatch load applications
   * @param event 
   */
  protected loadApplications(event: any) {
    // load all applications from a content return
    this.applicationsService.GetAllFromContent("", <Map<string,string>> {size: 1000}).subscribe(
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
    this.applicationsService.GetAllFromContent('/' + application.domain + '/' + application.name, <Map<string,string>> {size: 1}).subscribe(
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
