import { Injectable } from '@angular/core';
import { ConfigurationService } from './configuration.service';
import { DefaultResource } from '../interfaces/default-resources.interface';
import { DataCoreResource } from './data-core-resources.service';
/**
 * data model
 */
import { EnvironmentBean } from '../models/commons/applications-bean';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class DataEnvironmentService extends DataCoreResource<EnvironmentBean> implements DefaultResource<EnvironmentBean> {
  constructor(
    private _http: HttpClient,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ApiUrl + 'v1/environments', _http);
  }
}
