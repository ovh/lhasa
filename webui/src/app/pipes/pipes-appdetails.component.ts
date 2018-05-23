import { Pipe, PipeTransform } from '@angular/core';
import { DeploymentBean } from '../models/commons/applications-bean';

@Pipe({
  name: 'activeDeployment'
})
export class AppdetailsActiveDeploymentsPipe implements PipeTransform {
  transform(deployments: Array<DeploymentBean>): Array<DeploymentBean> {
    const filtered = [];
    deployments.forEach(value => {
      if (!value.undeployedAt || value.undeployedAt === null) {
        filtered.push(value);
      }
    });
    return filtered;
  }
}
