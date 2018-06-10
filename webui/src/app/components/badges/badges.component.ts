import { Router } from '@angular/router';
import { Component, OnInit, ViewChild, PipeTransform, Pipe } from '@angular/core';
import { BadgesStoreService } from '../../stores/badges-store.service';
import { Store } from '@ngrx/store';
import { BadgeBean, BadgePagesBean } from '../../models/commons/badges-bean';
import { DataApplicationService } from '../../services/data-application-version.service';
import { ContentListResponse, PageMetaData } from '../../models/commons/entity-bean';
import { find } from 'lodash';

import { DataDeploymentService } from '../../services/data-deployment.service';
import { UiKitPaginate } from '../../models/kit/paginate';
import { ActivatedRoute } from '@angular/router';
import { BadgesResolver } from '../../resolver/resolve-badges';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { OuiPaginationComponent } from '../../kit/oui-pagination/oui-pagination.component';
import { BadgeShieldsIOBean } from '../../widget/badgewidget/badgewidget.component';

export class BadgeUIBean {
  title: string;
  type: string;
  levels: BadgeShieldsIOBean[];
}

@Component({
  selector: 'app-badges',
  templateUrl: './badges.component.html',
  styleUrls: ['./badges.component.css'],
})
export class BadgesComponent implements OnInit {

  @ViewChild('pagination') pagination: OuiPaginationComponent;

  /**
   * internal streams and store
   */
  protected badgesStream: Store<BadgePagesBean>;
  public badges: BadgeUIBean[] = [];
  public metadata: UiKitPaginate = {
    totalElements: 0,
    totalPages: 0,
    size: 0,
    number: 0
  };

  public param = { target: 'badgess' };
  public page = 0;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private badgesStoreService: BadgesStoreService,
    private badgesResolver: BadgesResolver
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
    this.badgesStream = this.badgesStoreService.badgePages();
    
    this.badgesStream.subscribe(  
      (elem: BadgePagesBean) => {
        this.badges = []
        elem.badges.forEach((bdg) => {
          var levels: BadgeShieldsIOBean[] = []
          bdg.levels.forEach((lvl) => {
            levels.push(<BadgeShieldsIOBean>{
              id: lvl.id,
              title: bdg.title,
              label: lvl.label,
              color: lvl.color,
              comment: "-",
              description: lvl.description,
              value: lvl.id,
            })
          })
          this.badges.push(<BadgeUIBean>{
            title: bdg.title,
            type: bdg.type,
            levels: levels,
          });
        })

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
      this.badgesResolver.selectBadges(metadata, new BehaviorSubject<any>('select another page on badges'));
      this.page = page;
    }
  }
  
}
