import {catchError} from 'rxjs/operators';

import {throwError as observableThrowError,  Observable } from 'rxjs';

import { DefaultResource, DefaultGraphResource } from '../interfaces/default-resources.interface';
import { ConfigurationService } from '../services/configuration.service';
/**
 * data model
 */
import { ContentListResponse, EntityBean } from '../models/commons/entity-bean';
import { HttpClient, HttpErrorResponse, HttpHeaders, HttpParams } from '@angular/common/http';

/**
 * default class to handle default behaviour or resource
 * component
 */
export class DataGraphResource<T> implements DefaultGraphResource<T> {

  protected actionUrl: string;
  protected headers: HttpHeaders;
  protected http: HttpClient;
  protected configuration: ConfigurationService;

  /**
   * constructor
   */
  constructor(_configuration: ConfigurationService, actionUrl: string, _http: HttpClient) {
    this.http = _http;
    this.actionUrl = actionUrl;
    this.configuration = _configuration;

    this.headers = new HttpHeaders();
    this.headers.append('Content-Type', 'application/json');
    this.headers.append('Accept', 'application/json');
    this.headers.append('AuthToken', this.configuration.getAuthToken());
  }

  /**
   * get single resource
   */
  public Get = (params: any): Observable<any> => {
    this.headers.set('AuthToken', this.configuration.getAuthToken());
    return this.http.get(this.actionUrl , {headers: this.headers}).pipe(
      catchError(this.handleError));
  }

  /**
   * error handler
   */
  protected handleError(error: HttpErrorResponse) {
    return observableThrowError(error || 'Server error');
  }
}
