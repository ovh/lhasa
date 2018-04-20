import {Injectable} from '@angular/core';
import {Http} from '@angular/http';
import {ConfigurationService} from './configuration.service';
import {DefaultResource} from '../interfaces/default-resources.interface';
import {DataCoreResource} from './data-core-resources.service';
/**
 * data model
 */
import {EnvironmentBean} from '../models/commons/applications-bean';

@Injectable()
export class DataEnvironmentService extends DataCoreResource<EnvironmentBean> implements DefaultResource<EnvironmentBean> {
  constructor(
    private _http: Http,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ServerCatalogWithApiUrl + 'v1/environments', _http);
  }
}
