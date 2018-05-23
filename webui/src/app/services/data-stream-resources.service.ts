import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';
import { Observable } from 'rxjs/Observable';

import { DefaultStreamResource } from '../interfaces/default-resources.interface';
import { ConfigurationService } from '../services/configuration.service';
/**
 * data model
 */
import { HttpClient, HttpErrorResponse, HttpHeaders, HttpParams } from '@angular/common/http';

/**
 * default class to handle default behaviour for streaml resource
 * component
 */
export class DataStreamResource<T> implements DefaultStreamResource<T> {

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
    this.headers.append('Content-Type', 'text/plain');
    this.headers.append('Accept', 'text/plain');
    this.headers.append('AuthToken', this.configuration.getAuthToken());
  }

  /**
   * get single resource
   */
  public GetSingle = (id: string): Observable<T> => {
    this.headers.set('AuthToken', this.configuration.getAuthToken());
    return this.http.get(this.actionUrl + '/' + id, {headers: this.headers, responseType: 'text'})
      .catch(this.handleError);
  }

  /**
   * error handler
   */
  protected handleError(error: HttpErrorResponse) {
    return Observable.throw(error || 'Server error');
  }
}
