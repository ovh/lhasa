import { Observable } from 'rxjs/Observable';
import { Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Injectable } from '@angular/core';
import { BadgeBean, BadgePagesBean } from '../models/commons/badges-bean';
import { LoadBadgesAction, BadgesStoreService } from '../stores/badges-store.service';
import { ContentListResponse, PageMetaData } from '../models/commons/entity-bean';
import { DataDeploymentService } from '../services/data-deployment.service';
import { DataBadgeService } from '../services/data-badge.service';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';
import { LoadersStoreService } from '../stores/loader-store.service';

@Injectable()
export class BadgesResolver implements Resolve<BadgeBean[]> {
    constructor(
        private badgesStoreService: BadgesStoreService,
        private badgesService: DataBadgeService,
        private loadersStoreService: LoadersStoreService
    ) { }

    resolve(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {
        return this.selectBadges({
            number: route.queryParams.page || 0,
            size: 20
        }, new BehaviorSubject<any>('select all badges'));
    }

    /**
     * dispatch load badges
     * @param event
     */
    public selectBadges(metadata: PageMetaData, subject: Subject<any>): Subject<any> {
        this.loadersStoreService.notify(subject);
        // load all badges
        const meta: {
            [key: string]: any | any[];
        } = {
            size: metadata.size,
            page: metadata.number
        };
        this.badgesService.GetAllFromContent('', meta).subscribe(
            (data: ContentListResponse<BadgeBean>) => {
                this.badgesStoreService.dispatch(
                    new LoadBadgesAction({
                        badges: data.content,
                        metadata: data.pageMetadata,
                    }, subject)
                );
            }
        );
        return subject;
    }
}
