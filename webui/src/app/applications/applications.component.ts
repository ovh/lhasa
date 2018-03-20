import { Component, OnInit } from '@angular/core';
import { ApplicationsService, Application } from '../app.service';

@Component({
  selector: 'app-applications',
  templateUrl: './applications.component.html',
  styleUrls: ['./applications.component.css']
})
export class ApplicationsComponent implements OnInit {

  applications: Application[] = []
  
  constructor(private applicationsService: ApplicationsService) { 
    
  }

  ngOnInit() {
    this.getApplications()
  }

  getApplications(): void {
    this.applicationsService.listApplications()
    .subscribe(
      (data) => {
        this.applications = data['content']
      }
    );
  }
}
