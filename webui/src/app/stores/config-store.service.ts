import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Injectable } from '@angular/core';
import { createFeatureSelector, createSelector, Selector, Store } from '@ngrx/store';

import { ActionWithPayloadAndPromise } from './action-with-payload';
import { Subject } from 'rxjs/Subject';
import { ContentBean } from '../models/commons/content-bean';
import { DataContentService } from '../services/data-content.service';
import { LoadersStoreService } from './loader-store.service';
import { ErrorsStoreService, NewErrorAction, ErrorBean } from './errors-store.service';

// Config
export class ConfigBean {
  config: any;
}

/**
 * states
 */
export interface ConfigState {
  /**
   * help
   */
  config: ConfigBean;
}

/**
 * actions
 */
export class SetConfigAction implements ActionWithPayloadAndPromise<ConfigBean> {
  readonly type = SetConfigAction.getType();

  public static getType(): string {
    return 'SetConfigAction';
  }

  constructor(public payload: ConfigBean, public subject?: Subject<any>) {
  }
}

export type AllStoreActions = SetConfigAction;

/**
 * main store for this application
 */
@Injectable()
export class ConfigStoreService {

  readonly getConfig: Selector<object, ConfigBean>;

  /**
   *
   * @param _store constructor
   */
  constructor(
    private _store: Store<ConfigState>,
    private _content: DataContentService,
    private loadersStoreService: LoadersStoreService,
    private errorsStoreService: ErrorsStoreService,
  ) {
    this.getConfig = ConfigStoreService.create((state: ConfigState) => state.config);
  }

  /**
   * create a selector
   * @param handler internal static
   */
  private static create(handler: (S1: ConfigState) => any) {
    return createSelector(createFeatureSelector<ConfigState>('config'), handler);
  }

  /**
   * metareducer (Cf. https://www.concretepage.com/angular-2/ngrx/ngrx-store-4-angular-5-tutorial)
   * @param state
   * @param action
   */
  public static reducer(state: ConfigState = {
    config: new ConfigBean(),
  }, action: AllStoreActions): ConfigState {

    switch (action.type) {
      /**
       * add a new loader
       */
      case SetConfigAction.getType(): {
        const n = <ConfigBean>Object.assign({}, action.payload);
        action.subject.complete();
        return {
          config: n,
        };
      }
      default:
        return state;
    }
  }

  /**
   * select this store service
   */
  public help(): Store<ConfigBean> {
    return this._store.select(this.getConfig);
  }

  /**
   * dispatch
   * @param action dispatch action
   */
  public dispatch(action: AllStoreActions) {
    this._store.dispatch(action);
  }

  /**
   * dispatch
   * @param action dispatch action
   */
  public request(key: string) {
    const subject: Subject<any> = new BehaviorSubject<any>('select config token ' + key);
    this.loadersStoreService.notify(subject);
    this._content.GetSingle(key).subscribe(
      (data: ContentBean) => {
        let configuration;
        try {
          configuration = JSON.parse(data.toString());
        } catch(e) {
          console.warn(e);
          configuration = {};
        }
        this.dispatch(
          new SetConfigAction(<ConfigBean>{
            config: configuration
          }, subject));
      },
      (error) => {
        // When errors load a default config
        this.dispatch(
          new SetConfigAction(<ConfigBean>{
            config: {}
          }, subject));
      }
    );
  }

}
