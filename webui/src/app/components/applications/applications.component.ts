import { Component, OnInit, Pipe, ViewChild } from '@angular/core';
import { ApplicationsStoreService, SelectApplicationAction, LoadApplicationsAction } from '../../stores/applications-store.service';
import { Store } from '@ngrx/store';
import { ApplicationBean, DeploymentBean, DomainBean, ApplicationPagesBean } from '../../models/commons/applications-bean';
import { DataApplicationService } from '../../services/data-application-version.service';
import { ContentListResponse, PageMetaData } from '../../models/commons/entity-bean';

import { DataDeploymentService } from '../../services/data-deployment.service';
import { UiKitPaginate } from '../../models/kit/paginate';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { ApplicationsResolver } from '../../resolver/resolve-applications';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { OuiMessageComponent } from '../../kit/oui-message/oui-message.component';
import { OuiPaginationComponent } from '../../kit/oui-pagination/oui-pagination.component';

@Component({
  selector: 'app-applications',
  templateUrl: './applications.component.html',
  styleUrls: ['./applications.component.css'],
})
export class ApplicationsComponent implements OnInit {

  /**
   * internal streams and store
   */
  protected applicationsStream: Store<ApplicationPagesBean>;
  public applications: ApplicationBean[];
  public metadata: UiKitPaginate = {
    totalElements: 0,
    totalPages: 0,
    size: 0,
    number: 0
  };

  @ViewChild('paginationtop') paginationtop: OuiPaginationComponent;
  @ViewChild('paginationbottom') paginationbottom: OuiPaginationComponent;
  @ViewChild('message') msg: OuiMessageComponent;


  public domain = '';
  public param = { target: '' };
  public searchString = ''

  protected orderedDomains = new Map<string, ApplicationBean[]>();
  public domains: DomainBean[] = [];
  public page = 0;

  constructor(
    private router: Router,
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsResolver: ApplicationsResolver,
    private route: ActivatedRoute
  ) {
    /**
     * subscribe
     */
    this.applicationsStream = this.applicationsStoreService.applications();

    // Subscribe to retrieve page asked
    this.route
      .queryParams
      .subscribe(params => {
        let changeDomain = false;
        // Defaults to 0 if no query param provided.
        this.page = +params['page'] || 0;
        this.searchString = params['search'] || '';
        // verify if filter by domain
        if (!params['domain'] && this.domain !== '') {
          changeDomain = true;
        }
        if (params['domain']) {
          this.domain = '/' + params['domain'];
          this.param.target = params['domain'];
        } else {
          this.domain = '';
          this.param.target = null;
        }
        // reload data
        if (changeDomain) {
          this.refreshApplications();
        }
      });
  }

  ngOnInit() {
    // store
    this.applicationsStream.subscribe(
      (element: ApplicationPagesBean) => {
        this.applications = element.applications;
        this.applications.forEach((app) => {
          app.description = 'No description provided. Please fill the `description` field of the manifest.';
          if (app.manifest && app.manifest.description) {
            app.description = (app.manifest.description.length > 200) ?
              (app.manifest.description.substr(0, 200) + '...') : (app.manifest.description);
          }
        });
        this.orderedDomains = new Map<string, ApplicationBean[]>();
        this.applications.forEach((app) => {
          if (!this.orderedDomains.has(app.domain)) {
            this.orderedDomains.set(app.domain, []);
          }
          this.orderedDomains.get(app.domain).push(app);
        });
        this.domains = [];
        this.orderedDomains.forEach((v, k) => {
          this.domains.push({ id: null, timestamp: null, name: k, applications: v });
        });
        // Fix meta data
        this.metadata = {
          totalElements: element.metadata.totalElements,
          totalPages: element.metadata.totalPages,
          size: element.metadata.size,
          number: element.metadata.number
        };
        // if page different from 0
        if (this.page !== 0) {
          this.paginationtop.RefreshMetadata(this.metadata, 'select', this.page);
          this.paginationbottom.RefreshMetadata(this.metadata, 'select', this.page);
        }
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }

  /**
   * change selection
   */
  public onSelect(event: any) {
    const metadata: PageMetaData = {
      totalElements: event.data.metadata.totalElements,
      totalPages: event.data.metadata.totalPages,
      size: event.data.metadata.size,
      number: event.data.page
    };
    // Change page
    this.navigate(metadata, event.data.page);
  }

  /**
   * change navigation
   */
  public navigate(metadata: PageMetaData, page: number) {
    // navigate if needed
    if (page !== this.page) {
      // Refresh query params
      this.router.navigate([], { queryParams: { page: page } });
      this.applicationsResolver.selectApplications(metadata,
        this.domain, this.searchString,
        new BehaviorSubject<any>('select another page on domains'));
      this.page = page;
    }
  }

  /**
   * change selection
   */
  public onMessageEvent(event: any) {
    if (event.data.type === 'close') {
      this.msg.hide();
      this.refreshApplications();
    }
  }

  /**
   * refresh applications
   */
  public refreshApplications() {
    // Reset nav
    this.router.navigate([], { queryParams: { page: 0, search: this.searchString } });
    return this.applicationsResolver.selectApplications({
      number: 0,
      size: 100
    }, '', this.searchString, new BehaviorSubject<any>('refresh all applications'));
  }

  public onFilterKeyPressed(event: any) {
    this.searchString = event.target.value;
    this.refreshApplications()
  }
}
