import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { ConfigurationService } from './configuration.service';
import { DefaultResource } from '../interfaces/default-resources.interface';
import { DataCoreResource } from './data-core-resources.service';
import { BadgeRatingBean } from '../models/commons/applications-bean';

/**
 * data model
 */


@Injectable()
export class DataBadgeRatingsService{
  private http 
  private headers
  private configuration
  private actionUrl
  constructor(
    private _http: HttpClient,
    private _configuration: ConfigurationService
  ) {
    this.http = _http;
    this.configuration = _configuration;
    this.headers = new HttpHeaders();
    this.headers.append('Content-Type', 'application/json');
    this.headers.append('Accept', 'application/json');
    this.headers.append('AuthToken', this.configuration.getAuthToken());
    this.actionUrl = _configuration.ApiUrl  + 'v1/applications'
  }

  /**
   * get all badge values
   */
  public GetBadgeRatings = (id: string): Observable<BadgeRatingBean[]> => {
    this.headers.set('AuthToken', this.configuration.getAuthToken());
    return this.http.get(this.actionUrl + '/' + id, {headers: this.headers})
      .catch(this.handleError);
  }

  
  /**
   * error handler
   */
  protected handleError(error: HttpErrorResponse) {
    return Observable.throw(error || 'Server error');
  }
}
