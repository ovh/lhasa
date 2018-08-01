import { Observable ,  BehaviorSubject ,  Subject } from 'rxjs';
import { Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Injectable } from '@angular/core';
import { ApplicationBean, ApplicationPagesBean } from '../models/commons/applications-bean';
import { LoadDomainsAction, ApplicationsStoreService, LoadApplicationsAction } from '../stores/applications-store.service';
import { DataApplicationService } from '../services/data-application-version.service';
import { ContentListResponse, PageMetaData } from '../models/commons/entity-bean';
import { DataDeploymentService } from '../services/data-deployment.service';
import { DataDomainService } from '../services/data-domain.service';
import { LoadersStoreService } from '../stores/loader-store.service';
import { ErrorsStoreService, ErrorBean, NewErrorAction } from '../stores/errors-store.service';
import { Params } from '@angular/router/src/shared';
import { HttpParams } from '@angular/common/http/src/params';

@Injectable()
export class ApplicationsResolver implements Resolve<ApplicationPagesBean> {
    constructor(
        private applicationsStoreService: ApplicationsStoreService,
        private applicationsService: DataApplicationService,
        private loadersStoreService: LoadersStoreService,
        private errorsStoreService: ErrorsStoreService,
    ) {

    }

    resolve(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {
        return this.selectApplications( route.queryParams, new BehaviorSubject<any>('select all applications'));
    }

    /**
     * dispatch load domains
     * @param event
     */
    public selectApplications(params: Params, subject: Subject<any>): Subject<any> {
        this.loadersStoreService.notify(subject);
        const paramsClone = Object.assign({}, params);
        paramsClone.sort = 'domain';
        if (paramsClone.size === undefined) {
            paramsClone.size = 100;
        }
        paramsClone.sort = 'domain';
        this.applicationsService.GetAllFromContent('', paramsClone).subscribe(
            (data: ContentListResponse<ApplicationBean>) => {
                this.applicationsStoreService.dispatch(
                    new LoadApplicationsAction({
                        applications: data.content,
                        metadata: data.pageMetadata,
                    }, subject)
                );
            },
            (error) => {
                this.errorsStoreService.dispatch(new NewErrorAction(
                    <ErrorBean>{
                        code: 'ERROR-APPLICATION',
                        stack: JSON.stringify(error, null, 2),
                    }, subject
                ));
            }
        );
        return subject;
    }
}
