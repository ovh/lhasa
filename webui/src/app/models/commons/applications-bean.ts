import { EntityBean, PageMetaData, HrefLinks, AbstractPaginatedResource } from './entity-bean';
import { Timestamp } from 'rxjs/operators/timestamp';
import { BadgeLevelBean } from './badges-bean';


// Application
export class ApplicationBean extends EntityBean {
  domain: string;
  name: string;
  type?: string;
  language?: string;
  project?: string;
  repo?: string;
  description?: string;
  manifest: ManifestBean;
  properties: any;
  tags?: string[];
  deployments?: DeploymentBean[];
  badgeRatings?: BadgeRatingBean[];
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

// Domain
export class DomainBean extends EntityBean {
  name: string;
  applications?: ApplicationBean[];
  _links?: HrefLinks[];
}

// Domain for page browse
export class ApplicationPagesBean extends AbstractPaginatedResource {
  applications: ApplicationBean[] = [];
}

// Domain for page browse
export class DomainPagesBean extends AbstractPaginatedResource {
  domains: DomainBean[] = [];
}

// BadgeRatingBean for a gamification badge value for an application version
export class BadgeRatingBean {
  badgeslug: string;
  badgetitle: string;
  value: string;
  comment: string;
  level: BadgeLevelBean;
}
