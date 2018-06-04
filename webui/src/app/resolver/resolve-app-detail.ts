import { Observable } from 'rxjs/Observable';
import { Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Injectable } from '@angular/core';
import { ApplicationBean, DeploymentBean } from '../models/commons/applications-bean';
import { ApplicationsStoreService, SelectApplicationAction } from '../stores/applications-store.service';
import { DataApplicationService } from '../services/data-application-version.service';
import { ContentListResponse } from '../models/commons/entity-bean';
import { DataDeploymentService } from '../services/data-deployment.service';
import { LoadersStoreService } from '../stores/loader-store.service';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class ApplicationResolver implements Resolve<ApplicationBean> {
    constructor(
        private applicationsStoreService: ApplicationsStoreService,
        private applicationsService: DataApplicationService,
        private deploymentService: DataDeploymentService,
        private loadersStoreService: LoadersStoreService
    ) { }

    resolve(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {
        this.selectApplication(route.params.domain, route.params.name, route.params.version, new BehaviorSubject<any>('select an application'));
    }

    /**
     * dispatch  SelectApplicationAction
     * @param event
     */
    protected selectApplication(domain: string, name: string, version: string, subject: Subject<any>) {
        this.loadersStoreService.notify(subject);
        this.applicationsService.GetSingle(`${domain}/${name}/versions/${version}`)
            .subscribe(
                (data: ApplicationBean) => {
                    this.deploymentService.GetAllFromContent(
                        '/?q=%7B%22properties._app_domain%22%3A%20%22' + domain +
                        '%22%2C%20%22properties._app_name%22%3A%20%22' + name + '%22%7D',
                        { 'size': '20' }).subscribe(
                            (deployments: ContentListResponse<DeploymentBean>) => {
                                this.applicationsStoreService.dispatch(
                                    new SelectApplicationAction(
                                        data,
                                        deployments.content, 
                                        subject)
                                );
                            }
                        );
                }
            );
    }
}
