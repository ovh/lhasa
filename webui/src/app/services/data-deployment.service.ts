import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import { ConfigurationService } from './configuration.service';
import { DefaultResource, DefaultLinkResource } from '../interfaces/default-resources.interface';
import { DataCoreResource } from './data-core-resources.service';

/**
 * data model
 */
import { DeploymentBean } from '../models/commons/applications-bean';

@Injectable()
export class DataDeploymentService extends DataCoreResource<DeploymentBean> implements DefaultResource<DeploymentBean> {
  constructor(
    private _http: Http,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ServerCatalogWithApiUrl + 'v1/deployments', _http);
  }
}
