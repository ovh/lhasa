import { Component, OnInit, Pipe, ViewChild } from '@angular/core';
import { ApplicationsStoreService, SelectApplicationAction, LoadApplicationsAction } from '../../stores/applications-store.service';
import { Store } from '@ngrx/store';
import { ApplicationBean, DeploymentBean, DomainBean, ApplicationPagesBean } from '../../models/commons/applications-bean';
import { DataApplicationService } from '../../services/data-application-version.service';
import { ContentListResponse, PageMetaData } from '../../models/commons/entity-bean';

import { DataDeploymentService } from '../../services/data-deployment.service';
import { UiKitPaginate } from '../../models/kit/paginate';
import { BehaviorSubject } from 'rxjs';
import { ApplicationsResolver } from '../../resolver/resolve-applications';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { OuiMessageComponent } from '../../kit/oui-message/oui-message.component';
import { OuiPaginationComponent } from '../../kit/oui-pagination/oui-pagination.component';
import { Observable, Subject } from 'rxjs';
import { Event } from '@angular/router/src/events';

@Component({
  selector: 'app-applications',
  templateUrl: './applications.component.html',
  styleUrls: ['./applications.component.css'],
})
export class ApplicationsComponent implements OnInit {

  /**
   * internal streams and store
   */
  protected applicationsStream: Observable<ApplicationPagesBean>;
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
  public searchString = '';
  public keyUp = new Subject<Event>();

  protected orderedDomains = new Map<string, ApplicationBean[]>();
  public domains: DomainBean[] = [];
  public page = 0;

  constructor(
    private router: Router,
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsResolver: ApplicationsResolver,
    private route: ActivatedRoute
  ) {

    const subscription = this.keyUp
      .map(event => event['target'].value)
      .debounceTime(1000)
      .distinctUntilChanged()
      .flatMap(search => Observable.of(search).delay(500))
      .subscribe((searchString: string) => {
        this.searchString = searchString;
        this.refreshApplications();
      });

    /**
     * subscribe
     */
    this.applicationsStream = this.applicationsStoreService.applications();

    // Subscribe to retrieve page asked
    this.route
      .queryParams
      .subscribe(params => {
        if (this.searchString !== '') {
          return;
        }
        for (const k of Object.keys(params)) {
          if (k === 'page') {
            this.page = +params[k];
            continue;
          }
          if (k === 'freesearch') {
            this.searchString += `${params[k]} `;
            continue;
          }
          this.searchString += `${k}:${params[k]} `;
        }
        this.searchString.trim();
      });
  }

  ngOnInit() {
    // store
    this.applicationsStream.subscribe(
      (element: ApplicationPagesBean) => {
        this.applications = element.applications;
        this.applications.forEach((app) => {
          app.description = 'No description.';
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
    if (this.page !== event.data.page) {
      this.page = event.data.page;
      this.refreshApplications();
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
    const params = {};
    this.searchString.split(' ').forEach((part: string) => {
      if (part === '') { return; }
      const blocks = part.split(':');
      if (blocks.length === 2 && blocks[0] !== '' && blocks[1] !== '') {
        params[blocks[0]] = blocks[1];
        return;
      }
      params['freesearch'] = part;
    });
    params['page'] = this.page;
    this.router.navigate([], { queryParams: params});
    return this.applicationsResolver.selectApplications(params, new BehaviorSubject<any>('refresh all applications'));
  }
}
