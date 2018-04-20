import {Injectable} from '@angular/core';
import {ActionReducer, Action, State} from '@ngrx/store';
import {Store} from '@ngrx/store';
import {createFeatureSelector, createSelector, MemoizedSelector} from '@ngrx/store';

import * as _ from 'lodash';

import {ActionWithPayload} from './action-with-payload';
import {ApplicationBean, DeploymentBean, EnvironmentBean} from '../models/commons/applications-bean';
import {ApplicationState} from './applications-store.service';

/**
 * states
 */
export interface EnvironmentState {
  environments: Map<string, EnvironmentBean>;
}

/**
 * actions
 */
export class LoadEnvironmentsAction implements ActionWithPayload<Map<string, EnvironmentBean>> {
  readonly type = 'LoadEnvironmentsAction';

  constructor(public payload: Map<string, EnvironmentBean>) {
  }
}

export type AllStoreActions = LoadEnvironmentsAction;

/**
 * main store for this application
 */
@Injectable()
export class EnvironmentsStoreService {

  private getAll: MemoizedSelector<object, Map<string, EnvironmentBean>>;

  /**
   * metareducer (Cf. https://www.concretepage.com/angular-2/ngrx/ngrx-store-4-angular-5-tutorial)
   * @param state
   * @param action
   */
  public static reducer(state: EnvironmentState = {
    environments: new Map<string, EnvironmentBean>(),
  }, action: AllStoreActions): EnvironmentState {

    switch (action.type) {
      /**
       * message incomming
       */
      case 'LoadEnvironmentsAction': {
        const environments = Object.assign(new Map<string, EnvironmentBean>(), action.payload);
        return {
          environments: environments,
        };
      }

      default:
        return state;
    }
  }

  /**
   *
   * @param _store constructor
   */
  constructor(
    private _store: Store<EnvironmentBean>
  ) {
    this.getAll = createSelector(createFeatureSelector<EnvironmentState>('environments'), (state: EnvironmentState) => state.environments);
  }

  /**
   * select this store service
   */
  public environments(): Store<Map<string, EnvironmentBean>> {
    return this._store.select(this.getAll);
  }

  /**
   * dispatch
   * @param action dispatch action
   */
  public dispatch(action: AllStoreActions) {
    this._store.dispatch(action);
  }
}
