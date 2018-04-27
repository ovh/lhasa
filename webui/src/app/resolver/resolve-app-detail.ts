import { Observable } from 'rxjs/Observable';
import { Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Injectable } from '@angular/core';
import { ApplicationBean, DeploymentBean } from '../models/commons/applications-bean';
import { ApplicationsStoreService, SelectApplicationAction } from '../stores/applications-store.service';
import { DataApplicationService } from '../services/data-application-version.service';
import { ContentListResponse } from '../models/commons/entity-bean';
import { DataDeploymentService } from '../services/data-deployment.service';

@Injectable()
export class ApplicationResolver implements Resolve<ApplicationBean> {
    constructor(
        private applicationsStoreService: ApplicationsStoreService,
        private applicationsService: DataApplicationService,
        private deploymentService: DataDeploymentService
    ) { }

    resolve(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {
        this.selectApplication(route.params.domain, route.params.name);
    }

    /**
     * dispatch load applications
     * @param event
     */
    protected selectApplication(domaine: string, name: string) {
        // load all applications from a content return
        this.applicationsService.GetAllFromContent('/' + domaine + '/' + name, <Map<string, string>>{ size: 1 })
            .subscribe(
                (data: ContentListResponse<ApplicationBean>) => {
                    this.deploymentService.GetAllFromContent(
                        '/?q=%7B%22properties._app_domain%22%3A%20%22' + domaine +
                        '%22%2C%20%22properties._app_name%22%3A%20%22' + name + '%22%7D',
                        <Map<string, string>>{ size: 20 }).subscribe(
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
