import { Component, OnInit, Output, EventEmitter } from '@angular/core';

import { ActivatedRoute } from '@angular/router';
import { ApplicationBean } from '../../models/commons/applications-bean';
import { Store } from '@ngrx/store';
import { ApplicationsStoreService } from '../../stores/applications-store.service';

@Component({
  selector: 'app-oui-progress-tracker',
  templateUrl: './oui-progress-tracker.component.html',
  styleUrls: ['./oui-progress-tracker.component.css']
})
export class OuiProgressTrackerComponent implements OnInit {

  visible = true;
  @Output() select: EventEmitter<any> = new EventEmitter();

  constructor(
  ) {
  }

  ngOnInit() {
  }

  onSelect(event: any, tabs: string) {
    event.data = tabs;
    this.select.emit(event);
  }

}
