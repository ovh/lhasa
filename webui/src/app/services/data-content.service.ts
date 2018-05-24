import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import { ConfigurationService } from './configuration.service';
import { DefaultResource, DefaultStreamResource } from '../interfaces/default-resources.interface';
import { DataCoreResource } from './data-core-resources.service';

/**
 * data model
 */
import { ContentBean } from '../models/commons/content-bean';
import { HttpClient } from '@angular/common/http';
import { DataStreamResource } from './data-stream-resources.service';

@Injectable()
export class DataContentService extends DataStreamResource<ContentBean> implements DefaultStreamResource<ContentBean> {
  constructor(
    private _http: HttpClient,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ServerCatalogWithApiUrl + 'v1/contents', _http);
  }
}
