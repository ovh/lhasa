import { Injectable } from '@angular/core';
import { createFeatureSelector, createSelector, Selector, Store } from '@ngrx/store';
import { ActionWithPayloadAndPromise } from './action-with-payload';
import { ApplicationBean, ApplicationPagesBean, DeploymentBean, DomainBean, DomainPagesBean } from '../models/commons/applications-bean';
import { Subject, Observable } from 'rxjs';
import { GraphBean } from '../models/graph/graph-bean';

import { remove } from 'lodash';

// Error
export class ErrorBean {
  code: string;
  stack?: any;
}

/**
 * states
 */
export interface ErrorsState {
  /**
   * store each graph with a key
   */
  errors: ErrorBean[];
}

/**
 * actions
 */
export class NewErrorAction implements ActionWithPayloadAndPromise<ErrorBean> {
  readonly type = NewErrorAction.getType();

  public static getType(): string {
    return 'NewErrorAction';
  }

  constructor(public payload: ErrorBean, public subject?: Subject<any>) {
  }
}

export class DropErrorAction implements ActionWithPayloadAndPromise<ErrorBean> {
  readonly type = DropErrorAction.getType();

  public static getType(): string {
    return 'DropErrorAction';
  }

  constructor(public payload: ErrorBean, public subject?: Subject<any>) {
  }
}

export type AllStoreActions = NewErrorAction | DropErrorAction;

/**
 * main store for this Graph
 */
@Injectable()
export class ErrorsStoreService {

  readonly getErrors: Selector<object, ErrorBean[]>;

  /**
   *
   * @param _store constructor
   */
  constructor(
    private _store: Store<ErrorsState>
  ) {
    this.getErrors = ErrorsStoreService.create((state: ErrorsState) => state.errors);
  }

  /**
   * create a selector
   * @param handler internal static
   */
  private static create(handler: (S1: ErrorsState) => any)  {
    return createSelector(createFeatureSelector<ErrorsState>('errors'), handler);
  }

  /**
   * metareducer (Cf. https://www.concretepage.com/angular-2/ngrx/ngrx-store-4-angular-5-tutorial)
   * @param state
   * @param action
   */
  public static reducer(state: ErrorsState = {
    errors: <ErrorBean[]>[],
  }, action: AllStoreActions): ErrorsState {

    switch (action.type) {
      /**
       * update store
       */
      case NewErrorAction.getType(): {
        const errors = <ErrorBean[]> Object.assign([], state.errors);
        errors.push(action.payload);

        // Complete load action
        action.subject.complete();
        return {
          errors: errors,
        };
      }

      /**
       * update store
       */
      case DropErrorAction.getType(): {
        const errors = <ErrorBean[]> Object.assign([], state.errors);

        remove(errors, (error) => {
          return error.code === action.payload.code;
        });
        errors.push(action.payload);

        // Complete load action
        action.subject.complete();
        return {
          errors: errors,
        };
      }

      default:
        return state;
    }
  }

  /**
   * select this store service
   */
  public errors(): Observable<ErrorBean[]> {
    return this._store.select(this.getErrors);
  }

  /**
   * dispatch
   * @param action dispatch action
   */
  public dispatch(action: AllStoreActions) {
    this._store.dispatch(action);
  }
}
