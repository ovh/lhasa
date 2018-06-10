import { Component, OnInit } from '@angular/core';

import { ActivatedRoute } from '@angular/router';
import { ApplicationBean, DeploymentBean, EnvironmentBean, BadgeRatingBean } from '../../models/commons/applications-bean';
import { Store } from '@ngrx/store';
import { ApplicationsStoreService } from '../../stores/applications-store.service';
import { EnvironmentsStoreService } from '../../stores/environments-store.service';
import { element } from 'protractor';
import { ISubscription } from 'rxjs/Subscription';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { BadgeShieldsIOBean } from '../../widget/badgewidget/badgewidget.component';

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
  protected applicationStream: Store<ApplicationBean>;
  protected applicationSubscription: ISubscription;
  public application: ApplicationBean;
  private _badgeRatingShields: BadgeShieldsIOBean[];

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
        this.application = app;
        if (this.application.deployments == undefined) {
          this.application.deployments = []
        }
        if (this.application.badgeRatings == undefined){
          this.application.badgeRatings = []
        }
        this._badgeRatingShields = []
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
          )
        })
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }
}
