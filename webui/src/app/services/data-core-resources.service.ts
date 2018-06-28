import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';
import { Observable } from 'rxjs/Observable';

import { DefaultResource } from '../interfaces/default-resources.interface';
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
export class DataCoreResource<T extends EntityBean> implements DefaultResource<T> {

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
   * get all resources
   */
  public GetAll = (): Observable<T[]> => {
    this.headers.set('AuthToken', this.configuration.getAuthToken());
    return this.http.get(this.actionUrl, {headers: this.headers})
      .catch(this.handleError);
  }

  /**
   * get all resources
   */
  public GetAllFromContent = (filter: string, params: { [p: string]: string }): Observable<ContentListResponse<T>> => {
    this.headers.set('AuthToken', this.configuration.getAuthToken());
    return this.http.get<ContentListResponse<T>>(this.actionUrl + filter,
      {headers: this.headers, params: new HttpParams({fromObject: params})})
      .catch(this.handleError);
  }

  /**
   * get single resource
   */
  public GetSingle = (id: string): Observable<T> => {
    this.headers.set('AuthToken', this.configuration.getAuthToken());
    return this.http.get<Observable<T>>(this.actionUrl + '/' + id, {headers: this.headers})
      .catch(this.handleError);
  }

  /**
   * get single resource with any reseult
   */
  public GetSingleAny = (id: string): Observable<any> => {
    this.headers.set('AuthToken', this.configuration.getAuthToken());
    return this.http.get<Observable<T>>(this.actionUrl + '/' + id, {headers: this.headers})
      .catch(this.handleError);
  }

  /**
   * call api backend with a POST to execute task
   */
  public Task = (path: String, payload: any): Observable<T> => {
    this.headers.set('AuthToken', this.configuration.getAuthToken());
    return this.http.post(this.actionUrl + '/' + path, JSON.stringify(payload), {headers: this.headers})
      .catch(this.handleError);
  }

  /**
   * add a new resource
   */
  public Add = (itemToAdd: T): Observable<T> => {
    this.headers.set('AuthToken', this.configuration.getAuthToken());
    return this.http.post(this.actionUrl, JSON.stringify(itemToAdd), {headers: this.headers})
      .catch(this.handleError);
  }

  /**
   * update this resource
   */
  public Update = (id: string, itemToUpdate: T): Observable<T> => {
    this.headers.append('AuthToken', this.configuration.getAuthToken());
    return this.http.put(this.actionUrl + '/' + id, JSON.stringify(itemToUpdate), {headers: this.headers})
      .catch(this.handleError);
  }

  /**
   * delete this resource
   */
  public Delete = (id: string): Observable<T> => {
    this.headers.set('AuthToken', this.configuration.getAuthToken());
    return this.http.delete(this.actionUrl + '/' + id, {headers: this.headers})
      .catch(this.handleError);
  }

  /**
   * error handler
   */
  protected handleError(error: HttpErrorResponse) {
    return Observable.throw(error || 'Server error');
  }
}
