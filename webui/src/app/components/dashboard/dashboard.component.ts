import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router, Params } from '@angular/router';
import * as _ from 'lodash';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  constructor(
    private route: ActivatedRoute,
    private router: Router
  ) {
  }

  ngOnInit() {
    this.route.queryParams.forEach((params: Params) => {
      if (params.redirect) {
        const copy: Params = {};
        _.each(params, prm => {
          if (prm !== 'redirect') {
            copy[prm.key] = params[prm];
          }
        });
        this.router.navigate([params.redirect], {
          queryParams: copy
        });
      }
    });
  }

}
