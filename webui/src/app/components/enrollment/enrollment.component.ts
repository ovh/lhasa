import {Component, OnInit} from '@angular/core';
import {ApplicationsStoreService, SelectApplicationAction} from '../../stores/applications-store.service';
import {Store} from '@ngrx/store';
import {ApplicationBean, DeploymentBean} from '../../models/commons/applications-bean';
import {DataApplicationService} from '../../services/data-application-version.service';
import {ContentListResponse} from '../../models/commons/entity-bean';

import {FormBuilder, FormGroup, Validators} from '@angular/forms';

import * as _ from 'lodash';
import {ActivatedRoute} from '@angular/router';
import {URLSearchParams} from '@angular/http';
import {BitbucketService} from '../../services/data-bitbucket.service';
import {MatSnackBar} from '@angular/material';
import {environment} from '../../../environments/environment';

@Component({
  selector: 'app-enrollment',
  templateUrl: './enrollment.component.html',
  styleUrls: ['./enrollment.component.css']
})
export class EnrollmentComponent implements OnInit {

  private selected: string = '0';

  /**
   * internal streams and store
   */
  protected applicationStream: Store<ApplicationBean>;
  protected application: ApplicationBean;
  protected enrollment: string;

  constructor(
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService,
    private bitbucketService: BitbucketService,
    private _formBuilder: FormBuilder,
    private activatedRoute: ActivatedRoute,
    public snackBar: MatSnackBar
  ) {
    /**
     * subscribe
     */
    this.applicationStream = this.applicationsStoreService.active();

    this.applicationStream.subscribe(
      (element: ApplicationBean) => {
        this.application = element;

        if (this.application.manifest && this.application.manifest.repository) {
          let url = this.application.manifest.repository.split(/\//);
          this.application.project = url[3];
          this.application.repo = url[4].split(/\./)[0];
        } else {
          this.application.project = 'UNKNOWN';
          this.application.repo = 'UNKNOWN';
        }
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );

    this.firstFormGroup = this._formBuilder.group({
      firstCtrl: ['', Validators.required]
    });
    this.secondFormGroup = this._formBuilder.group({
      secondCtrl: ['', Validators.required]
    });

    this.enrollment = `
   Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod
   tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
   quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
   Duis aute irure dolor in reprehenderit in voluptate velit esse cillum
   dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident,
   sunt in culpa qui officia deserunt mollit anim id est laborum.
`;
  }

  isLinear = false;
  firstFormGroup: FormGroup;
  secondFormGroup: FormGroup;

  ngOnInit() {
    this.activatedRoute.queryParams.subscribe(params => {
      let domain = params['domain'];
      let application = params['application'];
      if (domain && application) {
        this.selectApplication(
          domain,
          application
        );
      }
    });
  }

  /**
   * dispatch load applications
   * @param event
   */
  protected createPullRequest(domain: string, application: string) {
    this.application.manifest.profile = this.application.domain;
    this.application.manifest.name = this.application.name;
    // load all applications from a content return
    this.bitbucketService.Tasks('create-pull-request', {
      domain: this.application.domain,
      application: this.application.name,
      project: this.application.project,
      repo: this.application.repo,
      manifest: this.application.manifest
    }).subscribe(
      (data: any) => {
        this.snackBar.open('Create pull request', 'Ok', {
          duration: 2000,
        });
      }
    );
  }

  /**
   * dispatch load applications
   * @param event
   */
  protected selectApplication(domain: string, application: string) {
    // load all applications from a content return
    this.applicationsService.GetAllFromContent('/' + domain + '/' + application, new Map<string, string>([['size', '1']])).subscribe(
      (data: ContentListResponse<ApplicationBean>) => {
        this.applicationsStoreService.dispatch(
          new SelectApplicationAction(
            data.content[0],
            new Array<DeploymentBean>(),
          )
        );
      }
    );
  }

  protected resource(path: string) {
    return environment.appcatalog.baseUrlUi + '/ui' + path;
  }

  protected onSelect(event: any) {
    this.selected = event.data;
  }

  ngOnDestroy() {
    // this.sub.unsubscribe();
  }

}
