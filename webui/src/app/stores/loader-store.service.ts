import { Injectable } from '@angular/core';
import { createFeatureSelector, createSelector, Store, Selector } from '@ngrx/store';

import { ActionWithPayload, ActionWithPayloadAndPromise } from './action-with-payload';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';
import { remove } from 'lodash';

// Loader
export class LoaderBean {
  id: number;
  store?: Store<LoaderState>;
  subject: Subject<any>;
}

/**
 * states
 */
export interface LoaderState {
  /**
   * loaders
   */
  loaders: LoaderBean[];
}

/**
 * actions
 */
export class AddLoaderAction implements ActionWithPayloadAndPromise<LoaderBean> {
  readonly type = AddLoaderAction.getType();

  public static getType(): string {
    return 'AddLoaderAction';
  }

  constructor(public payload: LoaderBean, public subject?: Subject<any>) {
  }
}

export class LoaderResolvedAction implements ActionWithPayloadAndPromise<LoaderBean> {
  readonly type = LoaderResolvedAction.getType();

  public static getType(): string {
    return 'LoaderResolvedAction';
  }

  constructor(public payload: LoaderBean, public subject?: Subject<any>) {
  }
}

export type AllStoreActions = AddLoaderAction | LoaderResolvedAction;

/**
 * main store for this application
 */
@Injectable()
export class LoadersStoreService {

  readonly getLoaders: Selector<object, LoaderBean[]>;
  private counters = 0;

  /**
   *
   * @param _store constructor
   */
  constructor(
    private _store: Store<LoaderState>
  ) {
    this.getLoaders = LoadersStoreService.create((state: LoaderState) => state.loaders);
  }

  /**
   * create a selector
   * @param handler internal static
   */
  private static create(handler: (S1: LoaderState) => any) {
    return createSelector(createFeatureSelector<LoaderState>('loaders'), handler);
  }

  /**
   * metareducer (Cf. https://www.concretepage.com/angular-2/ngrx/ngrx-store-4-angular-5-tutorial)
   * @param state
   * @param action
   */
  public static reducer(state: LoaderState = {
    loaders: [],
  }, action: AllStoreActions): LoaderState {

    switch (action.type) {
      /**
       * add a new loader
       */
      case AddLoaderAction.getType(): {
        const loaders = <LoaderBean[]>Object.assign([], state.loaders);
        loaders.push(action.payload);

        const target = action.payload;
        target.subject.asObservable().subscribe(
          (data) => {
          },
          (error) => {
            console.error('errors', target.subject);
          },
          () => {
            target.store.dispatch(new LoaderResolvedAction({
              id: target.id,
              subject: target.subject
            }, new BehaviorSubject<any>('loaded')));
          }
        );

        action.subject.complete();
        return {
          loaders: loaders,
        };
      }

      /**
       * add a new loader
       */
      case LoaderResolvedAction.getType(): {
        const loaders = <LoaderBean[]>Object.assign([], state.loaders);
        console.error('complete', action.payload, loaders);
        remove(loaders, (loader) => {
          return loader.subject.isStopped;
        });
        console.error('complete', action.payload, loaders);
        return {
          loaders: loaders,
        };
      }

      default:
        return state;
    }
  }

  /**
   * select this store service
   */
  public loaders(): Store<LoaderBean[]> {
    return this._store.select(this.getLoaders);
  }

  /**
   * dispatch
   * @param action dispatch action
   */
  public dispatch(action: AllStoreActions) {
    this._store.dispatch(action);
  }

  /**
   * notify
   * @param action notify action
   */
  public notify(subject: Subject<any>): Subject<any> {
    this._store.dispatch(new AddLoaderAction({
      id: this.counters++,
      store: this._store,
      subject: subject
    }, new BehaviorSubject<any>('loading')));
    return subject;
  }
}
