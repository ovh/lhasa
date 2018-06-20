import { Observable } from 'rxjs/Observable';
import { Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';
import { LoadersStoreService } from '../stores/loader-store.service';
import { GraphBean } from '../models/graph/graph-bean';
import { DataGraphService } from '../services/data-graph.service';
import { GraphsStoreService, LoadGraphDeploymentAction } from '../stores/graphs-store.service';
import { ErrorsStoreService, ErrorBean, NewErrorAction } from '../stores/errors-store.service';

@Injectable()
export class GraphsResolver implements Resolve<GraphBean> {
    constructor(
        private dataGraphService: DataGraphService,
        private graphsStoreService: GraphsStoreService,
        private errorsStoreService: ErrorsStoreService,
        private loadersStoreService: LoadersStoreService,
    ) { }

    resolve(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {
        // Only one graph at this moment, deployment graph
        // Soon may be we have to see route params to select the graph
        return this.selectGraphDeployment({}, new BehaviorSubject<any>('select one graph with query'));
    }

    /**
     * dispatch load Graphs
     * @param event
     */
    public selectGraphDeployment(params: any, subject: Subject<any>): Subject<any> {
        this.loadersStoreService.notify(subject);
        // load all Graphs
        this.dataGraphService.Get(params).subscribe(
            (data: GraphBean) => {
                this.graphsStoreService.dispatch(
                    new LoadGraphDeploymentAction(data, subject)
                );
            },
            (error) => {
                this.errorsStoreService.dispatch(new NewErrorAction(
                    <ErrorBean>{
                        code: 'ERROR-GRAPHS',
                        stack: JSON.stringify(error, null, 2),
                    }, subject
                ));
            }
        );
        return subject;
    }
}
