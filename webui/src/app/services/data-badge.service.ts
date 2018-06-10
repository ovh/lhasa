import { Injectable } from '@angular/core';
import { ConfigurationService } from './configuration.service';
import { DefaultResource } from '../interfaces/default-resources.interface';
import { DataCoreResource } from './data-core-resources.service';
/**
 * data model
 */
import { BadgeBean } from '../models/commons/badges-bean';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class DataBadgeService extends DataCoreResource<BadgeBean> implements DefaultResource<BadgeBean> {
  constructor(
    private _http: HttpClient,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ApiUrl + 'v1/badges', _http);
  }
}
