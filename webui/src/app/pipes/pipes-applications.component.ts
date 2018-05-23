import { Pipe, PipeTransform } from '@angular/core';
import { ApplicationBean, DomainBean } from '../models/commons/applications-bean';

import { sortBy } from 'lodash';

@Pipe({
  name: 'orderByDomains'
})
export class DomainSortPipe implements PipeTransform {
  transform(domains: Array<DomainBean>): Array<DomainBean> {
    const ordered = sortBy(domains, (domain) => {
      return this.pad(domain.name, 'a', 64);
    });
    return ordered;
  }

  pad(str, padString, length) {
    while (str.length < length) {
      str = str + padString;
    }
    return str;
  }
}

@Pipe({
  name: 'orderByApps'
})
export class ApplicationSortPipe implements PipeTransform {
  transform(applications: Array<ApplicationBean>): Array<ApplicationBean> {
    const ordered = sortBy(applications, (app) => {
      return this.pad(app.name, 'a', 64) + this.pad(app.domain, 'a', 64);
    });
    return ordered;
  }

  pad(str, padString, length) {
    while (str.length < length) {
      str = str + padString;
    }
    return str;
  }
}
