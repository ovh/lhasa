import {Component, OnDestroy, OnInit} from '@angular/core';

import {ActivatedRoute} from '@angular/router';
import {ApplicationBean, DeploymentBean, EnvironmentBean} from '../../models/commons/applications-bean';
import {Store} from '@ngrx/store';
import {ApplicationsStoreService} from '../../stores/applications-store.service';
import {EnvironmentsStoreService} from '../../stores/environments-store.service';
import {element} from 'protractor';
import {ISubscription} from 'rxjs/Subscription';


@Component({
  selector: 'app-appdetail',
  templateUrl: './appdetail.component.html',
  styleUrls: ['./appdetail.component.css'],

})
export class AppdetailComponent implements OnInit, OnDestroy {

  /**
   * internal streams and store
   */
  protected applicationStream: Store<ApplicationBean>;
  protected applicationSubscription: ISubscription;
  protected deploymentStream: Store<DeploymentBean[]>;
  protected deploymentSubscription: ISubscription;
  protected application: ApplicationBean;
  protected deployments: DeploymentBean[];

  constructor(
    private applicationsStoreService: ApplicationsStoreService,
    private route: ActivatedRoute) {
    /**
     * subscribe
     */
    this.applicationStream = this.applicationsStoreService.active();
    this.deploymentStream = this.applicationsStoreService.deployments();
  }

  ngOnInit(): void {
    this.applicationSubscription = this.applicationStream.subscribe(
      (element: ApplicationBean) => {
        this.application = element;
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );

    this.deploymentSubscription = this.deploymentStream.subscribe(
      (element: DeploymentBean[]) => {
        this.deployments = element;
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }

  ngOnDestroy(): void {
    this.applicationSubscription.unsubscribe();
    this.deploymentSubscription.unsubscribe();
  }
}
