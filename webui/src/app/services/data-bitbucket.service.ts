import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import { ConfigurationService } from './configuration.service';
import { DefaultResource, DefaultLinkResource } from '../interfaces/default-resources.interface';
import { DataCoreResource } from './data-core-resources.service';
import { BitbucketBean } from '../models/commons/applications-bean';

/**
 * data model
 */

@Injectable()
export class BitbucketService extends DataCoreResource<BitbucketBean> implements DefaultResource<BitbucketBean> {
  constructor(
    private _http: Http,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ServerCatalogCompanionWithApiUrl + 'v1/bitbucket', _http);
  }
}
