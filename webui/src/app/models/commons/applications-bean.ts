import { EntityBean, PageMetaData, HrefLinks } from './entity-bean';
import { Timestamp } from 'rxjs/operators/timestamp';

// Application
export class ApplicationBean extends EntityBean {
  domain: string;
  name: string;
  type: string;
  language: string;
  project: string;
  repo: string;
  description: string;
  manifest: ManifestBean;
  tags: string[];
  deployments: DeploymentBean[];
  version: string;
}

// Manifest
export class ManifestBean {
  name?: string;
  profile?: string;
  description?: string;
  repository?: string;
  authors?: PersonBean[];
  support?: TeamBean;
}

// Deployment
export class DeploymentBean extends EntityBean {
  id: string;
  properties: Map<string, any>;
  undeployedAt: Date;
}

// Environment
export class EnvironmentBean extends EntityBean {
  slug: string;
  name: string;
  properties: Map<string, any>;
}

// Person
export class PersonBean {
  name: string;
  email: string;
  role: string;
  cisco: string;
}

// Team
export class TeamBean {
  name: string;
  email: string;
  cisco: string;
}

// Domain for page browse
export class ApplicationPagesBean {
  applications: ApplicationBean[] = [];
  metadata: PageMetaData = {
    totalElements: 0,
    totalPages: 0,
    size: 0,
    number: 0
  };
}

// Domain
export class DomainBean extends EntityBean {
  name: string;
  applications?: ApplicationBean[];
  _links?: HrefLinks[];
}

// Domain for page browse
export class DomainPagesBean {
  domains: DomainBean[] = [];
  metadata: PageMetaData = {
    totalElements: 0,
    totalPages: 0,
    size: 0,
    number: 0
  };
}

// Bitbucket
export class BitbucketBean extends EntityBean {
}
