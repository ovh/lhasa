import { Observable } from 'rxjs/Observable';
import { Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Injectable } from '@angular/core';
import { BadgeBean, BadgePagesBean } from '../models/commons/badges-bean';
import { LoadBadgesAction, BadgesStoreService } from '../stores/badges-store.service';
import { ContentListResponse, PageMetaData } from '../models/commons/entity-bean';
import { DataDeploymentService } from '../services/data-deployment.service';
import { DataBadgeService } from '../services/data-badge.service';
import { DataBadgeStatsService } from '../services/data-badgestats.service';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';
import { LoadersStoreService } from '../stores/loader-store.service';
import { NewErrorAction, ErrorBean, ErrorsStoreService } from '../stores/errors-store.service';

@Injectable()
export class BadgesResolver implements Resolve<BadgeBean[]> {
    constructor(
        private badgesStoreService: BadgesStoreService,
        private badgesService: DataBadgeService,
        private badgeStatsService: DataBadgeStatsService,
        private loadersStoreService: LoadersStoreService,
        private errorsStoreService: ErrorsStoreService,
    ) { }

    resolve(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {
        return this.selectBadges({
            number: route.queryParams.page || 0,
            size: 40
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
                //var stats = {}
                var i = 0
                data.content.forEach(badge => {
                    this.badgeStatsService.GetBadgeStats(badge.slug).subscribe(
                        (stats: Map<string, number>) => {
                            console.log("stats", stats);
                            badge._stats = stats
                            if (i == 3) {
                                this.badgesStoreService.dispatch(
                                    new LoadBadgesAction({
                                        badges: data.content,
                                        metadata: data.pageMetadata,
                                    }, subject)
                                );
                            }
                            i += 1;
                        })
                });

            },
            (error) => {
                this.errorsStoreService.dispatch(new NewErrorAction(
                    <ErrorBean>{
                        code: 'ERROR-BADGES',
                        stack: JSON.stringify(error, null, 2),
                    }, subject
                ));
            }
        );
        return subject;
    }
}
