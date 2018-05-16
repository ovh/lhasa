import { Component, OnInit, Output, EventEmitter, Input } from '@angular/core';
import { each } from 'lodash';

import { ActivatedRoute } from '@angular/router';
import { ApplicationBean } from '../../models/commons/applications-bean';
import { Store } from '@ngrx/store';
import { ApplicationsStoreService } from '../../stores/applications-store.service';
import { UiKitStep } from '../../models/kit/progress-tracker';

@Component({
  selector: 'app-oui-progress-tracker',
  templateUrl: './oui-progress-tracker.component.html',
  styleUrls: ['./oui-progress-tracker.component.css']
})
export class OuiProgressTrackerComponent implements OnInit {

  _steps: UiKitStep[] = [];
  current: UiKitStep;

  @Output() select: EventEmitter<any> = new EventEmitter();

  constructor(
  ) {
  }

  /**
   * setter
   */
  @Input() set steps(val: UiKitStep[]) {
    this._steps = val;
    this.selectStep(this._steps[0]);
  }

  ngOnInit() {
  }

  onSelect(event: any, step: UiKitStep) {
    event.data = step;
    this.selectStep(step);
    this.select.emit(event);
  }

  selectStep(sel: UiKitStep) {
    let oneActive = false;
    this.current = sel;
    this._steps.forEach((step) => {
      if (step.id === sel.id) {
        step.status = 'active';
        oneActive = true;
      } else {
        if (oneActive) {
          step.status = 'disabled';
        } else {
          step.status = 'complete';
        }
      }
      switch (step.status) {
        case 'complete':
          step.stepClass = 'oui-progress-tracker__step_complete';
          step.labelClass = 'oui-progress-tracker__label_complete';
          return;
        case 'active':
          step.stepClass = 'oui-progress-tracker__step_active';
          step.labelClass = 'oui-progress-tracker__label_active';
          return;
        case 'disabled':
          step.stepClass = 'oui-progress-tracker__step_disabled';
          step.labelClass = 'oui-progress-tracker__label_disabled';
          return;
        default:
          return;
      }
    });
  }

  index(find) {
    let found = 0;
    let index = 0;
    each(this._steps, (step) => {
      if (step.id === find.id) {
        found = index;
      } else {
      }
      index++;
    });
    return found;
  }

  prev() {
    let p = this.index(this.current) - 1;
    if (p < 0) {
      p = 0;
    }
    this.onSelect({}, this._steps[p]);
  }

  next() {
    let n = this.index(this.current) + 1;
    if (n >= this._steps.length) {
      n = this._steps.length - 1;
    }
    this.onSelect({}, this._steps[n]);
  }
}
