import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ConfigurationService } from './configuration.service';
import { DefaultResource } from '../interfaces/default-resources.interface';
import { DataCoreResource } from './data-core-resources.service';
/**
 * data model
 */
import { DeploymentBean } from '../models/commons/applications-bean';

@Injectable()
export class DataDeploymentService extends DataCoreResource<DeploymentBean> implements DefaultResource<DeploymentBean> {
  constructor(
    private _http: HttpClient,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ServerCatalogWithApiUrl + 'v1/deployments', _http);
  }
}
