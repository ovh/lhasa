import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import { environment } from '../environments/environment';

export class Content {
  content: Application[]
}

export class Application{
  id : number;
  domain: string;
  name: string;
  type: string;
  language: string;
  repositoryurl: string;
  avatarurl: string;
  description: string;
  manifest: Manifest
}

export class Manifest {
  description: string;
}

class ApplicationListResponse{
  _links : {}
  start: number;
  size: number;
  limit: number;
  results: Application[];
}

@Injectable()
export class ApplicationsService {
  constructor(private http: HttpClient) { }
  listApplications(): Observable<ApplicationListResponse[]> {
    return this.http.get<ApplicationListResponse[]>(environment.baseUrlUi+'/api/v1/applications',{ params: {"size": "1000"}});
  }
  getApplication(domain: string, name: string): Observable<Content> {
    return this.http.get<Content>(environment.baseUrlUi+'/api/v1/applications/' + domain + '/' + name,{ params: {"size": "1"}});
  }
}
