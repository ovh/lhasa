import { Component, OnInit, ViewChild } from '@angular/core';
import { ApplicationBean, DeploymentBean } from '../../models/commons/applications-bean';
import { ApplicationsStoreService } from '../../stores/applications-store.service';
import { Observable } from 'rxjs';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { BadgeShieldsIOBean } from '../../widget/badgewidget/badgewidget.component';
import { GraphComponent } from '../../widget/graph/graph.component';
import { GraphsStoreService } from '../../stores/graphs-store.service';
import { GraphBean, GraphVis, NodeBean } from '../../models/graph/graph-bean';


@Component({
  selector: 'app-appdetail',
  templateUrl: './appdetail.component.html',
  styleUrls: ['./appdetail.component.css'],
})

@AutoUnsubscribe()
export class AppdetailComponent implements OnInit {

  /**
   * internal streams and store
   */
  protected applicationStream: Observable<ApplicationBean>;
  public application: ApplicationBean;
  private _badgeRatingShields: BadgeShieldsIOBean[];
  public _selectedDeployment: DeploymentBean;
  public _activeDeployments: DeploymentBean[];

  public description: string;
  public readme: string;
  public links: Map<string, string>;

  @ViewChild('deploymentsGraph')
  private graph: GraphComponent;

  constructor(
    private applicationsStoreService: ApplicationsStoreService) {
    /**
     * subscribe
     */
    this.applicationStream = this.applicationsStoreService.active();
  }

  ngOnInit(): void {
    this.applicationStream.subscribe(
      (app: ApplicationBean) => {
        // With full text search, active application can be undefined
        // Just check if app is undefined
        if (!app) {
          return;
        }
        this.application = app;

        // Rule is
        // First use properties.description
        // Then manifest.description overrive this description
        if (app.properties && app.properties.description) {
          this.description = app.properties.description;
        }
        if (app.manifest && app.manifest.description) {
          this.description = app.manifest.description;
        }
        if (app.properties && app.properties.readme) {
          this.readme = app.properties.readme;
        }
        this._activeDeployments = [];
        this._selectedDeployment = null;
        if (this.application.deployments !== undefined) {
          this.application.deployments.forEach(deployment => {
            if (!deployment.undeployedAt || deployment.undeployedAt === null) {
              this._activeDeployments.push(deployment);
              if (this._selectedDeployment === null) {
                this._selectedDeployment = deployment;
                if (this._selectedDeployment.properties.links == null ) {
                  this._selectedDeployment.properties.links = new Map<string, string>();
                }
                console.log(deployment);
              }
            }
          });
        }
        if (this.application.badgeRatings === undefined) {
          this.application.badgeRatings = [];
        }
        this._badgeRatingShields = [];
        app.badgeRatings.forEach((bdgRating) => {
          this._badgeRatingShields.push(<BadgeShieldsIOBean>{
              id: bdgRating.badgeslug,
              value: bdgRating.value,
              title: bdgRating.badgetitle,
              comment: bdgRating.comment,
              label: bdgRating.level.label,
              color: bdgRating.level.color,
              description: bdgRating.level.description,
            }
          );
        });
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }

  shortLink(str: string, maxLen: number): string {
    str = str.replace('https://', '').replace('http://', '').replace('mailto:', '');
    if (str.length > maxLen) {
      return str.substr(0, maxLen / 2) + '…' + str.substr(str.length - maxLen / 2, str.length);
    }
    return str;
  }

  selectDeployment(deployment: DeploymentBean): void {
    this._selectedDeployment = deployment;
    this.graph.graph = deployment._graph;
    this.graph.update();
  }
}
