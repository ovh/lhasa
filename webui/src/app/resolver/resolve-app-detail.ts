import { Observable, BehaviorSubject, Subject } from 'rxjs';
import { Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Injectable } from '@angular/core';
import { ApplicationBean, DeploymentBean, BadgeRatingBean } from '../models/commons/applications-bean';
import { ApplicationsStoreService, SelectApplicationAction } from '../stores/applications-store.service';
import { DataApplicationService } from '../services/data-application-version.service';
import { ContentListResponse } from '../models/commons/entity-bean';
import { DataDeploymentService } from '../services/data-deployment.service';
import { DataBadgeRatingsService } from '../services/data-badgeratings.service';
import { LoadersStoreService } from '../stores/loader-store.service';
import { ErrorsStoreService, NewErrorAction, ErrorBean } from '../stores/errors-store.service';
import { DataGraphService } from '../services/data-graph.service';
import { GraphBean, GraphVis } from '../models/graph/graph-bean';

@Injectable()
export class ApplicationResolver implements Resolve<ApplicationBean> {
  constructor(
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService,
    private deploymentService: DataDeploymentService,
    private badgeRatingService: DataBadgeRatingsService,
    private loadersStoreService: LoadersStoreService,
    private errorsStoreService: ErrorsStoreService,
    private graphService: DataGraphService
  ) {
  }

  resolve(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<any> | Promise<any> | any {
    this.selectApplication(route.params.domain,
      route.params.name, new BehaviorSubject<any>('select latest application'));
  }

  /**
   * dispatch  SelectApplicationAction
   * @param event
   */
  protected selectApplication(domain: string, name: string, subject: Subject<any>) {
    this.loadersStoreService.notify(subject);
    this.applicationsService.GetSingle(`${domain}/${name}/latest`)
      .subscribe(
        (app: ApplicationBean) => {
          // application service should answer apps, but here we have deployments
          this.applicationsService.GetSingleAny(`${domain}/${name}/deployments`)
            .subscribe(
              (deployments: DeploymentBean[]) => {
                const deploymentsGraphed = deployments.map(
                  (value: DeploymentBean, index: number, array: DeploymentBean[]): DeploymentBean => {
                    this.graphService.GetDeployment(value.id).first().subscribe((graph: GraphBean) => {
                      // Compute data
                      const g = new GraphVis();
                      g.nodes = [];
                      graph.nodes.forEach(node => {
                        let group = node.properties.environment.slug;
                        if (node.id === value.id) {
                          group = 'highlight';
                        }
                        g.nodes.push({
                          id: node.id,
                          label: node.name,
                          group: group,
                          environment: node.properties.environment.slug,
                          domain: node.properties.application.domain,
                          application: node.properties.application.name,
                        });
                      });
                      // Compute data
                      g.edges = [];
                      graph.edges.forEach(edge => {
                        g.edges.push({
                          id: edge.id,
                          from: edge.from,
                          to: edge.to,
                          label: edge.type
                        });
                      });
                      value._graph = g;
                    });
                    return value;
                  });

                this.badgeRatingService.GetBadgeRatings(`${domain}/${name}/versions/${app.version}/badges`).subscribe(
                  (badgeRatings: BadgeRatingBean[]) => {
                    app.deployments = deploymentsGraphed;
                    app.badgeRatings = badgeRatings;
                    this.applicationsStoreService.dispatch(
                      new SelectApplicationAction(
                        app,
                        subject
                      )
                    );
                  }
                );
              }
            );
        },
        (error) => {
          this.errorsStoreService.dispatch(new NewErrorAction(
            <ErrorBean>{
              code: 'ERROR-APP-DETAILS',
              stack: JSON.stringify(error, null, 2),
            }, subject
          ));
        }
      );
  }
}
