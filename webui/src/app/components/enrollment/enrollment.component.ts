import { Observable } from 'rxjs/Observable';
import { DeploymentBean, PersonBean } from './../../models/commons/applications-bean';
import { Component, OnInit, ChangeDetectorRef, ViewChild } from '@angular/core';
import { ApplicationsStoreService, SelectApplicationAction } from '../../stores/applications-store.service';
import { Store } from '@ngrx/store';
import { ApplicationBean } from '../../models/commons/applications-bean';
import { DataApplicationService } from '../../services/data-application-version.service';

import { FormBuilder } from '@angular/forms';

import { cloneDeep, remove } from 'lodash';
import { ActivatedRoute } from '@angular/router';
import { BitbucketService } from '../../services/data-bitbucket.service';
import { MatSnackBar } from '@angular/material';
import { ISubscription } from 'rxjs/Subscription';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { UiKitStep } from '../../models/kit/progress-tracker';
import { OuiProgressTrackerComponent } from '../../kit/oui-progress-tracker/oui-progress-tracker.component';
import { DataContentService } from '../../services/data-content.service';
import { ContentBean } from '../../models/commons/content-bean';
import { Subject } from 'rxjs/Subject';

@Component({
  selector: 'app-enrollment',
  templateUrl: './enrollment.component.html',
  styleUrls: ['./enrollment.component.css']
})
@AutoUnsubscribe()
export class EnrollmentComponent implements OnInit {

  public selected = 'description';
  public steps: UiKitStep[] = [];
  observableStep1: Observable<ContentBean>;

  @ViewChild('progress') progress: OuiProgressTrackerComponent;

  /**
   * internal streams and store
   */
  protected applicationStream: Store<ApplicationBean>;
  protected application: ApplicationBean;
  protected enrollment: string;

  /**
   * internal streams and store
   */
  protected subscription: ISubscription;

  constructor(
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService,
    private bitbucketService: BitbucketService,
    private contentService: DataContentService,
    private _formBuilder: FormBuilder,
    private activatedRoute: ActivatedRoute,
    public snackBar: MatSnackBar
  ) {
    /**
     * subscribe
     */
    this.applicationStream = this.applicationsStoreService.active();
  }

  isLinear = false;

  ngOnInit() {
    this.steps = [];
    this.steps.push({
      id: 'description',
      label: 'DESCRIPTION-MANIFEST',
      status: 'active'
    });
    this.steps.push({
      id: 'update',
      label: 'UPDATE-MANIFEST',
      status: 'disabled'
    });
    this.steps.push({
      id: 'save',
      label: 'SAVE-MANIFEST',
      status: 'disabled'
    });

    this.subscription = this.applicationStream.subscribe(
      (element: ApplicationBean) => {
        this.application = element;
        if (this.application.manifest && this.application.manifest.repository) {
          // Analyse is base on stach url
          const url = this.application.manifest.repository.split(/\//);
          this.application.project = url[4];
          this.application.repo = url[6].split(/\./)[0];
        }
        // Add support info
        if (this.application.manifest) {
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
   * dispatch load applications
   * @param event
   */
  protected createPullRequest() {
    this.application.manifest.profile = this.application.domain;
    this.application.manifest.name = this.application.name;
    // load all applications from a content return
    this.applicationsService.Task(`${this.application.domain}/${this.application.name}/assistant`, {
      repository: 'bitbucket',
      manifest: this.application.manifest
    }).subscribe(
      (data: any) => {
        this.snackBar.open('Create pull request', 'Ok', {
          duration: 2000,
        });
      },
      (data: any) => {
        console.warn(data);
        this.snackBar.open(JSON.stringify(data, null, 2), 'Error', {
          duration: 10000,
        });
      }
    );
  }

  /**
   * change selection
   */
  public onSelect(event: any) {
    this.selected = event.data.id;
    this.steps.forEach((step) => {
      if (this.selected === step.id) {
        step.status = 'active';
      } else {
        step.status = 'complete';
      }
    });

    // upgrade manifest on each change
    const manifest = cloneDeep(this.application.manifest);
    if (manifest) {
      // remove empty data
      if (manifest.support.name === '') {
        delete manifest.support.name;
      }
      if (manifest.support.email === '') {
        delete manifest.support.email;
      }
      if (manifest.support.cisco === '') {
        delete manifest.support.cisco;
      }
      if (!manifest.support.name && !manifest.support.email && !manifest.support.cisco) {
        delete manifest.support;
      }
      this.enrollment = JSON.stringify(manifest, null, 2);
    } else {
      this.application.manifest = {};
    }
  }

  /**
   * copy localfile
   */
  protected copy() {
    const a: any = document.createElement('textarea');
    document.body.appendChild(a);
    a.value = this.enrollment;
    a.select();
    document.execCommand('copy');
    a.style = 'display: none';
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
    this.application.manifest.authors.push({
        email: 'nobody@exemple.com',
        name: 'name',
        role: 'MAINTAINER',
        cisco: undefined
      }
    );
  }

  /**
   * download localfile
   */
  protected download() {
    const fileName = 'manifest.json';
    const a: any = document.createElement('a');
    document.body.appendChild(a);
    a.style = 'display: none';
    const file = new Blob([this.enrollment], {type: 'application/text'});
    const fileURL = window.URL.createObjectURL(file);
    a.href = fileURL;
    a.download = fileName;
    a.click();
  }

  /**
   * previous step
   */
  public prev() {
    this.progress.prev();
  }

  /**
   * next step
   */
  public next() {
    this.progress.next();
  }
}
