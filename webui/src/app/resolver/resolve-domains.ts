import { Observable } from 'rxjs/Observable';
import { Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Injectable } from '@angular/core';
import { DomainBean, DomainPagesBean } from '../models/commons/applications-bean';
import { LoadDomainsAction, ApplicationsStoreService } from '../stores/applications-store.service';
import { DataApplicationService } from '../services/data-application-version.service';
import { ContentListResponse, PageMetaData } from '../models/commons/entity-bean';
import { DataDeploymentService } from '../services/data-deployment.service';
import { DataDomainService } from '../services/data-domain.service';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';
import { LoadersStoreService } from '../stores/loader-store.service';

@Injectable()
export class DomainsResolver implements Resolve<DomainBean[]> {
    constructor(
        private applicationsStoreService: ApplicationsStoreService,
        private domainsService: DataDomainService,
        private loadersStoreService: LoadersStoreService
    ) { }

    resolve(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {
        return this.selectDomains({
            number: route.queryParams.page || 0,
            size: 5
        }, new BehaviorSubject<any>('select all domains'));
    }

    /**
     * dispatch load domains
     * @param event
     */
    public selectDomains(metadata: PageMetaData, subject: Subject<any>): Subject<any> {
        this.loadersStoreService.notify(subject);
        // load all domains
        const meta: {
            [key: string]: any | any[];
        } = {
            size: metadata.size,
            page: metadata.number
        };
        this.domainsService.GetAllFromContent('', meta).subscribe(
            (data: ContentListResponse<DomainBean>) => {
                this.applicationsStoreService.dispatch(
                    new LoadDomainsAction({
                        domains: data.content,
                        metadata: data.pageMetadata,
                    }, subject)
                );
            }
        );
        return subject;
    }
}
