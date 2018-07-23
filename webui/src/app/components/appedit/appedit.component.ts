import { Observable ,  SubscriptionLike as ISubscription } from 'rxjs';
import { PersonBean } from './../../models/commons/applications-bean';
import { Component, OnInit, ViewChild } from '@angular/core';
import { ApplicationsStoreService } from '../../stores/applications-store.service';
import { Store } from '@ngrx/store';
import { ApplicationBean } from '../../models/commons/applications-bean';
import { DataApplicationService } from '../../services/data-application-version.service';

import { cloneDeep, remove } from 'lodash';
import { MatSnackBar } from '@angular/material';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { UiKitStep } from '../../models/kit/progress-tracker';
import { OuiProgressTrackerComponent } from '../../kit/oui-progress-tracker/oui-progress-tracker.component';
import { DataContentService } from '../../services/data-content.service';
import { ContentBean } from '../../models/commons/content-bean';

import { Router, CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';


@Component({
  selector: 'app-edit',
  templateUrl: './appedit.component.html',
  styleUrls: []
})
@AutoUnsubscribe()
export class AppEditComponent implements OnInit {

  public selected = 'description';
  public steps: UiKitStep[] = [];
  observableStep1: Observable<ContentBean>;

  @ViewChild('progress') progress: OuiProgressTrackerComponent;

  /**
   * internal streams and store
   */
  protected applicationStream: Observable<ApplicationBean>;
  public application: ApplicationBean;

  /**
   * internal streams and store
   */
  protected subscription: ISubscription;

  constructor(
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService,
    private contentService: DataContentService,
    private router: Router,
    public snackBar: MatSnackBar
  ) {
    /**
     * subscribe
     */
    this.applicationStream = this.applicationsStoreService.active();
  }

  ngOnInit() {

    this.subscription = this.applicationStream.subscribe(
      (element: ApplicationBean) => {
        this.application = element;
        if (this.application.manifest && this.application.manifest.repository) {
          // Analysis is based on the url of the repository
          const url = this.application.manifest.repository.split(/\//);
          this.application.project = url[4];
          this.application.repo = url[6].split(/\./)[0];
        }
        // Add support info
        if (this.application.manifest.support == undefined ){
          this.application.manifest.support = {
            name: '',
            email: '',
            cisco: ''
          };
        }
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );

    /**
     * resolve resource
     */
    this.observableStep1 = this.contentService.GetSingle('register-page');
  }

  /**
   * drop
   */
  protected drop(author: PersonBean) {
    remove(this.application.manifest.authors, (item: PersonBean) => {
      return item.email === author.email;
    });
  }

  /**
   * add
   */
  protected add(author: PersonBean) {
    if (this.application.manifest.authors == null){
      this.application.manifest.authors = []  
    }
    
    
    this.application.manifest.authors.push({
        email: null,
        name: undefined,
        role: 'MAINTAINER',
        cisco: undefined
      }
    );
  }


  /**
   * save
   */
  public save() {
    this.application.manifest.profile = this.application.domain;
    this.application.manifest.name = this.application.name;
    // load all applications from a content return
    this.applicationsService.Update(`${this.application.domain}/${this.application.name}/versions/${this.application.version}`, 
    this.application).subscribe(
      (data: any) => {
        this.snackBar.open('Application Saved', 'Ok', {
          duration: 2000,
        });
        // redirect route qui va rafraichir le store
        this.router.navigateByUrl(`/applications/${this.application.domain}/${this.application.name}/${this.application.version}`)
      },
      (data: any) => {
        console.warn(data);
        this.snackBar.open(JSON.stringify(data, null, 2), 'Error', {
          duration: 10000,
        });
      }
    );
  }
}
