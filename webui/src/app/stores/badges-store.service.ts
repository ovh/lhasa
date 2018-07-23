import { Injectable } from '@angular/core';
import { createFeatureSelector, createSelector, Selector, Store } from '@ngrx/store';

import { ActionWithPayloadAndPromise } from './action-with-payload';
import {  BadgeBean, BadgePagesBean } from '../models/commons/badges-bean';
import { Subject } from 'rxjs';
import { Observable } from 'rxjs';

/**
 * states
 */
export interface BadgeState {
  /**
   * badges of each application loaded in store
   */
  badgePages: BadgePagesBean;
  /**
   * badges of each application loaded in store
   */
  badges: Array<BadgeBean>;
}

/**
 * actions
 */


export class LoadBadgesAction implements ActionWithPayloadAndPromise<BadgePagesBean> {
  readonly type = LoadBadgesAction.getType();

  public static getType(): string {
    return 'LoadBadgesAction';
  }

  constructor(public payload: BadgePagesBean, public subject: Subject<any>) {
  }
}

export type AllStoreActions =  LoadBadgesAction ;

/**
 * main store for this application
 */
@Injectable()
export class BadgesStoreService {

  readonly getBadgePages: Selector<object, BadgePagesBean>;
  readonly getBadges: Selector<object, Array<BadgeBean>>;

  /**
   *
   * @param _store constructor
   */
  constructor(
    private _store: Store<BadgeState>
  ) {
    this.getBadges = createSelector(createFeatureSelector<BadgeState>('badges'), (state: BadgeState) => state.badges);
    this.getBadgePages = BadgesStoreService.create((state: BadgeState) => state.badgePages);
  }

  /**
   * create a selector
   * @param handler internal static
   */
  private static create(handler: (S1: BadgeState) => any)  {
    return createSelector(createFeatureSelector<BadgeState>('badges'), handler);
  }

  /**
   * metareducer (Cf. https://www.concretepage.com/angular-2/ngrx/ngrx-store-4-angular-5-tutorial)
   * @param state
   * @param action
   */
  public static reducer(state: BadgeState = {
    badgePages: new BadgePagesBean(),
    badges: new Array<BadgeBean>(),
  }, action: AllStoreActions): BadgeState {

    switch (action.type) {
      /**
       * update all badges in store
       */
      case LoadBadgesAction.getType(): {
        const badgePages = Object.assign(new BadgePagesBean(), action.payload);

        /**
         * notify badges change
         */
        action.subject.complete();
        return {
          badgePages: badgePages,
          badges: state.badges,
        };
      }

      default:
        return state;
    }
  }

  /**
   * select this store service
   */
  public badgePages(): Observable<BadgePagesBean> {
    return this._store.select(this.getBadgePages);
  }

  /**
   * select this store service
   */
  public badges(): Observable<Array<BadgeBean>> {
    return this._store.select(this.getBadges);
  }
  
  /**
   * dispatch
   * @param action dispatch action
   */
  public dispatch(action: AllStoreActions) {
    this._store.dispatch(action);
  }
}
