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
    this.getAppDetail()
  }

  getAppDetail(): void {
    let app_id = this.route.snapshot.paramMap.get('app_id');
    this.applicationsService.getApplication(app_id)
    .subscribe(data => this.app = data);
  }

}
