import { Injectable } from '@angular/core';
import { ConfigurationService } from './configuration.service';
import { DefaultResource } from '../interfaces/default-resources.interface';
import { DataCoreResource } from './data-core-resources.service';
/**
 * data model
 */
import { DomainBean } from '../models/commons/applications-bean';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class DataDomainService extends DataCoreResource<DomainBean> implements DefaultResource<DomainBean> {
  constructor(
    private _http: HttpClient,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ServerCatalogWithApiUrl + 'v1/domains', _http);
  }
}
