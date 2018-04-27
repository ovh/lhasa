import {Injectable} from '@angular/core';
import {createFeatureSelector, createSelector, MemoizedSelector, Store} from '@ngrx/store';

import {ActionWithPayload} from './action-with-payload';
import {ApplicationBean, DeploymentBean, DomainBean} from '../models/commons/applications-bean';

/**
 * states
 */
export interface ApplicationState {
  domains: Array<DomainBean>;
  applications: Array<ApplicationBean>;
  active: ApplicationBean;
  deployments: Array<DeploymentBean>;
}

/**
 * actions
 */
export class LoadApplicationsAction implements ActionWithPayload<Array<ApplicationBean>> {
  readonly type = LoadApplicationsAction.getType();

  public static getType(): string {
    return 'LoadApplicationsAction';
  }

  constructor(public payload: Array<ApplicationBean>) {
  }
}

export class SelectApplicationAction implements ActionWithPayload<ApplicationBean> {
  readonly type = SelectApplicationAction.getType();

  public static getType(): string {
    return 'SelectApplicationAction';
  }

  constructor(public payload: ApplicationBean, public deployments: Array<DeploymentBean>) {
  }
}

export type AllStoreActions = LoadApplicationsAction | SelectApplicationAction;

/**
 * main store for this application
 */
@Injectable()
export class ApplicationsStoreService {

  readonly getDomains: MemoizedSelector<object, Array<DomainBean>>;
  readonly getApplications: MemoizedSelector<object, Array<ApplicationBean>>;
  readonly getActive: MemoizedSelector<object, ApplicationBean>;
  readonly getDeployments: MemoizedSelector<object, Array<DeploymentBean>>;

  /**
   *
   * @param _store constructor
   */
  constructor(
    private _store: Store<ApplicationState>
  ) {
    this.getDomains = createSelector(createFeatureSelector<ApplicationState>('applications'), (state: ApplicationState) => state.domains);
    this.getApplications = createSelector(createFeatureSelector<ApplicationState>('applications'),
      (state: ApplicationState) => state.applications);
    this.getActive = createSelector(createFeatureSelector<ApplicationState>('applications'), (state: ApplicationState) => state.active);
    this.getDeployments = createSelector(createFeatureSelector<ApplicationState>('applications'),
      (state: ApplicationState) => state.deployments);
  }

  /**
   * metareducer (Cf. https://www.concretepage.com/angular-2/ngrx/ngrx-store-4-angular-5-tutorial)
   * @param state
   * @param action
   */
  public static reducer(state: ApplicationState = {
    domains: new Array<DomainBean>(),
    applications: new Array<ApplicationBean>(),
    active: new ApplicationBean(),
    deployments: new Array<DeploymentBean>(),
  }, action: AllStoreActions): ApplicationState {

    switch (action.type) {
      /**
       * message incomming
       */
      case LoadApplicationsAction.getType(): {
        const applications = Object.assign([], action.payload);

        const orderedDomains = new Map<string, ApplicationBean[]>();
        applications.forEach((app) => {
          if (!orderedDomains.has(app.domain)) {
            orderedDomains.set(app.domain, []);
          }
          orderedDomains.get(app.domain).push(app);
        });
        const domains = [];
        orderedDomains.forEach((v, k) => {
          domains.push({name: k, applications: v});
        });

        return {
          domains: domains,
          applications: applications,
          active: applications[0],
          deployments: new Array<DeploymentBean>(),
        };
      }

      case SelectApplicationAction.getType(): {
        action = action as SelectApplicationAction;
        const active = Object.assign({}, action.payload);

        const deployments = Object.assign([], action.deployments);
        return {
          domains: state.domains,
          applications: state.applications,
          active: active,
          deployments: deployments,
        };
      }

      default:
        return state;
    }
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
   * select this store service
   */
  public deployments(): Store<Array<DeploymentBean>> {
    return this._store.select(this.getDeployments);
  }

  /**
   * dispatch
   * @param action dispatch action
   */
  public dispatch(action: AllStoreActions) {
    this._store.dispatch(action);
  }
}
