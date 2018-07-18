import { Component, Input, OnInit } from '@angular/core';
import { EnvironmentsStoreService } from '../../stores/environments-store.service';
import { EnvironmentBean } from '../../models/commons/applications-bean';
import { Store } from '@ngrx/store';
import { SubscriptionLike as ISubscription } from 'rxjs';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { Observable } from 'rxjs/internal/Observable';

@Component({
  selector: 'app-env-chip',
  templateUrl: './env-chip.component.html',
  styleUrls: ['./env-chip.component.css', './flags.css'],
})
@AutoUnsubscribe()
export class EnvChipComponent implements OnInit {

  @Input() slug: string;

  public blankImageURL = require('./blank.gif');

  protected environmentStream: Observable<Map<string, EnvironmentBean>>;
  protected environmentSubscription: ISubscription;

  public environments: Map<string, EnvironmentBean>;

  constructor(
    private environmentsStoreService: EnvironmentsStoreService) {
    this.environmentStream = this.environmentsStoreService.environments();
  }

  ngOnInit() {
    this.environmentSubscription = this.environmentStream.subscribe(
      (element: Map<string, EnvironmentBean>) => {
        this.environments = element;
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }
}
