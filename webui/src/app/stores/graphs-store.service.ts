import { Injectable } from '@angular/core';
import { createFeatureSelector, createSelector, Selector, Store } from '@ngrx/store';

import { ActionWithPayloadAndPromise } from './action-with-payload';
import { ApplicationBean, ApplicationPagesBean, DeploymentBean, DomainBean, DomainPagesBean } from '../models/commons/applications-bean';
import { Subject } from 'rxjs/Subject';
import { GraphBean } from '../models/graph/graph-bean';

/**
 * states
 */
export interface GraphsState {
  /**
   * store each graph with a key
   */
  deployments: GraphBean;
}

/**
 * actions
 */
export class LoadGraphDeploymentAction implements ActionWithPayloadAndPromise<GraphBean> {
  readonly type = LoadGraphDeploymentAction.getType();

  public static getType(): string {
    return 'LoadGraphDeploymentAction';
  }

  constructor(public payload: GraphBean, public subject?: Subject<any>) {
  }
}

export type AllStoreActions = LoadGraphDeploymentAction;

/**
 * main store for this Graph
 */
@Injectable()
export class GraphsStoreService {

  readonly getDeploymentGraph: Selector<object, GraphBean>;

  /**
   *
   * @param _store constructor
   */
  constructor(
    private _store: Store<GraphsState>
  ) {
    this.getDeploymentGraph = GraphsStoreService.create((state: GraphsState) => state.deployments);
  }

  /**
   * create a selector
   * @param handler internal static
   */
  private static create(handler: (S1: GraphsState) => any)  {
    return createSelector(createFeatureSelector<GraphsState>('graphs'), handler);
  }

  /**
   * metareducer (Cf. https://www.concretepage.com/angular-2/ngrx/ngrx-store-4-angular-5-tutorial)
   * @param state
   * @param action
   */
  public static reducer(state: GraphsState = {
    deployments: new GraphBean(),
  }, action: AllStoreActions): GraphsState {

    switch (action.type) {
      /**
       * update all applications in store
       */
      case LoadGraphDeploymentAction.getType(): {
        const graph = <GraphBean> Object.assign({}, action.payload);

        // Complete load action
        action.subject.complete();
        return {
          deployments: graph,
        };
      }


      default:
        return state;
    }
  }

  /**
   * select this store service
   */
  public deployments(): Store<GraphBean> {
    return this._store.select(this.getDeploymentGraph);
  }

  /**
   * dispatch
   * @param action dispatch action
   */
  public dispatch(action: AllStoreActions) {
    this._store.dispatch(action);
  }
}
