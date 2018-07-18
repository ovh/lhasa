
import {catchError} from 'rxjs/operators';
import { Injectable } from '@angular/core';

import { Observable } from 'rxjs';
import { Router } from '@angular/router';

import { ConfigurationService } from '../services/configuration.service';
import { DefaultResource } from '../interfaces/default-resources.interface';
/**
 * data model
 */
import { EntityBean } from '../models/commons/entity-bean';
import { DataCoreResource } from './data-core-resources.service';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class SecurityService extends DataCoreResource<EntityBean> implements DefaultResource<EntityBean> {

  /**
   * constructor
   */
  constructor(
    private _http: HttpClient,
    private _router: Router,
    private _ConfigurationService: ConfigurationService) {
    super(_ConfigurationService, _ConfigurationService.ApiUrl, _http);
  }

  /**
   * get connect resource
   */
  public Connect = (): Observable<Object> => {
    return this.http.get(this.actionUrl + 'api/connect', {headers: this.headers}).pipe(
      catchError(this.handleError));
  }
}
