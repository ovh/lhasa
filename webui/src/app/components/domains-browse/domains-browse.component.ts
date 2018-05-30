import { Router } from '@angular/router';
import { Component, OnInit, ViewChild } from '@angular/core';
import { ApplicationsStoreService, LoadApplicationsAction, SelectApplicationAction } from '../../stores/applications-store.service';
import { Store } from '@ngrx/store';
import { ApplicationBean, DeploymentBean, DomainBean, DomainPagesBean } from '../../models/commons/applications-bean';
import { DataApplicationService } from '../../services/data-application-version.service';
import { ContentListResponse, PageMetaData } from '../../models/commons/entity-bean';
import { find } from 'lodash';

import { DataDeploymentService } from '../../services/data-deployment.service';
import { UiKitPaginate } from '../../models/kit/paginate';
import { ActivatedRoute } from '@angular/router';
import { DomainsResolver } from '../../resolver/resolve-domains';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { OuiPaginationComponent } from '../../kit/oui-pagination/oui-pagination.component';

@Component({
  selector: 'app-domains-browse',
  templateUrl: './domains-browse.component.html',
  styleUrls: ['./domains-browse.component.css'],
})
export class DomainsBrowseComponent implements OnInit {

  @ViewChild('pagination') pagination: OuiPaginationComponent;

  /**
   * internal streams and store
   */
  protected domainsStream: Store<DomainPagesBean>;
  public domains: DomainBean[] = [];
  public metadata: UiKitPaginate = {
    totalElements: 0,
    totalPages: 0,
    size: 0,
    number: 0
  };

  public param = { target: 'applications' };
  public page = 0;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private applicationsStoreService: ApplicationsStoreService,
    private domainsResolver: DomainsResolver
  ) {
    // Subscribe to retrieve page asked
    this.route
      .queryParams
      .subscribe(params => {
        // Defaults to 0 if no query param provided.
        this.page = +params['page'] || 0;
      });

  }

  ngOnInit() {
    /**
     * subscribe
     */
    this.domainsStream = this.applicationsStoreService.domainPages();

    this.domainsStream.subscribe(
      (elem: DomainPagesBean) => {
        this.domains = elem.domains;
        this.metadata = {
          totalElements: elem.metadata.totalElements,
          totalPages: elem.metadata.totalPages,
          size: elem.metadata.size,
          number: elem.metadata.number
        };
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
    // if page different from 0
    if (this.page !== 0) {
      this.pagination.RefreshMetadata(this.metadata, 'select', this.page);
    }
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
      this.domainsResolver.selectDomains(metadata, new BehaviorSubject<any>('select another page on domains'));
      this.page = page;
    }
  }

}
