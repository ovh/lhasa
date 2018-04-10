import { Injectable } from '@angular/core';
import 'rxjs/add/operator/map';
import { Observable } from 'rxjs/Observable';
import { Router } from '@angular/router';
import { Http, Response, Headers } from '@angular/http';

import { ConfigurationService } from '../services/configuration.service';
import { DefaultResource } from '../interfaces/default-resources.interface';

/**
 * data model
 */
import { EntityBean } from '../models/commons/entity-bean';
import { DataCoreResource } from './data-core-resources.service';

@Injectable()
export class SecurityService extends DataCoreResource<EntityBean> implements DefaultResource<EntityBean> {

  /**
   * constructor
   */
  constructor(
    private _http: Http,
    private _router: Router,
    private _ConfigurationService: ConfigurationService) {
    super(_ConfigurationService, _ConfigurationService.ServerAppCatalogCompanionWithUrl, _http);
  }

  /**
   * get connect resource
   */
  public Connect = (): Observable<boolean> => {
    return this.http.get(this.actionUrl + 'api/connect', { headers: this.headers })
      .map((response: Response) => <boolean> response.json())
      .catch(this.handleError);
  }
}
