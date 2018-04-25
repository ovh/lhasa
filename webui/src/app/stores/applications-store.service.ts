import { Injectable } from '@angular/core';
import { ActionReducer, Action, State } from '@ngrx/store';
import { Store } from '@ngrx/store';
import { createFeatureSelector, createSelector, MemoizedSelector } from '@ngrx/store';

import * as _ from 'lodash';

import { ActionWithPayload } from './action-with-payload';
import { ApplicationBean, DomainBean } from '../models/commons/applications-bean';

/**
 * states
 */
export interface AppState {
  feature: ApplicationState;
}

export interface ApplicationState {
  domains: Array<DomainBean>;
  applications: Array<ApplicationBean>;
  active: ApplicationBean;
}

/**
 * actions
 */
export class LoadApplicationsAction implements ActionWithPayload<Array<ApplicationBean>> {
  readonly type = 'LoadApplicationsAction';
  constructor(public payload: Array<ApplicationBean>) { }
}

export class SelectApplicationAction implements ActionWithPayload<ApplicationBean> {
  readonly type = 'SelectApplicationAction';
  constructor(public payload: ApplicationBean) { }
}

export type AllStoreActions = LoadApplicationsAction | SelectApplicationAction;

/**
 * main store for this application
 */
@Injectable()
export class ApplicationsStoreService {

  private getDomains: MemoizedSelector<object, Array<DomainBean>>;
  private getApplications: MemoizedSelector<object, Array<ApplicationBean>>;
  private getActive: MemoizedSelector<object, ApplicationBean>;

  /**
   * 
   * @param _store constructor
   */
  constructor(
    private _store: Store<ApplicationState>
  ) {
    this.getDomains = createSelector(createFeatureSelector<ApplicationState>('applications'), (state: ApplicationState) => state.domains);
    this.getApplications = createSelector(createFeatureSelector<ApplicationState>('applications'), (state: ApplicationState) => state.applications);
    this.getActive = createSelector(createFeatureSelector<ApplicationState>('applications'), (state: ApplicationState) => state.active);
  }

  /**
   * select this store service
   */
  public domains(): Store<Array<DomainBean>> {
    return this._store.select(this.getDomains);
  }

  /**
   * select this store service
   */
  public applications(): Store<Array<ApplicationBean>> {
    return this._store.select(this.getApplications);
  }

  /**
   * select this store service
   */
  public active(): Store<ApplicationBean> {
    return this._store.select(this.getActive);
  }

  /**
   * dispatch
   * @param action dispatch action
   */
  public dispatch(action: AllStoreActions) {
    this._store.dispatch(action);
  }

  /**
   * metareducer (Cf. https://www.concretepage.com/angular-2/ngrx/ngrx-store-4-angular-5-tutorial)
   * @param state 
   * @param action 
   */
  public static reducer(state: ApplicationState = { domains: new Array<DomainBean>(), applications: new Array<ApplicationBean>(), active: new ApplicationBean() }, action: AllStoreActions): ApplicationState {

    switch (action.type) {
      /**
       * message incomming
       */
      case 'LoadApplicationsAction':
        {
          let applications = Object.assign([], action.payload);

          let orderedDomains = new Map<string, ApplicationBean[]>()
          _.each(applications, (app) => {
            if (!orderedDomains.has(app.domain)) {
              orderedDomains.set(app.domain, [])
            }
            orderedDomains.get(app.domain).push(app)
          });
          let domains = []
          orderedDomains.forEach((v, k) => {
            domains.push({ name: k, applications: v })
          })

          return {
            domains: domains,
            applications: applications,
            active: applications[0]
          };
        }

      case 'SelectApplicationAction':
        {
          let active = Object.assign({}, action.payload);
          return {
            domains: state.domains,
            applications: state.applications,
            active: active
          };
        }

      default:
        return state;
    }
  }
}
