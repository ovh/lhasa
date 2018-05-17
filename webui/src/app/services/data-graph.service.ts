import { DataGraphResource } from './data-graph-resources.service';
import { Injectable } from '@angular/core';
import { ConfigurationService } from './configuration.service';
import { DefaultGraphResource } from '../interfaces/default-resources.interface';
/**
 * data model
 */
import { HttpClient } from '@angular/common/http';
import { GraphBean } from '../models/graph/graph-bean';

@Injectable()
export class DataGraphService extends DataGraphResource<GraphBean> implements DefaultGraphResource<GraphBean> {
  constructor(
    private _http: HttpClient,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ApiUrl + 'v1/graphs', _http);
  }
}
