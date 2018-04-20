import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {EnvironmentsStoreService} from '../../stores/environments-store.service';
import {EnvironmentBean} from '../../models/commons/applications-bean';
import {Store} from '@ngrx/store';
import {ISubscription} from 'rxjs/Subscription';

@Component({
  selector: 'app-env-chip',
  templateUrl: './env-chip.component.html',
  styleUrls: ['./env-chip.component.css', './flags.css'],
})
export class EnvChipComponent implements OnInit, OnDestroy {

  @Input() slug: string;

  public blankImageURL = require('./blank.gif');

  protected environmentStream: Store<Map<string, EnvironmentBean>>;
  protected environmentSubscription: ISubscription;

  protected environments: Map<string, EnvironmentBean>;

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

  ngOnDestroy(): void {
    this.environmentSubscription.unsubscribe();
  }
}
