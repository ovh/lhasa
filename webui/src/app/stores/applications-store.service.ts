import {Injectable} from '@angular/core';
import {createFeatureSelector, createSelector, Store, Selector} from '@ngrx/store';

import { ActionWithPayload, ActionWithPayloadAndPromise } from './action-with-payload';
import {ApplicationBean, DeploymentBean, DomainBean, DomainPagesBean, ApplicationPagesBean} from '../models/commons/applications-bean';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';

/**
 * states
 */
export interface ApplicationState {
  /**
   * domains of each application loaded in store
   */
  domainPages: DomainPagesBean;
  /**
   * domains of each application loaded in store
   */
  domains: Array<DomainBean>;
  applications: ApplicationPagesBean;
  active: ApplicationBean;
  deployments: Array<DeploymentBean>;
}

/**
 * actions
 */
export class LoadApplicationsAction implements ActionWithPayloadAndPromise<ApplicationPagesBean> {
  readonly type = LoadApplicationsAction.getType();

  public static getType(): string {
    return 'LoadApplicationsAction';
  }

  constructor(public payload: ApplicationPagesBean, public subject?: Subject<any>) {
  }
}

export class SelectApplicationAction implements ActionWithPayloadAndPromise<ApplicationBean> {
  readonly type = SelectApplicationAction.getType();

  public static getType(): string {
    return 'SelectApplicationAction';
  }

  constructor(public payload: ApplicationBean, public deployments: Array<DeploymentBean>, public subject?: Subject<any>) {
  }
}

export class LoadDomainsAction implements ActionWithPayloadAndPromise<DomainPagesBean> {
  readonly type = LoadDomainsAction.getType();

  public static getType(): string {
    return 'LoadDomainsAction';
  }

  constructor(public payload: DomainPagesBean, public subject: Subject<any>) {
  }
}

export type AllStoreActions = LoadApplicationsAction | LoadDomainsAction | SelectApplicationAction;

/**
 * main store for this application
 */
@Injectable()
export class ApplicationsStoreService {

  readonly getDomainPages: Selector<object, DomainPagesBean>;
  readonly getDomains: Selector<object, Array<DomainBean>>;
  readonly getApplications: Selector<object, ApplicationPagesBean>;
  readonly getActive: Selector<object, ApplicationBean>;
  readonly getDeployments: Selector<object, Array<DeploymentBean>>;

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
    this.getDomainPages = ApplicationsStoreService.create((state: ApplicationState) => state.domainPages);
  }

  /**
   * create a selector
   * @param handler internal static
   */
  private static create(handler: (S1: ApplicationState) => any)  {
    return createSelector(createFeatureSelector<ApplicationState>('applications'), handler);
  }

  /**
   * metareducer (Cf. https://www.concretepage.com/angular-2/ngrx/ngrx-store-4-angular-5-tutorial)
   * @param state
   * @param action
   */
  public static reducer(state: ApplicationState = {
    domainPages: new DomainPagesBean(),
    domains: new Array<DomainBean>(),
    applications: new ApplicationPagesBean(),
    active: new ApplicationBean(),
    deployments: new Array<DeploymentBean>(),
  }, action: AllStoreActions): ApplicationState {

    switch (action.type) {
      /**
       * update all applications in store
       */
      case LoadApplicationsAction.getType(): {
        const pages = <ApplicationPagesBean> Object.assign([], action.payload);

        const orderedDomains = new Map<string, ApplicationBean[]>();
        pages.applications.forEach((app) => {
          if (!orderedDomains.has(app.domain)) {
            orderedDomains.set(app.domain, []);
          }
          orderedDomains.get(app.domain).push(app);
        });
        const domains = [];
        orderedDomains.forEach((v, k) => {
          domains.push({name: k, applications: v});
        });

        action.subject.complete();
        return {
          domainPages: state.domainPages,
          domains: domains,
          applications: pages,
          active: pages.applications[0],
          deployments: new Array<DeploymentBean>(),
        };
      }

      /**
       * update all domains in store
       */
      case LoadDomainsAction.getType(): {
        const domainPages = Object.assign(new DomainPagesBean(), action.payload);

        /**
         * notify domains change
         */
        action.subject.complete();
        return {
          domainPages: domainPages,
          domains: state.domains,
          applications: state.applications,
          active: state.active,
          deployments: state.deployments,
        };
      }

      /**
       * select a single applications
       */
      case SelectApplicationAction.getType(): {
        action = action as SelectApplicationAction;
        const active = Object.assign({}, action.payload);

        const deployments = Object.assign([], action.deployments);
        return {
          domainPages: state.domainPages,
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
  public domainPages(): Store<DomainPagesBean> {
    return this._store.select(this.getDomainPages);
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
  public applications(): Store<ApplicationPagesBean> {
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
