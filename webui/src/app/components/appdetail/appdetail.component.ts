import { Component, OnInit } from '@angular/core';
import { CardModule } from 'primeng/card';
import { CodeHighlighterModule } from 'primeng/codehighlighter';
import { ActivatedRoute } from '@angular/router';
import { ApplicationBean, DeploymentBean, EnvironmentBean, BadgeRatingBean } from '../../models/commons/applications-bean';
import { Store } from '@ngrx/store';
import { ApplicationsStoreService } from '../../stores/applications-store.service';
import { EnvironmentsStoreService } from '../../stores/environments-store.service';
import { element } from 'protractor';
import { SubscriptionLike as ISubscription } from 'rxjs';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { BadgeShieldsIOBean } from '../../widget/badgewidget/badgewidget.component';
import {FieldsetModule} from 'primeng/fieldset';
import { Observable } from 'rxjs';


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

        // When read field exist use it as plain text
        // Or if it's an url call it to obtain data

        this._activeDeployments = [];
        this._selectedDeployment = null;
        if (this.application.deployments !== undefined) {
          this.application.deployments.forEach(value => {
            if (!value.undeployedAt || value.undeployedAt === null) {
              this._activeDeployments.push(value);
              if (this._selectedDeployment === null) {
                this._selectedDeployment = value;
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

  selectDeployment(deployment: DeploymentBean): void {
    this._selectedDeployment = deployment;
  }
}
