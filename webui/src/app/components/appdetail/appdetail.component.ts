import { Component, OnInit } from '@angular/core';

import { ActivatedRoute } from '@angular/router';
import { ApplicationBean } from '../../models/commons/applications-bean';
import { Store } from '@ngrx/store';
import { ApplicationsStoreService } from '../../stores/applications-store.service';


@Component({
  selector: 'app-appdetail',
  templateUrl: './appdetail.component.html',
  styleUrls: ['./appdetail.component.css']
})
export class AppdetailComponent implements OnInit {

  /**
   * internal streams and store
   */
  protected applicationStream: Store<ApplicationBean>;
  protected application: ApplicationBean

  constructor(
    private applicationsStoreService: ApplicationsStoreService,
    private route: ActivatedRoute) {
    /**
     * subscribe
     */
    this.applicationStream = this.applicationsStoreService.active();

    this.applicationStream.subscribe(
      (element: ApplicationBean) => {
        this.application = element;
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }

  ngOnInit() {
  }

}
