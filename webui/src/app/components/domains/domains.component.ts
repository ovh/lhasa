import { Component, OnInit, Pipe } from '@angular/core';
import { ApplicationsStoreService, SelectApplicationAction, LoadApplicationsAction } from '../../stores/applications-store.service';
import { Store } from '@ngrx/store';
import { ApplicationBean, DomainBean } from '../../models/commons/applications-bean';
import { DataApplicationService } from '../../services/data-application-version.service';
import { ContentListResponse } from '../../models/commons/entity-bean';

import * as _ from 'lodash';
import { Node, Edge } from '../../models/graph/graph-bean';

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

  /**
   * to graph domains
   */
  protected nodes: Node[];
  protected edges: Edge[];

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
        this.nodes = []
        this.edges = []
        this.domains = domains
        let indexApp = new Map<string, number>()
        _.each(domains, (domain) => {
          this.nodes.push({
            id: domain.name,
            label: domain.name
          })
          let applications = domain.applications;
          _.each(applications, (app) => {
            let keyApp = domain.name + "#" + app.name
            if(indexApp.has(keyApp)) {
              let instance = indexApp.get(keyApp)+1
              indexApp.set(keyApp,instance)
            } else {
              indexApp.set(keyApp,0)
            }
            this.nodes.push({
              id: keyApp + "#" + indexApp.get(keyApp),
              label: app.name + "#" + indexApp.get(keyApp)
            })
            this.edges.push({
              from: domain.name,
              to: keyApp + "#" + indexApp.get(keyApp),
              label: domain.name + " to " + app.name + "#" + indexApp.get(keyApp)
            })
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
