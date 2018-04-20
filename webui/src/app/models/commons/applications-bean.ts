import {EntityBean} from './entity-bean';
import {Timestamp} from 'rxjs/operators/timestamp';

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
}

// Manifest
export class ManifestBean {
  name: string;
  profile: string;
  description: string;
  repository: string;
  authors: PersonBean[];
  support: TeamBean;
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

// Domain
export class DomainBean {
  name: string;
  applications: ApplicationBean[];
}

// Bitbucket
export class BitbucketBean extends EntityBean {
}
