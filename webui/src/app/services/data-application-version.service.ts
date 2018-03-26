import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import { ConfigurationService } from './configuration.service';
import { DefaultResource, DefaultLinkResource } from '../interfaces/default-resources.interface';
import { DataCoreResource } from './data-core-resources.service';

/**
 * data model
 */
import { ApplicationBean } from '../models/commons/applications-bean';

@Injectable()
export class DataApplicationServiceService extends DataCoreResource<ApplicationBean> implements DefaultResource<ApplicationBean> {
  constructor(
    private _http: Http,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ServerWithApiUrl + 'v1/applications', _http);
  }
}
