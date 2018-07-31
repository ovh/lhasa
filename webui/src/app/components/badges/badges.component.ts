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
import { BehaviorSubject } from 'rxjs';
import { OuiPaginationComponent } from '../../kit/oui-pagination/oui-pagination.component';
import { BadgeShieldsIOBean } from '../../widget/badgewidget/badgewidget.component';
import { UIChart } from 'primeng/primeng';
import { Observable } from 'rxjs';

export class BadgeUIBean {
  slug: string;
  title: string;
  type: string;
  levels: BadgeShieldsIOBean[];
  piechartData: PieChartDataBean;
}

export class PieChartDataBean {
  labels: string[] = [];
  datasets: PieChartDatasetBean[] = [];
}

export class PieChartDatasetBean {
  label: string;
  data: number[] = [];
  backgroundColor: string[] = [];
  levelID: string[] = [];
}

@Component({
  selector: 'app-badges',
  templateUrl: './badges.component.html',
  styleUrls: ['./badges.component.css'],
})
export class BadgesComponent implements OnInit {

  @ViewChild('paginationtop') paginationtop: OuiPaginationComponent;
  @ViewChild('paginationbottom') paginationbottom: OuiPaginationComponent;

  /**
   * internal streams and store
   */
  protected badgesStream: Observable<BadgePagesBean>;
  public badges: BadgeUIBean[] = [];
  public metadata: UiKitPaginate = {
    totalElements: 0,
    totalPages: 0,
    size: 0,
    number: 0
  };

  public param = { target: 'badgess' };
  public page = 0;
  public PieChartOptions;

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

    this.PieChartOptions = {
      cutoutPercentage: 50,
      legend: {
        display: false
      },
      animation: {
        animateScale: true,
      }

    };
  }

  onPieChartClick(event) {
    const dataset = event.dataset[event.element._index]._chart.data.datasets[0];
    const params = {};
    params[`badge_ratings.${dataset.label}.value`] = dataset.levelID[event.element._index];
    this.router.navigate(['/applications'],  { queryParams: params} );
  }

  ngOnInit() {
    /**
     * subscribe
     */
    this.badgesStream = this.badgesStoreService.badgePages();

    this.badgesStream.subscribe(
      (elem: BadgePagesBean) => {
        this.badges = [];
        elem.badges.forEach((bdg) => {
          const levels: BadgeShieldsIOBean[] = [];
          const piechartData = new (PieChartDataBean);
          piechartData.datasets.push(new (PieChartDatasetBean));
          piechartData.datasets[0].label = bdg.slug;
          bdg.levels.forEach((lvl) => {
            levels.push(<BadgeShieldsIOBean>{
              id: lvl.id,
              title: bdg.title,
              label: lvl.label,
              color: lvl.color,
              comment: '-',
              description: lvl.description,
              value: lvl.id,
            });
            piechartData.labels.push(`${lvl.label} (${lvl.id})`);
            let val = bdg._stats[lvl.id];
            if (val === undefined) {
              val = 0;
            }
            piechartData.datasets[0].data.push(val);
            piechartData.datasets[0].backgroundColor.push(lvl.color);
            piechartData.datasets[0].levelID.push(lvl.id);
          });
          this.badges.push(<BadgeUIBean>{
            slug: bdg.slug,
            title: bdg.title,
            type: bdg.type,
            levels: levels,
            piechartData: piechartData,
          });
        });

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
      this.paginationtop.RefreshMetadata(this.metadata, 'select', this.page);
      this.paginationbottom.RefreshMetadata(this.metadata, 'select', this.page);
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
