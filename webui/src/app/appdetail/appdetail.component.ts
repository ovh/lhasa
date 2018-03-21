import { Component, OnInit } from '@angular/core';

import { ApplicationsService, Application } from '../app.service';
import {ActivatedRoute} from '@angular/router';


@Component({
  selector: 'app-appdetail',
  templateUrl: './appdetail.component.html',
  styleUrls: ['./appdetail.component.css']
})
export class AppdetailComponent implements OnInit {

  app: Application
  
  constructor(private applicationsService: ApplicationsService, private route: ActivatedRoute) { 
    
  }

  ngOnInit() {
    this.app = <Application> {}
    this.getAppDetail()
  }

  getAppDetail(): void {
    let domain = this.route.snapshot.paramMap.get('domain');
    let name = this.route.snapshot.paramMap.get('name');
    this.applicationsService.getApplication(domain, name)
    .subscribe(
      (data) => {
        this.app = data.content[0]
      });
  }

}
