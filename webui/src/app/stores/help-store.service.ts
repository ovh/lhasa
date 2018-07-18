import { BehaviorSubject ,  Subject } from 'rxjs';
import { Injectable } from '@angular/core';
import { createFeatureSelector, createSelector, Selector, Store } from '@ngrx/store';

import { ActionWithPayloadAndPromise } from './action-with-payload';
import { ContentBean } from '../models/commons/content-bean';
import { DataContentService } from '../services/data-content.service';
import { LoadersStoreService } from './loader-store.service';
import { ErrorsStoreService, NewErrorAction, ErrorBean } from './errors-store.service';
import { Observable } from 'rxjs/internal/Observable';

// Help
export class HelpBean {
  token: string;
  content: ContentBean;
}

/**
 * states
 */
export interface HelpState {
  /**
   * help
   */
  help: HelpBean;
}

/**
 * actions
 */
export class GetHelpAction implements ActionWithPayloadAndPromise<HelpBean> {
  readonly type = GetHelpAction.getType();

  public static getType(): string {
    return 'GetHelpAction';
  }

  constructor(public payload: HelpBean, public subject?: Subject<any>) {
  }
}

export type AllStoreActions = GetHelpAction;

/**
 * main store for this application
 */
@Injectable()
export class HelpsStoreService {

  readonly getHelp: Selector<object, HelpBean>;

  /**
   *
   * @param _store constructor
   */
  constructor(
    private _store: Store<HelpState>,
    private _content: DataContentService,
    private loadersStoreService: LoadersStoreService,
    private errorsStoreService: ErrorsStoreService,
  ) {
    this.getHelp = HelpsStoreService.create((state: HelpState) => state.help);
  }

  /**
   * create a selector
   * @param handler internal static
   */
  private static create(handler: (S1: HelpState) => any) {
    return createSelector(createFeatureSelector<HelpState>('helps'), handler);
  }

  /**
   * metareducer (Cf. https://www.concretepage.com/angular-2/ngrx/ngrx-store-4-angular-5-tutorial)
   * @param state
   * @param action
   */
  public static reducer(state: HelpState = {
    help: new HelpBean(),
  }, action: AllStoreActions): HelpState {

    switch (action.type) {
      /**
       * add a new loader
       */
      case GetHelpAction.getType(): {
        const n = <HelpBean>Object.assign({}, action.payload);
        action.subject.complete();
        return {
          help: n,
        };
      }
      default:
        return state;
    }
  }

  /**
   * select this store service
   */
  public help(): Observable<HelpBean> {
    return this._store.select(this.getHelp);
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
    const subject: Subject<any> = new BehaviorSubject<any>('select help token ' + key);
    this.loadersStoreService.notify(subject);
    this._content.GetSingle(key).subscribe(
      (data: ContentBean) => {
        this.dispatch(
          new GetHelpAction(<HelpBean>{
            token: key,
            content: data
          }, subject));
      },
      (error) => {
        // When errors only print a default help value
        this.dispatch(
          new GetHelpAction(<HelpBean>{
            token: key,
            content: 'TBD ...'
          }, subject));
      }
    );
  }

}
