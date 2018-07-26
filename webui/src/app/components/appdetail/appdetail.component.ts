import { Component, OnInit } from '@angular/core';
import { CardModule } from 'primeng/card';
import { CodeHighlighterModule } from 'primeng/codehighlighter';
import { ActivatedRoute } from '@angular/router';
import { ApplicationBean, DeploymentBean, EnvironmentBean, BadgeRatingBean, DeploymentPropertiesBean } from '../../models/commons/applications-bean';
import { Store } from '@ngrx/store';
import { ApplicationsStoreService } from '../../stores/applications-store.service';
import { EnvironmentsStoreService } from '../../stores/environments-store.service';
import { element } from 'protractor';
import { SubscriptionLike as ISubscription } from 'rxjs';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { BadgeShieldsIOBean } from '../../widget/badgewidget/badgewidget.component';
import { FieldsetModule } from 'primeng/fieldset';
import { Observable } from 'rxjs';
import { strict } from 'assert';
import { forEach } from '@angular/router/src/utils/collection';


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
  protected applicationSubscription: ISubscription;
  public application: ApplicationBean;
  private _badgeRatingShields: BadgeShieldsIOBean[];
  public _selectedDeployment: DeploymentBean;
  public _activeDeployments: DeploymentBean[];

  public description: string;
  public readme: string;
  public links: Map<string, string>;

  constructor(
    private applicationsStoreService: ApplicationsStoreService,
    private route: ActivatedRoute) {
    /**
     * subscribe
     */
    this.applicationStream = this.applicationsStoreService.active();
  }

  ngOnInit(): void {
    this.applicationSubscription = this.applicationStream.subscribe(
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

  getKeys(myMap: Map<string, string>): Array<Object> {
    const res: Array<Object> = [];
    for (const key of Object.keys(myMap).sort()) {
      res.push({
        'label': key,
        'url': myMap[key],
      });
    }
    return res;
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
  }
}
