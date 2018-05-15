import {Component, OnInit} from '@angular/core';
import {ApplicationsStoreService, LoadApplicationsAction, SelectApplicationAction} from '../../stores/applications-store.service';
import {Store} from '@ngrx/store';
import {ApplicationBean, DeploymentBean, DomainBean} from '../../models/commons/applications-bean';
import {DataApplicationService} from '../../services/data-application-version.service';
import {ContentListResponse} from '../../models/commons/entity-bean';

import {Edge, Node} from '../../models/graph/graph-bean';
import {DataDeploymentService} from '../../services/data-deployment.service';

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
  public domains: DomainBean[] = [];

  /**
   * to graph domains
   */
  public nodes: Node[];
  public edges: Edge[];

  constructor(
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService,
    private deploymentService: DataDeploymentService
  ) {
    /**
     * subscribe
     */
    this.domainsStream = this.applicationsStoreService.domains();

    this.domainsStream.subscribe(
      (domains: DomainBean[]) => {
        this.nodes = [];
        this.edges = [];
        this.domains = domains;
        const indexApp = new Map<string, number>();
        domains.forEach((domain) => {
          this.nodes.push({
            id: domain.name,
            label: domain.name
          });
          const applications = domain.applications;
          applications.forEach((app) => {
            const keyApp = domain.name + '#' + app.name;
            if (indexApp.has(keyApp)) {
              const instance = indexApp.get(keyApp) + 1;
              indexApp.set(keyApp, instance);
            } else {
              indexApp.set(keyApp, 0);
            }
            this.nodes.push({
              id: keyApp + '#' + indexApp.get(keyApp),
              label: app.name + '#' + indexApp.get(keyApp)
            });
            this.edges.push({
              from: domain.name,
              to: keyApp + '#' + indexApp.get(keyApp),
              label: domain.name + ' to ' + app.name + '#' + indexApp.get(keyApp)
            });
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
    if (!this.domains || this.domains.length === 0) {
      this.loadApplications(null);
    }
  }

  /**
   * dispatch load applications
   * @param event
   */
  protected loadApplications(event: any) {
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


  /**
   * dispatch load applications
   * @param event
   */
  protected selectApplication(application: ApplicationBean) {
    // load all applications from a content return
    this.applicationsService.GetAllFromContent(`/${application.domain}/${application.name}/versions`, <Map<string, string>> {size: 1})
      .subscribe(
        (data: ContentListResponse<ApplicationBean>) => {
          this.deploymentService.GetAllFromContent(
            '/?q=%7B%22properties._app_domain%22%3A%20%22' + application.domain +
            '%22%2C%20%22properties._app_name%22%3A%20%22' + application.name + '%22%7D',
            <Map<string, string>> {size: 20}).subscribe(
            (deployments: ContentListResponse<DeploymentBean>) => {
              this.applicationsStoreService.dispatch(
                new SelectApplicationAction(
                  data.content[0],
                  deployments.content
                )
              );
            }
          );
        }
      );
  }
}
