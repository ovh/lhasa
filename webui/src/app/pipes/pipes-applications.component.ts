import { Component, OnInit, Pipe } from '@angular/core';
import { ApplicationBean, DomainBean } from '../models/commons/applications-bean';

import * as _ from 'lodash';

@Pipe({
    name: "orderByDomains"
})
export class DomainSortPipe {
    transform(domains: Array<DomainBean>): Array<DomainBean> {
        let ordered = _.sortBy(domains, (domain) => {
            return this.pad(domain.name, "a", 64)
        });
        return ordered;
    }

    pad(str, padString, length) {
        while (str.length < length)
            str = str + padString;
        return str;
    }
}

@Pipe({
    name: "orderByApps"
})
export class ApplicationSortPipe {
    transform(applications: Array<ApplicationBean>): Array<ApplicationBean> {
        let ordered = _.sortBy(applications, (app) => {
            return this.pad(app.name, "a", 64) + this.pad(app.domain, "a", 64)
        });
        return ordered;
    }

    pad(str, padString, length) {
        while (str.length < length)
            str = str + padString;
        return str;
    }
}
