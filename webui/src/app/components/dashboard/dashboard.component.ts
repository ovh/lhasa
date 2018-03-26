import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router, Params } from '@angular/router';

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
      if (params.target) {
        console.info("Routing to:", params.target)
        this.router.navigate([params.target], {
          queryParams: {
            domain: params.domain,
            application: params.application
          }
        });
      }
    });
  }

}
