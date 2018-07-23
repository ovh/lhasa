import { Observable ,  BehaviorSubject ,  Subject } from 'rxjs';
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

@Injectable()
export class ApplicationResolver implements Resolve<ApplicationBean> {
    constructor(
        private applicationsStoreService: ApplicationsStoreService,
        private applicationsService: DataApplicationService,
        private deploymentService: DataDeploymentService,
        private badgeRatingService: DataBadgeRatingsService,
        private loadersStoreService: LoadersStoreService,
        private errorsStoreService: ErrorsStoreService,
    ) { }

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
                        this.badgeRatingService.GetBadgeRatings(`${domain}/${name}/versions/${app.version}/badges`).subscribe(
                            (badgeRatings: BadgeRatingBean[]) => {
                                app.deployments = deployments;
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
